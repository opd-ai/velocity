# Ebitengine Game Audit Report
Generated: 2026-04-02T03:22:32Z

## Executive Summary
- **Total Issues**: 42
- **Critical**: 4 - Crashes, game-breaking bugs
- **High**: 8 - Major functionality/UX problems
- **Medium**: 10 - Noticeable bugs, moderate impact
- **Low**: 6 - Minor issues, edge cases
- **Optimizations**: 8 - Performance improvements
- **Code Quality**: 6 - Maintainability concerns

## Critical Issues

### [C-001] Global mutable state for key tracking causes data races
- **Location**: `pkg/game/game.go:555`
- **Category**: State
- **Description**: `prevKeys` is a package-level `var` (`var prevKeys = make(map[ebiten.Key]bool)`) shared across all `Game` instances. Since it is written to in `updatePrevKeys()` and read in `wasKeyPressed()`, any concurrent access (e.g., multiple Game instances in tests, or goroutine access) constitutes a data race. Even in single-instance use, package-level mutable state violates encapsulation and prevents safe testing.
- **Impact**: Undefined behavior under `-race`, potential crashes if multiple Game instances exist.
- **Reproduction**:
  1. Create two `Game` instances (e.g., in a test).
  2. Call `Update()` on both concurrently.
  3. Race detector triggers on `prevKeys` map access.
- **Root Cause**: Key-tracking state is stored in a package-level variable instead of on the `Game` struct.
- **Suggested Fix**: Move `prevKeys` into the `Game` struct as a field (e.g., `prevKeyState map[ebiten.Key]bool`), and pass it to `handleMenuInput` / `updatePrevKeys`.

### [C-002] Pause input not edge-detected -- toggles every frame while held
- **Location**: `pkg/game/game.go:546-548`
- **Category**: Input
- **Description**: The pause toggle checks `g.inputSystem.IsPausePressed()` which returns `true` every frame the Escape key is held. There is no edge detection (press-to-release transition). This causes rapid pause/unpause toggling at 60 FPS while Escape is held.
- **Impact**: Game becomes unplayable -- pause state oscillates rapidly, visual flickering, player cannot reliably pause.
- **Reproduction**:
  1. Start a game and enter playing state.
  2. Press and hold the Escape key.
  3. Observe rapid toggling between playing and paused states.
- **Root Cause**: `IsPausePressed()` returns the raw held state; the pause check at line 546 lacks edge detection unlike the menu input handling which uses `wasKeyPressed()`.
- **Suggested Fix**: Apply the same edge-detection pattern used in `handleMenuInput()` -- only trigger pause on the rising edge (pressed this frame, not pressed last frame).

### [C-003] Retry from Game Over state fails -- StartGame() guards against non-MainMenu transitions
- **Location**: `pkg/ux/game_state.go:78-82`, `pkg/ux/game_state.go:245-247`
- **Category**: Logic
- **Description**: `MenuController.handleAction("retry")` calls `stateManager.StartGame()`, but `StartGame()` only transitions if `state == StateMainMenu`. When the game is in `StateGameOver`, `StartGame()` silently does nothing. The "Retry" menu option is therefore non-functional.
- **Impact**: Game-breaking -- after dying, the player cannot restart; they are permanently stuck on the Game Over screen unless they select "Main Menu" first.
- **Reproduction**:
  1. Play until the player dies (Game Over screen).
  2. Select "Retry" from the Game Over menu.
  3. Nothing happens.
- **Root Cause**: `StartGame()` has a state guard that only allows transition from `StateMainMenu`, but `handleAction("retry")` calls it from `StateGameOver`.
- **Suggested Fix**: Either add `StateGameOver` as an allowed source state in `StartGame()`, or have `handleAction("retry")` call `ReturnToMainMenu()` followed by `StartGame()`, or add a dedicated `RestartGame()` method.

### [C-004] Double movement for projectiles -- physics system and ProjectileSystem both update position
- **Location**: `pkg/combat/projectile.go:77,90-102`, `pkg/engine/physics.go:55-88`, `pkg/game/game.go:300,606,613`
- **Category**: Logic
- **Description**: Projectiles have both `position` and `velocity` components. The `PhysicsSystem` (registered with the World at line 300, called via `world.Update(dt)` at line 606) applies `pos += vel * dt`. Then `ProjectileSystem.Update()` (line 613) also calls `moveProjectile()` which does `pos += vel * dt` again. Projectiles move at **2x intended speed**.
- **Impact**: Projectiles travel twice as fast as designed, breaking weapon balance, making combat feel wrong, and reducing projectile lifetime distance (they leave the screen faster).
- **Reproduction**:
  1. Fire a projectile and observe its speed.
  2. Compare against `projectileSpeed = 400.0` -- it will move at ~800 pixels/sec.
- **Root Cause**: Both the generic ECS `PhysicsSystem` and the domain-specific `ProjectileSystem` independently update position from velocity.
- **Suggested Fix**: Remove the `moveProjectile()` call from `ProjectileSystem.Update()`, or exclude projectile entities from `PhysicsSystem` by adding a tag check, or don't give projectiles a `velocity` component used by PhysicsSystem (use a separate component).

## High Priority Issues

### [H-001] SFX audio data regenerated on every playback -- no caching
- **Location**: `pkg/audio/audio.go:311-329`
- **Category**: Performance / Assets
- **Description**: `GetSFXData(name)` is called every time an SFX plays (via `playSFXNow()`). Each call regenerates the PCM byte buffer from scratch using synthesis functions (`GenerateLaserSFX()`, `GenerateExplosionSFX()`, etc.). These are CPU-intensive operations involving per-sample math loops.
- **Impact**: Repeated allocations and heavy math during gameplay. In combat with many explosions, this causes per-frame allocation spikes and potential frame drops.
- **Suggested Fix**: Cache the generated PCM data in a `map[string][]byte` on first generation. SFX data is deterministic and identical every time, so caching is safe.

### [H-002] Drag coefficient applied per-frame instead of per-second -- physics is frame-rate dependent
- **Location**: `pkg/engine/physics.go:73-75`
- **Category**: Logic
- **Description**: `vel.VX *= ps.config.DragCoeff` is applied each `Update()` call. With `DragCoeff = 0.98`, this means velocity is multiplied by `0.98` every frame. At 60 FPS, effective per-second drag is `0.98^60 ~ 0.30`. If TPS changes (e.g., lag causes fewer ticks), drag behavior changes. The same issue affects particle drag at `rendering.go:286-287`.
- **Impact**: Ship handling feel changes if frame rate deviates from 60 FPS. Speed decay is not time-based.
- **Root Cause**: Drag is applied as a per-frame multiplier without incorporating `dt`.
- **Suggested Fix**: Use `vel.VX *= math.Pow(ps.config.DragCoeff, dt * 60)` or convert to exponential decay: `vel.VX *= math.Exp(-dragRate * dt)`.

### [H-003] BoundingBox offset semantics are inconsistent -- collision detection is unreliable
- **Location**: `pkg/game/game.go:360-362`, `pkg/combat/projectile.go:132-133,182-183,218-219`, `pkg/combat/projectile.go:187-191`
- **Category**: Collision
- **Description**: For the player, the bounding box is set with `X: -8, Y: -8` (offset from position). But in `getProjectileBounds()` when no box component exists, the default box is `{X: pos.X - 2, Y: pos.Y - 2, Width: 4, Height: 4}` -- mixing absolute position into the offset field. Then `boxesOverlap()` does `posA.X + boxA.X`, which for the default case becomes `pos.X + (pos.X - 2) = 2*pos.X - 2`, placing the collision box far off from the entity.
- **Impact**: Projectiles without explicit bounding box components have collision boxes displaced far from their actual position, making hit detection essentially random.
- **Root Cause**: Inconsistent convention -- sometimes `BoundingBox.X/Y` is a relative offset from entity position, sometimes it's absolute.
- **Suggested Fix**: Always use relative offsets in `BoundingBox`. Change the default in `getProjectileBounds()` to `&BoundingBox{X: -2, Y: -2, Width: 4, Height: 4}` and similarly for `getTargetBounds()`.

### [H-004] Sprite cache key does not include size -- different-sized sprites overwrite each other
- **Location**: `pkg/rendering/rendering.go:130-132`, `pkg/game/game.go:862`
- **Category**: Assets
- **Description**: `keyString(SpriteKey)` generates `"%s:%d:%d"` from `GenreID`, `Type`, and `Variant`. The `Size` parameter is not included. If `GetOrCreateShipSprite(variant=0, size=16)` and `GetOrCreateShipSprite(variant=0, size=32)` are called, they would collide in cache.
- **Impact**: If any future code requests sprites at different sizes with the same type and variant, it gets the wrong cached image. Currently a latent bug.
- **Suggested Fix**: Include `Size` in the `SpriteKey` struct and in `keyString()`.

### [H-005] Enemy entities spawn off-screen and may never enter the arena in `wrap` mode
- **Location**: `pkg/procgen/spawner.go:94-117`, `pkg/engine/arena.go:61-72`
- **Category**: Logic
- **Description**: Enemies spawn at positions like `y = -50` (50px off-screen). In wrap mode, the arena system teleports entities with `pos.Y < 0` to `pos.Y += height`. So enemies spawn, immediately get wrapped to the bottom of the screen, then start moving toward the player. This creates an unintended spawn pattern where all "top" spawns appear at the bottom.
- **Impact**: Enemy spawn patterns are incorrect in wrap mode. All top-edge spawns appear at the bottom, right-edge spawns appear at the left, etc.
- **Root Cause**: Arena wrap logic doesn't distinguish between intentional off-screen entities (spawning) and entities that have crossed the boundary during gameplay.
- **Suggested Fix**: Either delay arena processing for newly spawned enemies (e.g., with a "spawning" flag), or spawn enemies at the screen edge (margin = 0), or exempt spawning entities from arena wrapping.

### [H-006] Audio backend `Initialize()` called every frame
- **Location**: `pkg/audio/audio.go:136-140`
- **Category**: Performance
- **Description**: `Manager.Update()` calls `m.audioBackend.Initialize()` on every frame. While `ebitenBackend.Initialize()` has an `if b.initialized { return }` early-exit, this is still a method call + nil check + field read on every frame for the lifetime of the game.
- **Impact**: Minor per-frame overhead, but the pattern is fragile -- if the backend implementation changes, it could accidentally re-initialize.
- **Suggested Fix**: Call `Initialize()` once during `NewManager()` or on first `PlaySFX`/`PlayMusic` call. Use a `sync.Once` or an initialization flag on the Manager.

### [H-007] `loadAndResumeGame` calls `StartNextWave` in a loop without clearing enemies -- entities pile up
- **Location**: `pkg/game/game.go:474-479`
- **Category**: Logic
- **Description**: When loading a save, the code advances waves by calling `StartNextWave()` in a loop for `state.Wave - 1` iterations, then `clearEnemies()` after each. But `StartNextWave()` increments the wave counter AND spawns enemies. The inner `clearEnemies()` removes them, but the entity IDs are consumed. This is wasteful and may leave orphan components.
- **Impact**: Temporary entity churn during load, potential for leaked entity state.
- **Suggested Fix**: Add a `WaveManager.SetWave(n int)` method that sets `currentWave` directly without spawning.

### [H-008] `cosApprox` and `sinApprox` Taylor series are inaccurate for large angles
- **Location**: `pkg/rendering/rendering.go:332-343`
- **Category**: Logic
- **Description**: The Taylor series approximations `cosApprox` and `sinApprox` are only accurate near x=0. For angles like pi (3.14) or 2*pi (6.28), the error is significant. Particles emitted at angles > 1 radian will have noticeably wrong velocities. For example, `sinApprox(pi) ~ 3.14 - 5.17 + 2.55 = 0.52` instead of `0.0`.
- **Impact**: Particle emission directions are distorted for most angles, causing asymmetric or clumped particle effects instead of uniform radial explosions.
- **Root Cause**: Taylor series requires angle normalization to [-pi, pi] range first.
- **Suggested Fix**: Use `math.Cos` and `math.Sin` directly (they are fast on modern hardware), or normalize the input angle before applying the approximation.

## Medium Priority Issues

### [M-001] Combo display always shows +1 more than actual multiplier
- **Location**: `pkg/game/game.go:901`
- **Category**: UI
- **Description**: HUD displays `g.hud.Combo+1` as the combo multiplier text, but the actual score multiplier at line 413 is `1 + g.combo/ComboTierDivisor`. These formulas diverge -- the HUD shows raw `combo+1`, while actual scoring uses a tiered system. For combo=3, HUD shows "x4" but actual multiplier is 1 (3/5=0, 1+0=1).
- **Impact**: Player sees a misleading combo multiplier that doesn't match actual scoring.
- **Suggested Fix**: Display the actual multiplier: `1 + g.hud.Combo/ComboTierDivisor`.

### [M-002] Particle rendering uses per-pixel Set() calls -- extremely slow
- **Location**: `pkg/game/game.go:772-783`
- **Category**: Performance / Rendering
- **Description**: `renderParticlePixels()` calls `screen.Set(x, y, col)` for each pixel of each particle. In Ebitengine, `Set()` is extremely slow as it must lock the image, convert the color, and write one pixel at a time. With 20 particles x ~4-9 pixels each = up to 180 `Set()` calls per explosion.
- **Impact**: Significant frame-time cost during combat. Multiple simultaneous explosions could drop FPS below 60.
- **Suggested Fix**: Use a small pre-generated `ebiten.Image` for particles and draw it with `DrawImage()` + tint via `ColorScale`. Or batch all particle pixels into a single `WritePixels()` call.

### [M-003] Menu selection index not reset on state change
- **Location**: `pkg/ux/game_state.go:202-224,227-256`
- **Category**: UI
- **Description**: When transitioning from Pause menu (3 items) to Game Over menu (2 items), the `selectedIdx` may be 2 (third item). Since Game Over only has 2 items, `selectedIdx` is out of bounds until the player navigates. The `Select()` method at line 229 guards against this, but the visual selection indicator renders at the wrong position.
- **Impact**: Visual glitch -- selection cursor appears on a non-existent menu item after state transitions.
- **Suggested Fix**: Reset `selectedIdx = 0` in the `transition()` method, or in `GetCurrentItems()` clamp it.

### [M-004] Entity removal during iteration -- modifying `entities` map while iterating
- **Location**: `pkg/engine/engine.go:66-70`, `pkg/combat/damage.go:90-94`
- **Category**: State
- **Description**: `DamageSystem.Update()` calls `world.RemoveEntity()` at line 94 while other systems may be iterating via `ForEachEntity`. While the current call ordering avoids direct concurrent iteration, the pattern is fragile. Any future system that iterates entities while damage processing occurs could hit a map modification during iteration panic.
- **Impact**: Currently safe due to sequential system execution, but a latent bug if system ordering changes.
- **Suggested Fix**: Collect entities to remove and defer actual removal to end-of-frame.

### [M-005] Seed 0 produces identical content across all sessions -- no randomization
- **Location**: `pkg/config/config.go:97`, `pkg/procgen/procgen.go:32`
- **Category**: Logic
- **Description**: The default seed is 0. When seed is 0, `GenerateWave()` produces `Seed: 0 + waveNumber`, meaning wave 1 always seeds as 1, wave 2 as 2, etc. Every player with default config gets identical enemy patterns every game.
- **Impact**: No variety for default configuration.
- **Suggested Fix**: In `config.Load()` or `NewGame()`, if seed is 0, generate a random seed from `time.Now().UnixNano()`.

### [M-006] `InputSystem` not registered with World but manually called -- inconsistent pattern
- **Location**: `pkg/game/game.go:605-606`
- **Category**: Logic
- **Description**: `g.inputSystem.Update(dt)` is called explicitly at line 605, and `g.world.Update(dt)` is called at line 606. But `inputSystem` is NOT registered with `world.AddSystem()` (only physics and arena are, lines 300-301). Some systems are in the world, others are called manually.
- **Impact**: No bug currently, but confusing architecture. Adding new systems may lead to double-update errors.
- **Suggested Fix**: Document the pattern clearly or standardize all system updates.

### [M-007] `EbitenInputReader` screen size defaults don't update from config
- **Location**: `pkg/engine/input_ebiten.go:20-24`, `pkg/game/game.go:247-248`
- **Category**: Input
- **Description**: `NewEbitenInputReader()` defaults to 800x600 for touch region calculation. The `Game` never calls `inputReader.SetScreenSize()` with the actual config dimensions. If the config specifies a different resolution, touch regions will be calculated incorrectly.
- **Impact**: Touch input zones are wrong for any non-800x600 resolution.
- **Suggested Fix**: Call `inputReader.SetScreenSize(cfg.Display.Width, cfg.Display.Height)` after creating the input reader.

### [M-008] `SpriteCache.GetOrCreate` has a TOCTOU race condition
- **Location**: `pkg/rendering/rendering.go:150-157`
- **Category**: State
- **Description**: `GetOrCreate()` calls `Get()` (acquires RLock, checks, releases), then if missing calls `gen()` and `Set()` (acquires Lock). Between the RLock release and Lock acquisition, another goroutine could also find the key missing and call `gen()`. Both would generate and store the sprite, wasting work.
- **Impact**: Potential duplicate sprite generation in concurrent scenarios.
- **Suggested Fix**: Hold the write lock for the entire check-and-create operation.

### [M-009] `parseKey` returns `KeyW` as default for unknown key names
- **Location**: `pkg/engine/input_ebiten.go:159-164`
- **Category**: Input
- **Description**: If a user configures an unsupported key name in `config.yaml` (e.g., `fire: "Tab"`), `parseKey()` silently returns `ebiten.KeyW`. This means the unsupported binding becomes the thrust key.
- **Impact**: Misconfigured controls silently map to thrust, causing confusing behavior.
- **Suggested Fix**: Log a warning for unrecognized key names, or return a sentinel value.

### [M-010] Player death while paused can still process via death callback
- **Location**: `pkg/game/game.go:541,614`, `pkg/combat/damage.go:89-95`
- **Category**: State
- **Description**: Potential issue -- if damage is queued before pause and then `damageSystem.Update(dt)` is called on the same frame the pause occurs, a death could process while transitioning to pause.
- **Impact**: Edge case where player death and pause occur on the same frame could lead to unexpected state.
- **Suggested Fix**: Check if still playing before processing death callbacks, or queue state transitions.

## Low Priority Issues

### [L-001] `Layout()` ignores `outsideWidth` and `outsideHeight` parameters
- **Location**: `pkg/game/game.go:983-985`
- **Category**: Ebitengine
- **Description**: `Layout()` returns fixed dimensions from config, ignoring the provided `outsideWidth` and `outsideHeight`. This means the game doesn't adapt its logical resolution to window resizes.
- **Impact**: Minor -- game uses fixed internal resolution which Ebitengine scales. Valid design choice for pixel-art games.
- **Suggested Fix**: Acceptable as-is. Document that the game uses fixed logical resolution.

### [L-002] `audit.Logger.samples` grows unbounded
- **Location**: `pkg/audit/audit.go:25-30`
- **Category**: Performance
- **Description**: `Logger.Record()` appends to `samples` without any size limit. In a long play session, this slice grows continuously, consuming memory.
- **Impact**: Slow memory leak proportional to play duration. At ~60 samples/sec with ~40 bytes each, ~8.6 MB/hour.
- **Suggested Fix**: Add a ring buffer or cap the slice at a maximum size.

### [L-003] `Weapon.timer` can go negative without bound
- **Location**: `pkg/combat/combat.go:42-44`
- **Category**: Logic
- **Description**: `Weapon.Update(dt)` decrements `timer` every frame even when already <= 0. Over time, `timer` becomes a large negative number.
- **Impact**: Minimal -- float64 underflow is not practically reachable. Code hygiene issue.
- **Suggested Fix**: Clamp: `if w.timer > 0 { w.timer -= dt; if w.timer < 0 { w.timer = 0 } }`.

### [L-004] `os.Exit(0)` called directly in game callback -- bypasses deferred cleanup
- **Location**: `pkg/game/game.go:216`
- **Category**: State
- **Description**: The "quit" menu action calls `os.Exit(0)`, which terminates the process immediately without running deferred functions (e.g., the recovery handler in `main()`).
- **Impact**: Minor -- cleanup functions don't run. Could miss flushing logs or audit data.
- **Suggested Fix**: Return a sentinel error from `Update()` to let Ebitengine shut down gracefully.

### [L-005] `clearAllEntities` iterates then removes -- two passes over entity map
- **Location**: `pkg/game/game.go:381-389`
- **Category**: Performance
- **Description**: `clearAllEntities()` first collects all entity IDs into a slice, then iterates that slice to remove each. For a full clear, replacing the map would be more efficient.
- **Impact**: Negligible for typical entity counts (<200).
- **Suggested Fix**: Add a `World.Clear()` method that replaces the `entities` map with a new empty one.

### [L-006] `encodePlayerHealth` / `decodePlayerHealth` use string formatting for serialization
- **Location**: `pkg/game/game.go:515-529`
- **Category**: Code Quality
- **Description**: Player health is encoded as `fmt.Sprintf("%.2f", health)` and decoded with `fmt.Sscanf`. This is fragile and loses precision.
- **Impact**: Potential precision loss on save/load. Extremely unlikely to affect gameplay.
- **Suggested Fix**: Use `strconv.FormatFloat` / `strconv.ParseFloat` for locale-independent serialization.

## Performance Optimization Opportunities

### [P-001] `screen.Set()` per-pixel particle rendering
- **Location**: `pkg/game/game.go:772-783`
- **Current Impact**: Each explosion (20 particles x ~4-9 pixels) generates 80-180 individual `Set()` calls.
- **Optimization**: Pre-generate small `ebiten.Image` particles and render with `DrawImage()` + `ColorScale`.
- **Expected Improvement**: 10-50x faster particle rendering, ~2-5ms saved per frame during heavy combat.

### [P-002] SFX synthesis on every playback
- **Location**: `pkg/audio/audio.go:156-178,311-329`
- **Current Impact**: Per-SFX CPU cost: laser ~441 samples x math ops, explosion ~13,230 samples x math ops.
- **Optimization**: Cache generated PCM data in a `sync.Once`-initialized map.
- **Expected Improvement**: Eliminate ~0.5-2ms per SFX playback after first generation.

### [P-003] `ProjectileCount()` and `CountEnemies()` iterate all entities
- **Location**: `pkg/combat/projectile.go:234-243`, `pkg/procgen/spawner.go:245-254`
- **Current Impact**: O(n) scan of all entities to count specific types.
- **Optimization**: Maintain running counters incremented on spawn and decremented on removal.
- **Expected Improvement**: O(1) count lookups instead of O(n) iteration.

### [P-004] `ForEachEntity` iterates the entire entity map for every system
- **Location**: `pkg/engine/engine.go:66-70`
- **Current Impact**: With 5+ systems each iterating all entities, total iterations = `entities x systems`.
- **Optimization**: Implement component queries -- index entities by component presence.
- **Expected Improvement**: Significant reduction in iteration overhead as entity count grows.

### [P-005] `GetParticles()` copies entire particle slice every frame
- **Location**: `pkg/rendering/rendering.go:297-304`
- **Current Impact**: Allocates `len(particles) * sizeof(Particle)` bytes per frame.
- **Optimization**: Return a read-only view or use double-buffering.
- **Expected Improvement**: Eliminate per-frame allocation of particle copy.

### [P-006] `drawGameplay` iterates entities twice
- **Location**: `pkg/game/game.go:715-736`
- **Current Impact**: All entities iterated by `CreateDrawBatches()`, then again for unbatched projectiles.
- **Optimization**: Handle all entity rendering in a single pass.
- **Expected Improvement**: Eliminate one full entity iteration per frame.

### [P-007] Collision detection is O(n^2)
- **Location**: `pkg/combat/projectile.go:106-118`
- **Current Impact**: Each projectile iterates all entities. With P projectiles and E entities, cost is O(P x E).
- **Optimization**: Use spatial partitioning (grid, quadtree) or maintain separate collidable entity lists.
- **Expected Improvement**: 5-20x reduction in collision checks for typical gameplay.

### [P-008] `CullContext.ShouldRender` allocates a new `Viewport` struct per call
- **Location**: `pkg/rendering/culling.go:56-63`
- **Current Impact**: Every entity render check creates a new `Viewport` on the heap.
- **Optimization**: Pre-compute the expanded viewport once per frame in the `CullContext` constructor.
- **Expected Improvement**: Eliminate 100+ small allocations per frame.

## Code Quality Observations

### [Q-001] Magic numbers in weapon/projectile spawning
- **Location**: `pkg/combat/weapon_system.go:112,117-118`
- **Issue**: `offset := 12.0`, `projectileSpeed := 400.0`, `projectileLifetime := 2.0` are hardcoded local variables without constants.
- **Suggestion**: Define named constants or make them configurable via `WeaponSystem` fields.

### [Q-002] Inconsistent system update patterns
- **Location**: `pkg/game/game.go:593-641`
- **Issue**: Some systems are registered with `world.AddSystem()` and updated via `world.Update(dt)`, while others are called manually. Mixed pattern makes system ordering hard to understand.
- **Suggestion**: Standardize on one approach and document it.

### [Q-003] `Game` struct has too many responsibilities
- **Location**: `pkg/game/game.go:117-161`
- **Issue**: The `Game` struct holds 20+ fields spanning rendering, audio, combat, AI, UI, save/load, tutorial, and sprite caching.
- **Suggestion**: Extract subsystem coordinators (e.g., `CombatManager`, `RenderManager`).

### [Q-004] Config controls bindings defined but never used
- **Location**: `pkg/config/config.go:43-50`, `pkg/game/game.go:246`
- **Issue**: `ControlsConfig` is parsed from `config.yaml` but `initializeSystems()` uses `engine.DefaultKeyBindings()` instead. Custom key bindings are silently ignored.
- **Suggestion**: Construct `KeyBindings` from `cfg.Controls` fields.

### [Q-005] No graceful error handling for component type assertions
- **Location**: Throughout -- e.g., `pkg/engine/arena.go:50`, `pkg/combat/projectile.go:67`
- **Issue**: All component retrievals use unchecked type assertions like `posComp.(*Position)`. Wrong types cause panics.
- **Suggestion**: Use comma-ok pattern or add typed accessor methods to World.

### [Q-006] Several TODO comments for v5.0 stubs
- **Location**: `pkg/companion/companion.go:3`, `pkg/networking/networking.go:4`, `pkg/hostplay/hostplay.go:3`, `pkg/security/security.go:3`, `pkg/social/social.go:3`
- **Issue**: These packages have `TODO(v5.0)` stubs. They return success for operations that do nothing.
- **Suggestion**: Ensure v5.0 stub packages are not accidentally wired into runtime code paths.

## Recommendations by Priority

1. **Immediate Action Required**
   - [C-001]: Move `prevKeys` from package-level to `Game` struct to prevent data races
   - [C-002]: Add edge detection for pause input to prevent rapid toggling
   - [C-003]: Fix "Retry" action to work from `StateGameOver` state
   - [C-004]: Eliminate double position update for projectiles (2x speed bug)

2. **High Priority (Next Sprint)**
   - [H-001]: Cache SFX PCM data to prevent regeneration on every playback
   - [H-002]: Make drag coefficient time-based instead of frame-based
   - [H-003]: Fix BoundingBox offset convention for consistent collision detection
   - [H-005]: Fix enemy spawn interaction with arena wrap mode
   - [H-008]: Replace inaccurate Taylor series with `math.Sin`/`math.Cos`

3. **Medium Priority (Backlog)**
   - [M-001]: Fix combo display to match actual scoring multiplier
   - [M-002]: Replace per-pixel `Set()` with batched particle rendering
   - [M-005]: Auto-generate seed when default seed is 0
   - [M-007]: Pass config screen dimensions to `EbitenInputReader`
   - [Q-004]: Wire config key bindings to input system

4. **Technical Debt**
   - [Q-001]: Extract magic numbers into named constants
   - [Q-002]: Standardize system update patterns (world-registered vs manual)
   - [Q-003]: Decompose `Game` struct into focused subsystems
   - [Q-005]: Add safe type assertion patterns for component access
   - [P-004]: Implement component-based entity queries
   - [P-007]: Add spatial partitioning for collision detection

## Testing Recommendations

- **Pause Toggle Test**: Simulate holding Escape for 10 frames; verify exactly one pause transition occurs.
- **Projectile Speed Test**: Spawn a projectile, advance 1 second of frames, measure distance traveled; should equal `projectileSpeed` (400px), not 800px.
- **Retry Flow Test**: Transition to GameOver state, invoke Retry action; verify state transitions to Playing.
- **Collision Accuracy Test**: Place two entities at known positions with known bounding boxes; verify collision detection matches expected overlap.
- **Combo Scoring Test**: Kill 6 enemies in sequence (combo=6); verify score multiplier is `1 + 6/5 = 2`, not 7.
- **Wrap Mode Spawn Test**: Spawn enemy at `y = -50` with wrap mode; verify it appears at expected position.
- **Performance Benchmark**: Spawn 100 particles, render for 60 frames; measure frame time with and without particle optimization.
- **Input Edge Cases**: Configure an invalid key name; verify graceful fallback behavior.
- **Long Session Test**: Run for 10,000 frames; verify no unbounded memory growth in audit logger or particle system.
- **Save/Load Round-Trip**: Save game at wave 5 with score 5000 and health 75.5; load and verify all values match.

## Audit Methodology Notes

- **Approach**: Static code analysis of all `.go` source files in the repository. Every file in `cmd/`, `pkg/`, and `mods/` was reviewed line-by-line.
- **Tools Used**: Manual code inspection, cross-referencing between files to trace call chains, data flow analysis for state management.
- **Areas Not Covered**: Runtime performance profiling (no execution), GPU/shader behavior, actual gamepad/touch device testing, WASM-specific behavior, network latency simulation (v5.0 stubs only).
- **Assumptions**: Analysis assumes Ebitengine v2.9.8 semantics. Fixed TPS of 60 assumed per Ebitengine defaults. Single-threaded game loop assumed (Ebitengine calls Update/Draw on main thread).
- **Limitations**: Static analysis cannot detect all runtime race conditions, timing-dependent bugs, or hardware-specific rendering issues. Some issues marked "Potential" indicate analysis-based suspicion rather than confirmed runtime observation.

## Positive Observations

- **Clean ECS Architecture**: The `engine.World` implementation is simple, correct, and well-documented. The Entity-Component-System pattern is consistently applied.
- **Genre System Design**: The five-genre system with `GenreSetter` interface and `genre.GetPreset()` is well-architected. Color palettes are distinct and genre-appropriate.
- **Sprite Generation**: The symmetric pixel-mirroring algorithm in `sprites.go` is elegant and produces visually distinctive results. The `shouldFillPixel` center-bias heuristic is a nice touch.
- **Build Tag Strategy**: The `noebiten` / `!noebiten` build tag pattern for CI testing is well-executed. Stub implementations match the real interface.
- **Audio Synthesis**: The procedural SFX generation (laser, explosion, powerup, menu) using basic waveforms and envelopes is well-implemented. The spatial audio system with volume attenuation and stereo panning is a solid foundation.
- **Save/Load**: Simple, correct JSON serialization with version field for future migration. Error handling is proper.
- **Deterministic RNG**: Consistent use of `engine.DeterministicRNG(seed)` throughout generation code. Seed derivation pattern (`seed + waveNumber`) is simple and effective.
- **Particle System**: Pre-allocation with capacity hints, mutex protection, dead-particle compaction using the `alive = ps.particles[:0]` pattern -- good performance-aware design.
- **Edge Detection for Menu Input**: The `prevKeys` + `wasKeyPressed()` pattern for menu navigation correctly prevents key repeat, showing awareness of the input edge-detection problem (though it needs to be extended to pause).
- **Viewport Culling**: The `CullContext` with margin-expanded viewport is a correct approach to preventing pop-in artifacts.

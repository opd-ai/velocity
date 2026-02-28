# Implementation Plan: v1.0 — Core Engine + Playable Single-Player

## Phase Overview
- **Objective**: Deliver a playable single-player arcade shooter with the `scifi` genre, featuring ship physics, basic combat, procedural waves, and a complete game-loop from main menu to game-over.
- **Prerequisites**: Go 1.24+, Ebitengine v2.9+, Viper v1.21+ (all present in go.mod)
- **Estimated Scope**: Large

## Procedural Generation Requirement

**All gameplay assets must be 100% procedurally generated at runtime using deterministic algorithms.** This applies to:

- **Audio**: Music layers, sound effects, and positional audio generated via waveform synthesis (sine, square, noise) — no audio files
- **Visual**: Ship sprites, enemy sprites, projectiles, particles, UI elements generated via pixel algorithms and drawing primitives — no image files  
- **Narrative**: Quest text, objective descriptions, environmental lore generated algorithmically from seed values — no embedded text assets or dialogue files

The game must produce identical output given identical seed inputs, with zero reliance on external, bundled, or pre-authored media files.

## Implementation Steps

### Ship Physics & Player Control

1. Implement ship physics component
   - Deliverable: `pkg/engine/physics.go` — `Position`, `Velocity`, `Rotation` components and a `PhysicsSystem` that applies thrust, rotation, inertia, and drag each tick using Newtonian 2-D flight model
   - Dependencies: ECS core (`pkg/engine/engine.go`)

2. Wire input system to player entity
   - Deliverable: `pkg/engine/input.go` — `InputSystem` that reads keyboard state via Ebitengine's `ebiten.IsKeyPressed`, maps it to `InputState`, and applies thrust/rotation/fire commands to the player entity
   - Dependencies: Step 1 (physics component), `pkg/config` (control bindings)

3. Implement screen-wrap and bounded-arena modes
   - Deliverable: `pkg/engine/arena.go` — `ArenaSystem` that checks entity positions against screen bounds; wraps coordinates in "wrap" mode or applies bounce/clamp in "bounded" mode, configurable via `config.yaml` `arena_mode`
   - Dependencies: Step 1 (physics component), `pkg/config`

### Rendering Pipeline

4. Implement procedural sprite generation for ships and enemies
   - Deliverable: `pkg/rendering/sprites.go` — functions to procedurally generate ship and enemy sprite images as `*ebiten.Image` from seed values using pixel algorithms (no external image files), keyed by genre + variant, cached in `SpriteCache`
   - Dependencies: `pkg/rendering/rendering.go` (SpriteCache), `pkg/procgen/genre`

5. Implement sprite drawing in the render loop
   - Deliverable: Update `pkg/rendering/rendering.go` — `Renderer.Draw(screen, world, camera)` method that iterates entities with Position + Sprite components, applies camera offset and rotation, and draws via `ebiten.DrawImageOptions`
   - Dependencies: Step 4 (sprite generation), Step 1 (position components)

6. Implement particle system update and draw
   - Deliverable: Complete `ParticleSystem.Emit()` and `ParticleSystem.Update()` stubs in `pkg/rendering/rendering.go`; add `ParticleSystem.Draw(screen)` method for thruster trails, explosion debris, and impact sparks
   - Dependencies: Step 5 (render loop)

7. Implement viewport culling and draw batching
   - Deliverable: `pkg/rendering/culling.go` — skip draw calls for entities outside the viewport; batch draw calls per entity type per frame to reduce GPU overhead
   - Dependencies: Step 5 (render loop), camera system

### Combat & Collision

8. Implement hit detection and projectile system
   - Deliverable: `pkg/combat/projectile.go` — `ProjectileSystem` that spawns, moves, and despawns projectile entities; axis-aligned bounding-box (AABB) collision detection between projectiles and ships/enemies
   - Dependencies: Step 1 (physics), ECS core

9. Integrate weapon firing with player input
   - Deliverable: Update `pkg/combat/combat.go` — wire `Weapon.Fire()` to spawn a projectile entity when the player presses fire; respect cooldown timer
   - Dependencies: Step 2 (input), Step 8 (projectiles)

10. Implement damage calculation and entity destruction
    - Deliverable: `pkg/combat/damage.go` — `DamageSystem` that applies hit-points reduction on collision, triggers entity removal on death, and emits particle effects (explosions) on destruction
    - Dependencies: Step 8 (collision), Step 6 (particles)

### Procedural Wave Generation

11. Implement wave spawning system
    - Deliverable: `pkg/procgen/spawner.go` — `WaveSpawner` that uses `Generator.GenerateWave()` to determine enemy count, then creates enemy entities at randomised off-screen positions with basic linear-approach AI
    - Dependencies: `pkg/procgen/procgen.go`, ECS core, Step 4 (enemy sprites)

12. Implement wave progression and difficulty ramp
    - Deliverable: Update `pkg/procgen/procgen.go` — enhance `GenerateWave` to increase enemy speed, health, and spawn rate per wave; add a wave-complete check (all enemies destroyed) to trigger next wave
    - Dependencies: Step 11 (spawning), Step 10 (enemy destruction)

### Audio

13. Implement adaptive music playback
    - Deliverable: Update `pkg/audio/audio.go` — procedurally generate music tones via waveform synthesis (no audio files); implement intensity-driven layering that increases with wave number using algorithmically generated audio streams
    - Dependencies: Ebitengine audio API

14. Implement SFX triggers
    - Deliverable: Update `pkg/audio/audio.go` — implement `PlaySFX()` for weapon fire, explosions, and powerup collect using procedurally generated waveforms (sine/square/noise bursts) with no audio files
    - Dependencies: Step 13 (audio manager initialisation)

15. Implement positional audio cues
    - Deliverable: Update `pkg/audio/audio.go` — implement volume/pan adjustment based on entity distance and angle from player for off-screen enemy audio cues
    - Dependencies: Step 14 (SFX), Step 1 (positions)

### UI / HUD / Menus

16. Implement HUD rendering
    - Deliverable: `pkg/ux/hud_render.go` — draw health bar, shield meter, score display, wave counter, and combo multiplier on screen using `ebitenutil.DebugPrint` or procedural UI drawing
    - Dependencies: `pkg/ux/ux.go` (HUD state), player entity health, score state

17. Implement menu flow and game state management
    - Deliverable: `pkg/ux/menu_render.go` — render main menu, pause menu, and game-over screen; handle state transitions (MainMenu → Playing → Paused → GameOver); integrate with `cmd/velocity/main.go` game loop
    - Dependencies: `pkg/ux/ux.go` (Menu state), input system

18. Implement first-run tutorial
    - Deliverable: `pkg/ux/tutorial_render.go` — guided first wave teaching thrust, fire, and dodge with overlay text prompts; triggered on first run
    - Dependencies: Step 16 (HUD), Step 17 (menu flow), Step 2 (input)

### Game Loop Integration

19. Wire all systems into the main game loop
    - Deliverable: Update `cmd/velocity/main.go` — integrate `PhysicsSystem`, `InputSystem`, `ArenaSystem`, `ProjectileSystem`, `DamageSystem`, `WaveSpawner` into `Game.Update()`; integrate `Renderer.Draw()`, `ParticleSystem.Draw()`, HUD and Menu rendering into `Game.Draw()`
    - Dependencies: All steps above

20. Implement save/load integration
    - Deliverable: Update `cmd/velocity/main.go` — save `RunState` on pause/exit via `pkg/saveload`; load and resume on startup if save file exists
    - Dependencies: `pkg/saveload/saveload.go`, Step 17 (menu flow)

### Infrastructure & Quality

21. Add unit tests for core systems
    - Deliverable: `*_test.go` files for `pkg/engine`, `pkg/combat`, `pkg/procgen`, `pkg/rendering`, `pkg/ux`, `pkg/config`, `pkg/saveload`, `pkg/errors`, `pkg/validation`, `pkg/version` — targeting 82%+ coverage
    - Dependencies: All implementation steps above

22. Set up CI/CD pipeline
    - Deliverable: `.github/workflows/ci.yml` — GitHub Actions workflow with `GOOS` matrix (linux, darwin, windows) running `go build ./...`, `go test ./...`, and `go vet ./...`
    - Dependencies: None (can be done early)

23. Add cross-platform build targets
    - Deliverable: Update CI workflow or add `Makefile` — build single binary for Linux, macOS, Windows, and WASM (`GOOS=js GOARCH=wasm`)
    - Dependencies: Step 22 (CI)

24. Write project documentation
    - Deliverable: `CONTROLS.md` (keybindings reference), `CHANGELOG.md` (v1.0 release notes), `FAQ.md` (common questions)
    - Dependencies: Step 19 (finalised controls and gameplay)

## Technical Specifications

- **Physics model**: Newtonian 2-D — thrust applies acceleration along ship facing vector; velocity accumulates with configurable drag coefficient (0.98 per tick default); rotation is angular velocity with instant response
- **Collision detection**: AABB (axis-aligned bounding box) for v1.0; sufficient for top-down sprite rectangles; upgrade to circle-based or SAT in v2.0+ if needed
- **Sprite generation**: Procedural pixel-art via Go `image` package rendered to `*ebiten.Image`; symmetric ship hulls generated from seed using mirrored random pixel placement (no external image files); cached in `SpriteCache` by `genreID:variant` key
- **Audio generation**: Procedural waveforms using Ebitengine's `audio` package with no audio files; sine waves for music tones, noise bursts for explosions, square waves for laser SFX — all generated algorithmically at runtime
- **Arena bounds**: Screen dimensions from `config.yaml` (`800×600` default); wrap mode teleports entities crossing an edge to the opposite side; bounded mode reverses velocity component on contact
- **Wave formula**: Wave N spawns `N + 2` enemies with health `10 + N*5` and speed `1.0 + N*0.1`; single enemy type for v1.0 (basic drone)
- **Frame rate**: Fixed 60 TPS via Ebitengine default; `dt = 1.0/60.0` used for physics integration
- **Entity budget**: Target ≤200 entities on screen for v1.0 (player + enemies + projectiles + particles)

## Validation Criteria
- [ ] Game launches to main menu with "Start", "Settings", and "Quit" options
- [ ] Player ship spawns at centre and responds to thrust (W), rotate (A/D), and fire (Space)
- [ ] Ship exhibits Newtonian physics: inertia, drift, and drag
- [ ] Screen-wrap mode wraps entities at screen edges; bounded mode bounces them
- [ ] Enemies spawn in waves; wave counter increments when all enemies are destroyed
- [ ] Projectiles collide with enemies and destroy them; explosions produce particles
- [ ] HUD displays health, score, wave number, and combo multiplier
- [ ] Pause menu accessible via Escape; resume and quit-to-menu options work
- [ ] Game-over screen appears when player health reaches zero; shows final score
- [ ] Save/load persists and resumes a run across application restarts
- [ ] Procedural sprites render correctly for the `scifi` genre
- [ ] Audio plays: background music, weapon-fire SFX, explosion SFX
- [ ] All packages have unit tests; `go test ./...` passes with ≥82% coverage
- [ ] CI pipeline builds and tests on Linux, macOS, and Windows
- [ ] Single binary with zero external asset files (all assets procedurally generated)

## Known Gaps
- **Procedural audio generation**: Ebitengine's audio API supports PCM streaming, but the exact approach for generating music layers procedurally (e.g., frequency tables, beat sequencing) is not specified in the roadmap — see GAPS.md
- **Sprite generation algorithm**: The roadmap requires "procedurally drawn" sprites with no external assets, but the specific algorithm (pixel-mirroring, shape primitives, L-system) is not defined — see GAPS.md
- **Tutorial trigger mechanism**: How to detect "first run" vs returning player is unspecified; save-file existence is the likely heuristic but needs confirmation
- **Gamepad and touch input**: Roadmap lists keyboard/gamepad/touch support in v1.0, but the Ebitengine gamepad and touch APIs require platform-specific handling that is not detailed

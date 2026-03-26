# Implementation Gaps — 2026-03-26

This document tracks gaps between the project's stated goals (README.md, ROADMAP.md) and the current implementation.

---

## Audio Playback Not Connected to Ebitengine

- **Stated Goal**: "Adaptive music with intensity-driven music layers; genre-specific instrumentation" and "SFX: weapon fire, explosions, powerup collect, UI feedback" (ROADMAP.md v1.0 Tier 1).
- **Current State**: `pkg/audio/audio.go` contains complete procedural synthesis functions (`GenerateTone`, `GenerateLaserSFX`, `GenerateExplosionSFX`, `GeneratePowerupSFX`, `GenerateMenuSelectSFX`) that produce valid 16-bit stereo PCM data. However, `Manager.Update()` (line 107-113) only clears the SFX queue without ever playing audio. No `ebiten/v2/audio.Context` is initialized. `PlayMusic()` only sets a boolean flag.
- **Impact**: The game runs in complete silence. Players receive no audio feedback for actions, reducing gameplay impact and genre immersion.
- **Closing the Gap**:
  1. Add `audioContext *audio.Context` field to `Manager` struct
  2. Initialize in `NewManager()`: `m.audioContext = audio.NewContext(sampleRate)`
  3. In `Update()`, for each queued SFX, create a player: `player := m.audioContext.NewPlayerFromBytes(GetSFXData(name))`
  4. Call `player.Play()` (or use a player pool for efficiency)
  5. For adaptive music, implement a looping background track generator that changes based on `intensity` field
  6. Test with: visual confirmation of audio playback; unit test that verifies audio context initialization

---

## Procedural Sprites Not Rendered

- **Stated Goal**: "Procedurally drawn ship/enemy/projectile sprites (no external assets)" (ROADMAP.md v1.0 Tier 1).
- **Current State**: `pkg/rendering/sprites.go` implements complete symmetric pixel generation with genre color palettes. Sprites are generated and cached in `SpriteCache`. Entities have `SpriteComponent` attached (`cmd/velocity/main.go:266-270`). However, `drawGameplay()` (lines 555-585) ignores sprites entirely and renders all entities as monochrome 8×8 rectangles using `ebitenutil.DrawRect()`.
- **Impact**: The game has no visual variety. All ships, enemies, and projectiles look identical. Genre-specific visual theming is invisible to players.
- **Closing the Gap**:
  1. In `drawGameplay()`, retrieve `sprite` component: `spriteComp, ok := g.world.GetComponent(e, "sprite")`
  2. Based on `SpriteType`, call appropriate `renderer.GetOrCreate*Sprite(variant, size)`
  3. Convert `*image.RGBA` to `*ebiten.Image` (cache this conversion)
  4. Apply rotation from `rotation` component using `ebiten.DrawImageOptions{}.GeoM.Rotate(rot.Angle)`
  5. Draw with `screen.DrawImage(img, opts)`
  6. Test with: visual confirmation that ships have distinct pixel patterns; screenshot comparison across genres

---

## Gamepad and Touch Input Not Implemented

- **Stated Goal**: "Keyboard/gamepad/touch; rebindable controls" (ROADMAP.md v1.0 Tier 1). CONTROLS.md explicitly documents "Gamepad Support (Planned)".
- **Current State**: `pkg/engine/input_ebiten.go` only implements keyboard input via `ebiten.IsKeyPressed()`. The `InputReader` interface exists but no gamepad or touch implementation is provided. `config.yaml` only defines keyboard bindings.
- **Impact**: Players without keyboards (mobile, console, Steam Deck) cannot play the game despite this being a v1.0 deliverable per the roadmap.
- **Closing the Gap**:
  1. Extend `EbitenInputReader.ReadState()` to poll gamepads:
     - `ids := ebiten.GamepadIDs(); for _, id := range ids { ... }`
     - Map `ebiten.IsGamepadButtonPressed(id, button)` to `InputState` fields
  2. Add touch support:
     - `touchIDs := ebiten.TouchIDs()`
     - Define screen regions for virtual buttons
  3. Extend `config.yaml` with `gamepad:` and `touch:` sections
  4. Test with: connected gamepad input; touch simulator on WASM build

---

## Tutorial System Not Integrated

- **Stated Goal**: "First-run guided wave teaching thrust, fire, and dodge" (ROADMAP.md v1.0 Tier 1).
- **Current State**: `pkg/ux/ux.go:60-79` defines a `Tutorial` struct with `Active`, `Step`, `Advance()`, and `Complete()` methods. However, `Tutorial` is never instantiated in `cmd/velocity/main.go`. No first-run detection exists. No tutorial UI is rendered.
- **Impact**: New players receive no onboarding. The game assumes familiarity with thrust-based space games, creating a steep learning curve.
- **Closing the Gap**:
  1. Add first-run detection: check if `savePath()` exists; if not, this is first run
  2. Create `Tutorial` instance in `NewGame()` when first run detected
  3. During first wave, render tutorial overlays: "Press W to THRUST", "Press A/D to ROTATE", "Press SPACE to FIRE"
  4. Advance tutorial steps as player demonstrates each action
  5. Mark tutorial complete after first enemy kill; persist marker to prevent re-triggering
  6. Test with: delete save file, verify tutorial appears; complete tutorial, restart, verify tutorial does not appear

---

## Adaptive Music System

- **Stated Goal**: "Intensity-driven music layers" that change based on combat state (ROADMAP.md).
- **Current State**: `Manager.intensity` field exists and `SetIntensity()` is implemented, but nothing reads `intensity` to affect audio output. `PlayMusic()` sets `musicPlaying = true` but generates no audio. `GetGenreParams()` returns tempo and waveform parameters per genre but these are unused.
- **Impact**: No background music plays. The intensity-driven layering system is conceptualized but not implemented.
- **Closing the Gap**:
  1. Implement a music generator that produces continuous PCM based on `GenreAudioParams`
  2. Use `intensity` to control layer mixing: low intensity = ambient pad only; high intensity = full percussion
  3. Stream music via `audio.InfiniteLoop` or a custom `io.Reader` implementation
  4. Call `SetIntensity()` from `updateGameplay()` based on enemy count or combat state
  5. Test with: audio playback that audibly changes when waves start/end

---

## Tests Fail in Headless Environments

- **Stated Goal**: "82%+ test coverage" and "CI/CD: Multi-platform GitHub Actions" (ROADMAP.md v1.0).
- **Current State**: Tests in `pkg/combat`, `pkg/engine`, `pkg/procgen`, `pkg/rendering` panic with "GLFW library is not initialized" when `DISPLAY` environment variable is unset. This occurs because test files import Ebitengine, which initializes GLFW at package init time.
- **Impact**: CI pipelines in headless environments cannot run tests. Coverage metrics cannot be collected. The 82% target cannot be verified.
- **Closing the Gap**:
  1. Refactor tests to use interface abstractions instead of concrete Ebiten types where possible
  2. For tests requiring Ebiten, add build tags: `//go:build !headless`
  3. Create mock implementations of `InputReader`, image types for unit tests
  4. In CI, use `xvfb-run -a go test ./...` to provide a virtual display
  5. Alternatively, use `GOOS=js GOARCH=wasm go test` (WASM tests don't need display)
  6. Test with: `go test ./...` succeeds in headless Docker container

---

## Viewport Culling Not Applied

- **Stated Goal**: "Skip update/draw for off-screen entities" (ROADMAP.md v1.0 — Performance).
- **Current State**: `pkg/rendering/culling.go` implements `IsOnScreen()` function. However, `cmd/velocity/main.go:drawGameplay()` iterates all entities via `world.ForEachEntity()` without any culling check.
- **Impact**: When entity count increases (later waves, many projectiles), all entities are processed for drawing regardless of visibility. This will cause frame drops at scale.
- **Closing the Gap**:
  1. In `drawGameplay()`, retrieve position and check `rendering.IsOnScreen(x, y, entitySize, screenWidth, screenHeight)`
  2. Skip `DrawRect` / sprite rendering for entities outside viewport
  3. Apply same culling in physics update for entities far outside bounds
  4. Test with: spawn 500 off-screen entities, verify FPS remains at 60

---

## Draw Batching Not Utilized

- **Stated Goal**: "Batch draw calls per entity type per frame" (ROADMAP.md v1.0 — Performance).
- **Current State**: `pkg/rendering/rendering.go:304-325` implements `CreateDrawBatches()` that groups entities by `SpriteType`. This function is never called from `cmd/velocity/main.go`.
- **Impact**: Each entity is drawn with a separate draw call. Modern graphics benefit from batched drawing; current approach will bottleneck at high entity counts.
- **Closing the Gap**:
  1. Call `batches := rendering.CreateDrawBatches(g.world)` at start of `drawGameplay()`
  2. Iterate batches; for each batch, get/cache the sprite image once
  3. Draw all entities in batch using the same sprite with different positions/rotations
  4. Test with: profile draw calls before/after; verify reduced GPU command count

---

## Config Validation Not Enforced

- **Stated Goal**: "Input validation for config, save data, network messages" (ROADMAP.md v1.0).
- **Current State**: `pkg/validation/validation.go` provides `ValidateGenre()` and `ValidateArenaMode()`. `pkg/config/config.go:Load()` does not call these validators. Invalid genre strings would silently use default presets.
- **Impact**: User configuration errors are not caught at startup. A typo in `config.yaml` (e.g., `genre: "scify"`) would not produce an error, leading to unexpected behavior.
- **Closing the Gap**:
  1. After `viper.Unmarshal(&cfg)` in `Load()`, add validation:
     ```go
     if err := validation.ValidateGenre(cfg.Gameplay.Genre); err != nil {
         return nil, err
     }
     if err := validation.ValidateArenaMode(cfg.Gameplay.ArenaMode); err != nil {
         return nil, err
     }
     ```
  2. Test with: set invalid genre in config.yaml, verify startup error

---

## v5.0+ Features Are Stubs

- **Stated Goal**: v5.0+ features (multiplayer, security, social) are documented in ROADMAP.md as future milestones.
- **Current State**: The following packages contain only stub implementations:
  - `pkg/networking/networking.go` — Server/Client with Start()/Connect() that set booleans
  - `pkg/security/security.go` — Encrypt/Decrypt return `ErrNotImplemented`
  - `pkg/social/social.go` — Squadron/Leaderboard with in-memory-only storage
  - `pkg/hostplay/hostplay.go` — Host with no actual networking
  - `pkg/companion/companion.go` — Wingman.Update() is empty
- **Impact**: These are correctly scoped to v5.0+ and do not affect v1.0 functionality. However, their presence in the codebase without clear documentation could confuse contributors.
- **Closing the Gap**:
  1. This is acceptable for v1.0 release scope
  2. Ensure ROADMAP.md clearly marks these as v5.0+ milestones
  3. Add prominent `// TODO(v5.0):` comments in stub implementations
  4. Consider moving to `internal/stub/` to signal non-production status

---

## Score Persistence and Leaderboards

- **Stated Goal**: "Per-seed and global high-score tracking" (ROADMAP.md v4.0).
- **Current State**: `RunState` includes `Score` which is saved/loaded. `pkg/social/social.go` has `Leaderboard` struct. However, high scores are not persisted between sessions, and no leaderboard UI exists.
- **Impact**: Players cannot track their best runs or compare with others. This is a v4.0 feature but basic local high-score tracking would enhance v1.0.
- **Closing the Gap** (v1.0 scope):
  1. Add `HighScore int64` to a persistent settings file
  2. On game over, compare score to high score; update if exceeded
  3. Display "HIGH SCORE: X" on main menu
  4. Full leaderboard system deferred to v4.0 per roadmap

---

## Magic Numbers Throughout Codebase

- **Stated Goal**: Maintainable, readable code following Go conventions.
- **Current State**: `go-stats-generator` reports 645 magic numbers. Examples:
  - `pkg/class/class.go:16-19` — Ship stats (60, 300, 2, 2)
  - `pkg/audio/audio.go:16` — Pentatonic frequencies
  - `cmd/velocity/main.go:146-149` — Physics tuning (200.0, 4.0, 0.98, 300.0)
- **Impact**: Tuning values are scattered throughout code, making balance adjustments difficult.
- **Closing the Gap**:
  1. Extract constants to package-level vars: `const ScoutHealth = 60.0`
  2. For tunable values, add to `config.yaml` schema
  3. Target: reduce magic numbers by 50% in v1.1
  4. Test with: `go-stats-generator analyze . --format json | jq '.maintenance.magic_numbers'` shows <350

---

## Main Entry Point Oversized

- **Stated Goal**: Clean, maintainable architecture following Go project layout conventions.
- **Current State**: `cmd/velocity/main.go` has 475 lines with 28 functions. It contains game logic, rendering, input handling, state management, and menu rendering.
- **Impact**: Difficult to navigate, test, and maintain. Violates single-responsibility principle.
- **Closing the Gap**:
  1. Extract `Game` struct and methods to `pkg/game/game.go`
  2. Keep `cmd/velocity/main.go` to ~20 lines: config load, game init, run
  3. Extract `drawGameplay()`, `drawHUD()`, `drawMenu()` to `pkg/rendering/` or `pkg/ux/`
  4. Test with: `wc -l cmd/velocity/main.go` shows <100 lines

---

*This document is updated as gaps are identified or closed. Cross-reference with AUDIT.md findings.*

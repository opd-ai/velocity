# Implementation Gaps — 2026-03-26

This document tracks gaps between the project's stated goals (README.md, ROADMAP.md) and the current implementation.

---

## Adaptive Background Music Not Playing

- **Stated Goal**: "Intensity-driven music layers; genre-specific instrumentation" (ROADMAP.md v1.0 Tier 1).
- **Current State**: `pkg/audio/audio.go:123-125` — `PlayMusic()` sets `musicPlaying = true` but no audio is generated. `SetIntensity()` exists at line 91-93 but is never called from game code. `GetGenreParams()` returns tempo and waveform parameters per genre (lines 433-478) but these are unused.
- **Impact**: The game has no background music. Players experience silence except for SFX, reducing atmosphere and genre immersion.
- **Closing the Gap**:
  1. Implement `generateMusicSamples()` in `Manager.Update()` that produces continuous PCM based on `GenreAudioParams`
  2. Use `intensity` field to control layer mixing: low = ambient pad only; high = full percussion
  3. Stream music via `audio.InfiniteLoop` or a custom `io.Reader`
  4. Call `audio.SetIntensity(0.3)` in `updateGameplay()` when idle, `audio.SetIntensity(0.8)` when `waveManager.WaveInProgress()`
  5. Call `audio.PlayMusic()` in `startNewGame()`
  6. Test with: audio playback that audibly changes when waves start/end

---

## Touch Input Not Implemented

- **Stated Goal**: "Keyboard/gamepad/touch; rebindable controls" (ROADMAP.md v1.0 Tier 1).
- **Current State**: `pkg/engine/input_ebiten.go` implements keyboard (lines 23-31) and gamepad (lines 40-79) input. No `ebiten.TouchIDs()` calls exist. CONTROLS.md documents "Gamepad Support (Planned)" but touch is not mentioned.
- **Impact**: Players on mobile devices (iOS, Android via ebitenmobile) and touch-enabled laptops cannot play the game, limiting the reach of WASM builds.
- **Closing the Gap**:
  1. Add touch handling to `EbitenInputReader.ReadState()`:
     ```go
     touchIDs := ebiten.TouchIDs()
     for _, id := range touchIDs {
         x, y := ebiten.TouchPosition(id)
         // Map screen regions to virtual buttons
     }
     ```
  2. Define touch control regions (thrust zone, fire zone, rotate zones)
  3. Add `touch:` section to `config.yaml` for region customization
  4. Update CONTROLS.md with touch control documentation
  5. Test with: WASM build on mobile browser, verify thrust/rotate/fire work

---

## ParticleSystem Not Integrated

- **Stated Goal**: "Explosion debris, thruster trails, impact sparks" (ROADMAP.md v1.0 Tier 1 — Rendering particles).
- **Current State**: `pkg/rendering/rendering.go:148-296` — Complete `ParticleSystem` implementation exists with `Emit()`, `EmitDirectional()`, `Update()`, `GetParticles()`, genre-specific colors. However, it is never instantiated in `cmd/velocity/main.go`. No particles appear during gameplay.
- **Impact**: Explosions and thruster effects are invisible. Visual feedback for destruction and movement is absent, making gameplay feel flat.
- **Closing the Gap**:
  1. Add `particleSystem *rendering.ParticleSystem` field to `Game` struct
  2. Initialize in `NewGame()`: `g.particleSystem = rendering.NewParticleSystem()`
  3. Set genre: `g.particleSystem.SetGenre(genre)`
  4. Wire to death callback in `initializeSystems()`:
     ```go
     g.damageSystem.SetDeathCallback(func(event combat.DeathEvent) {
         if pos, ok := g.world.GetComponent(event.Entity, "position"); ok {
             p := pos.(*engine.Position)
             g.particleSystem.Emit(p.X, p.Y, 20)
         }
         // ... existing scoring logic
     })
     ```
  5. Call `g.particleSystem.Update(dt)` in `updateGameplay()`
  6. Render particles in `drawGameplay()` after entity rendering
  7. Test with: destroy enemy, verify particle burst appears

---

## Tests Fail in Headless Environments

- **Stated Goal**: "82%+ test coverage" and "CI/CD: Multi-platform GitHub Actions" (ROADMAP.md v1.0).
- **Current State**: Tests in `pkg/combat`, `pkg/engine`, `pkg/procgen`, `pkg/rendering` panic with "GLFW library is not initialized" when `DISPLAY` environment variable is unset. This occurs because test files import packages that transitively import `github.com/hajimehoshi/ebiten/v2`, which initializes GLFW at package init time.
- **Impact**: CI pipelines in headless environments (GitHub Actions, Docker) cannot run tests. Coverage cannot be measured or enforced. Test reliability is unverifiable.
- **Closing the Gap**:
  1. For immediate CI fix: use `xvfb-run -a go test ./...` to provide a virtual display
  2. Long-term: add build tag `//go:build !headless` to test files that require Ebiten graphics
  3. Create mock implementations of `InputReader` and image types for unit tests
  4. Alternatively, test pure-logic functions separately from Ebiten-dependent code
  5. Test with: `DISPLAY="" go test ./pkg/procgen/...` succeeds without panic (after splitting tests)

---

## No CI/CD Workflow

- **Stated Goal**: "CI/CD: Multi-platform GitHub Actions: Linux, macOS, Windows, WASM" (ROADMAP.md v1.0).
- **Current State**: No `.github/workflows/` directory exists. Build and test automation is not configured.
- **Impact**: Code quality cannot be automatically verified on push. Cross-platform compatibility is not tested. Regressions can be merged undetected.
- **Closing the Gap**:
  1. Create `.github/workflows/ci.yml`:
     ```yaml
     name: CI
     on: [push, pull_request]
     jobs:
       test:
         strategy:
           matrix:
             os: [ubuntu-latest, macos-latest, windows-latest]
         runs-on: ${{ matrix.os }}
         steps:
           - uses: actions/checkout@v4
           - uses: actions/setup-go@v5
             with:
               go-version: '1.24'
           - name: Test (Linux)
             if: runner.os == 'Linux'
             run: xvfb-run -a go test ./...
           - name: Test (other)
             if: runner.os != 'Linux'
             run: go test ./...
           - name: Build
             run: go build ./cmd/velocity/
       wasm:
         runs-on: ubuntu-latest
         steps:
           - uses: actions/checkout@v4
           - uses: actions/setup-go@v5
             with:
               go-version: '1.24'
           - run: GOOS=js GOARCH=wasm go build -o velocity.wasm ./cmd/velocity/
     ```
  2. Test with: push to branch, verify Actions run and pass

---

## Magic Numbers Throughout Codebase

- **Stated Goal**: Maintainable, readable code following Go conventions.
- **Current State**: `go-stats-generator` reports 648 magic numbers. Key examples:
  - `pkg/class/class.go:16-19` — Ship stats (60, 300, 2, 2)
  - `cmd/velocity/main.go:288-289` — Bounding box dimensions (-8, 16)
  - `pkg/audio/audio.go:16-27` — Frequency constants (already named)
  - `cmd/velocity/main.go:31-40` — Physics constants (already named)
- **Impact**: Tuning values scattered throughout code make balance adjustments difficult. Non-obvious numbers reduce readability.
- **Closing the Gap**:
  1. Review magic numbers in `pkg/class/class.go` — extract to `ShipClassConfig` with named fields
  2. For bounding box sizes, define `const DefaultEntitySize = 16`
  3. For gameplay tuning (wave formula coefficients), consider adding to `config.yaml`
  4. Target: reduce magic number count by 50% in v1.1
  5. Test with: `go-stats-generator analyze . --format json | jq '.maintenance.magic_numbers'` < 350

---

## Main Entry Point Oversized

- **Stated Goal**: Clean, maintainable architecture following Go project layout conventions.
- **Current State**: `cmd/velocity/main.go` has 820 lines with 34 functions. It contains:
  - Game struct and core loop
  - Physics integration
  - Rendering logic (`drawGameplay`, `drawHUD`, `drawMenu`, `drawTutorial`)
  - Input handling
  - Menu navigation
  - State management
  - Save/load orchestration
- **Impact**: Difficult to navigate, test in isolation, and maintain. Violates single-responsibility principle.
- **Closing the Gap**:
  1. Extract `Game` struct and methods to `pkg/game/game.go`
  2. Keep `cmd/velocity/main.go` to ~20 lines: config load, game init, `ebiten.RunGame()`
  3. Extract draw functions to `pkg/ux/draw.go` (can use Ebiten types)
  4. Extract menu input handling to `pkg/ux/input.go`
  5. Test with: `wc -l cmd/velocity/main.go` < 100

---

## v5.0+ Features Are Stubs

- **Stated Goal**: v5.0+ features (multiplayer, security, social, companion AI) are documented in ROADMAP.md as future milestones.
- **Current State**: The following packages contain only stub implementations:
  - `pkg/networking/networking.go` — Server/Client with Start()/Connect() that set booleans
  - `pkg/security/security.go` — Encrypt/Decrypt return `ErrNotImplemented`
  - `pkg/social/social.go` — Squadron/Leaderboard with in-memory-only storage
  - `pkg/hostplay/hostplay.go` — Host with no actual networking
  - `pkg/companion/companion.go` — Wingman.Update() is empty
- **Impact**: These are correctly scoped to v5.0+ per ROADMAP.md. Each has appropriate `// TODO(v5.0):` comments. Their presence does not affect v1.0 functionality.
- **Closing the Gap**:
  1. No action required for v1.0 release
  2. Ensure ROADMAP.md continues to clearly mark these as v5.0+ milestones
  3. Consider moving to `internal/stub/` to signal non-production status (optional)

---

## High Score Persistence

- **Stated Goal**: "Per-seed and global high-score tracking" (ROADMAP.md v4.0).
- **Current State**: `saveload.RunState` includes `Score` which is saved/loaded during gameplay sessions. `pkg/social/social.go` has `Leaderboard` struct. However, no persistent high score is tracked between separate game sessions. No high score display on main menu.
- **Impact**: Players cannot track their best runs across sessions. Replayability incentive is reduced.
- **Closing the Gap** (optional for v1.0, required for v4.0):
  1. Add `HighScores map[int64]int64` to a persistent settings file (keyed by seed)
  2. On game over, compare score to high score for that seed; update if exceeded
  3. Display "HIGH SCORE: X" on main menu and game over screen
  4. Full leaderboard system (network-backed) deferred to v4.0 per roadmap
  5. Test with: complete game, restart, verify high score persists

---

*This document is updated as gaps are identified or closed. Cross-reference with AUDIT.md findings.*

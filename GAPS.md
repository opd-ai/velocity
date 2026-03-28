# Implementation Gaps — 2026-03-28

This document tracks gaps between the project's stated goals (README.md, ROADMAP.md, FAQ.md) and the current implementation.

---

## Adaptive Background Music Not Generating Audio

- **Stated Goal**: "Adaptive music and SFX — intensity-driven music layers, spatial audio" (README.md pkg/audio description). "Intensity-driven music layers; genre-specific instrumentation" (ROADMAP.md v1.0).
- **Current State**: 
  - `pkg/audio/audio.go:126-128`: `PlayMusic()` sets `musicPlaying = true` but generates no audio.
  - `pkg/audio/audio.go:94-96`: `SetIntensity()` stores value but is never consumed by any music generator.
  - `pkg/audio/audio.go:434-478`: `GetGenreParams()` returns tempo, scale, and waveform parameters per genre but these are unused.
  - Game loop correctly calls `SetIntensity(0.3)` idle and `SetIntensity(0.8)` during combat at `pkg/game/game.go:624-628`.
- **Impact**: The game has no background music. Players experience silence except for SFX, reducing atmosphere and genre immersion. This is the primary missing v1.0 feature.
- **Closing the Gap**:
  1. Implement `generateMusicStream()` in `pkg/audio/audio.go` that returns an `io.Reader` producing continuous PCM samples.
  2. Use `GenreAudioParams` (tempo, waveform, scale) from `GetGenreParams()` to vary output per genre.
  3. Create 3-layer system:
     - Base layer (pad, always playing): Low-frequency sine wave at genre tempo.
     - Combat layer (intensity > 0.5): Add percussion (noise bursts at beat intervals).
     - Boss layer (intensity > 0.8): Add melody (pentatonic scale arpeggios).
  4. Stream via `audio.NewInfiniteLoop` or custom `io.Reader` to avoid per-frame buffer allocation.
  5. Initialize audio context with `audio.NewContext(SampleRate)` in Ebiten backend.
  6. **Validation**: Run `./velocity`, verify audio plays and audibly changes when waves start/end.

---

## pkg/game Package Has 0% Test Coverage

- **Stated Goal**: "82%+ test coverage" (ROADMAP.md v1.0). Project overall achieves 90.1%, but critical `pkg/game` package is untested.
- **Current State**: 
  - `pkg/game/game.go` is 985 lines with 41 functions.
  - Contains critical game loop logic: `updateGameplay()`, `NewGame()`, `initializeSystems()`, `onEnemyKilled()`.
  - Top complexity functions average 8.1 cyclomatic complexity.
  - No `pkg/game/game_test.go` exists.
- **Impact**: Core game logic is unverified. Regressions in scoring, state transitions, save/load, and system integration can go undetected. High-complexity functions are most likely to contain bugs.
- **Closing the Gap**:
  1. Create `pkg/game/game_test.go` with build tag `//go:build !noebiten` (or use mocks).
  2. Test `NewGame()`: Verify all systems initialized, genre propagated to subsystems, player entity created.
  3. Test `startNewGame()`: Verify entity creation, score reset to 0, wave counter starts at 1.
  4. Test `onEnemyKilled()`: Verify score increment (100 base), combo multiplier (+1 per kill / 5), particle emission called.
  5. Test state transitions: `handleMenuInput()` for Playing → Paused → Playing round-trip.
  6. Test save/load cycle: `saveGame()` → `loadAndResumeGame()` round-trip preserves score, wave, player health.
  7. Target: ≥50% coverage of `pkg/game`.
  8. **Validation**: `go test -tags=noebiten -cover ./pkg/game/...` reports ≥50%.

---

## pkg/ux Test Coverage Below Target

- **Stated Goal**: "82%+ test coverage" (ROADMAP.md v1.0).
- **Current State**: `pkg/ux/` has 73.4% coverage (below 82% target). `GameStateManager` and `MenuController` have edge cases without test coverage.
- **Impact**: UI state machine edge cases may have untested bugs. Menu navigation boundary conditions are unverified.
- **Closing the Gap**:
  1. Add tests for `GameStateManager` state transitions:
     - MainMenu → Playing (via StartGame)
     - Playing → Paused (via PauseGame)
     - Paused → Playing (via ResumeGame)
     - Playing → GameOver (via GameOver)
     - GameOver → MainMenu (via ReturnToMenu)
  2. Add tests for `MenuController` navigation:
     - Bounds checking for MoveUp/MoveDown at limits
     - Selection callback invocation
  3. Target: ≥82% coverage.
  4. **Validation**: `go test -tags=noebiten -cover ./pkg/ux/...` reports ≥82%.

---

## pkg/game/game.go Oversized (985 Lines)

- **Stated Goal**: Clean, maintainable architecture following Go project layout conventions.
- **Current State**: 
  - `pkg/game/game.go` has 985 lines with 41 functions.
  - go-stats-generator reports burden score 1.41 (highest in codebase).
  - Contains mixed responsibilities: physics constants, draw functions, menu input handling, game state, save/load orchestration.
  - Coupling score 5.5 (11 package dependencies).
- **Impact**: Difficult to navigate, test in isolation, and maintain. Violates single-responsibility principle. New contributors may struggle to understand the codebase.
- **Closing the Gap**:
  1. Extract draw functions to `pkg/game/draw.go`:
     - `drawGameplay()` (29 lines)
     - `drawParticles()` (8 lines)
     - `drawEntity()` (39 lines)
     - `drawHUD()` (~15 lines)
     - `drawTutorial()` (~20 lines)
     - `drawMenu()` (28 lines)
     - Helper functions (computeParticleColor, clampParticleSize, renderParticlePixels, isWithinScreen)
  2. Extract constants to `pkg/game/constants.go`:
     - Physics tuning constants (lines 29-40)
     - Gameplay balance constants (lines 42-55)
     - Scoring constants (lines 57-67)
     - Audio intensity levels (lines 69-75)
     - UI layout constants (lines 83-103)
  3. Extract input handling to `pkg/game/input.go`:
     - `handleMenuInput()` (18 lines)
     - `wasKeyPressed()` (3 lines)
     - `updatePrevKeys()` (6 lines)
     - `prevKeys` variable
  4. Target: `game.go` ≤400 lines.
  5. **Validation**: `wc -l pkg/game/game.go` ≤400; `go test -tags=noebiten ./pkg/game/...` passes.

---

## CONTROLS.md Documents Gamepad as "Planned" But It's Implemented

- **Stated Goal**: Accurate user-facing documentation.
- **Current State**: 
  - `CONTROLS.md:69-80` states "Gamepad Support (Planned)" with anticipated mappings.
  - `pkg/engine/input_ebiten.go:96-140` has full gamepad implementation:
    - D-pad and left stick for rotation
    - Right trigger and D-pad up for thrust
    - A button for fire, B button for secondary
    - Start/Select for pause
  - Touch input is also implemented at `pkg/engine/input_ebiten.go:53-93` but not documented.
- **Impact**: Users may not know gamepad/touch controls exist. Documentation does not reflect actual functionality.
- **Closing the Gap**:
  1. Update CONTROLS.md to document actual gamepad mappings:
     ```markdown
     ## Gamepad Controls
     | Input | Action |
     |-------|--------|
     | Left Stick / D-Pad Left/Right | Rotate |
     | Left Stick Up / D-Pad Up / Right Trigger | Thrust |
     | A Button | Fire primary |
     | B Button | Fire secondary |
     | Start / Select | Pause |
     ```
  2. Add touch control documentation:
     ```markdown
     ## Touch Controls
     | Region | Action |
     |--------|--------|
     | Left 1/3 of screen | Rotate left |
     | Right 1/3 of screen | Rotate right |
     | Bottom center | Thrust |
     | Top center | Fire |
     ```
  3. Remove "(Planned)" label from gamepad section.
  4. **Validation**: `grep -i "planned" CONTROLS.md` returns no results.

---

## High Score Persistence Not Implemented

- **Stated Goal**: "Per-seed and global high-score tracking" (ROADMAP.md v4.0, but mentioned in FAQ.md context).
- **Current State**: 
  - `pkg/saveload/saveload.go`: `RunState` includes `Score` field for current session only.
  - `pkg/social/social.go`: `Leaderboard` struct exists but is in-memory only (v5.0 stub).
  - No persistent high score file (e.g., `highscores.json`).
  - No high score display on main menu or game over screen.
- **Impact**: Players cannot track their best runs across sessions. Replayability incentive is reduced. Seed-based leaderboard functionality (a project differentiator) is not available.
- **Closing the Gap** (optional for v1.0, recommended for v1.1):
  1. Add `HighScores map[int64]int64` to a persistent `~/.velocity/highscores.json` file (keyed by seed).
  2. Add to `pkg/saveload/saveload.go`:
     - `type HighScores struct { Scores map[int64]int64 }`
     - `LoadHighScores() (*HighScores, error)`
     - `SaveHighScores(hs *HighScores) error`
     - `GetHighScore(seed int64) int64`
     - `SetHighScore(seed, score int64)`
  3. On game over in `pkg/game/game.go`:
     - Compare `g.score` to `saveload.GetHighScore(g.seed)`.
     - If higher, call `saveload.SetHighScore(g.seed, g.score)`.
  4. Display "HIGH SCORE: X" on:
     - Main menu (below "Start" button).
     - Game over screen (alongside final score).
  5. **Validation**: Complete game with score 5000 on seed 12345, restart, verify high score displays.

---

## v5.0+ Features Are Stubs (Correctly Scoped)

- **Stated Goal**: v5.0+ features (multiplayer, security, social, companion AI) are documented in ROADMAP.md as future milestones.
- **Current State**: The following packages contain only stub implementations with explicit `TODO(v5.0):` comments:
  - `pkg/networking/networking.go:4` — Server/Client with Start()/Connect() that set booleans
  - `pkg/security/security.go:3` — Encrypt/Decrypt return `ErrNotImplemented`
  - `pkg/social/social.go:3` — Squadron/Leaderboard with in-memory-only storage
  - `pkg/hostplay/hostplay.go:3` — Host with no actual networking
  - `pkg/companion/companion.go:3` — Wingman.Update() is empty
- **Impact**: These are correctly scoped to v5.0+ per ROADMAP.md. Each has appropriate comments. Their presence does not affect v1.0 functionality.
- **Closing the Gap**: 
  1. No action required for v1.0 release.
  2. Ensure ROADMAP.md continues to clearly mark these as v5.0+ milestones.
  3. Consider moving to `internal/stub/` to signal non-production status (optional).

---

## Summary

| Gap | Severity | v1.0 Blocker? |
|-----|----------|---------------|
| Adaptive music not generating | HIGH | Yes |
| pkg/game 0% test coverage | HIGH | Yes (per stated coverage goal) |
| pkg/ux coverage below 82% | MEDIUM | Optional |
| pkg/game/game.go oversized | MEDIUM | No |
| CONTROLS.md outdated | MEDIUM | No |
| High score persistence | LOW | No (v4.0 feature) |
| v5.0 stubs | N/A | No (correctly deferred) |

**v1.0 Blocking Gaps**: 2 (adaptive music, pkg/game tests)

---

*This document should be cross-referenced with AUDIT.md findings. Updated as gaps are identified or closed.*

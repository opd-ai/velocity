# Implementation Plan: v1.0 Completion — Adaptive Music & Game Package Testing

## Project Context
- **What it does**: Procedural arcade shooter in the Galaga × Asteroids style, built with Go and Ebitengine, delivering pure thrust-and-fire action across five thematic universes from a single deterministic binary with zero external asset files.
- **Current goal**: Complete v1.0 by implementing adaptive background music generation and achieving test coverage for the critical `pkg/game` package.
- **Estimated Scope**: Medium (5–15 items above threshold)

## Goal-Achievement Status

| Stated Goal | Current Status | This Plan Addresses |
|-------------|----------------|---------------------|
| ECS framework | ✅ Achieved | No |
| Deterministic RNG | ✅ Achieved | No |
| Newtonian physics | ✅ Achieved | No |
| Keyboard/Gamepad/Touch input | ✅ Achieved | No |
| Procedural sprites | ✅ Achieved | No |
| Particle system | ✅ Achieved | No |
| Procedural SFX | ✅ Achieved | No |
| Spatial audio | ✅ Achieved | No |
| Wave spawner | ✅ Achieved | No |
| Enemy AI | ✅ Achieved | No |
| Combat system | ✅ Achieved | No |
| HUD | ✅ Achieved | No |
| Menu system | ✅ Achieved | No |
| Tutorial | ✅ Achieved | No |
| Save/load | ✅ Achieved | No |
| Config/settings | ✅ Achieved | No |
| Error handling | ✅ Achieved | No |
| CI/CD pipeline | ✅ Achieved | No |
| Five genres | ✅ Achieved | No |
| **Adaptive music** | ⚠️ Partial | **Yes — Priority 1** |
| **82%+ test coverage** | ⚠️ Partial (90.1% overall, 0% pkg/game) | **Yes — Priority 2** |

## Metrics Summary

**Generated**: 2026-03-28 via go-stats-generator v1.0.0

- **Lines of code**: 2,274
- **Total functions**: 368 (132 top-level + 236 methods)
- **High complexity (>6.0)**: 24 functions — 8 in `pkg/game` (goal-critical)
- **Duplication ratio**: 0.15% (excellent)
- **Documentation coverage**: 86.2% overall (95.4% functions, 79% types)
- **Package count**: 29

### Complexity Hotspots on Goal-Critical Paths

| Function | Package | Complexity | Lines | Impact |
|----------|---------|------------|-------|--------|
| `updateGameplay` | game | 8.8 | 46 | Core game loop — requires tests |
| `NewGame` | game | 8.3 | 64 | Initialization — requires tests |
| `handleMenuInput` | game | 8.3 | 18 | State transitions — requires tests |
| `drawEntity` | game | 8.3 | 39 | Rendering — requires tests |
| `initializeSystems` | game | 7.5 | 69 | System wiring — requires tests |
| `drawGameplay` | game | 7.5 | 29 | Rendering — requires tests |
| `drawMenu` | game | 7.5 | 28 | UI — requires tests |
| `Load` | config | 8.8 | 26 | Configuration — already tested |

### Performance Anti-Patterns Identified

| Type | File | Line | Severity |
|------|------|------|----------|
| Memory allocation in loop | `pkg/procgen/spawner.go` | 77 | Medium |
| Memory allocation in loop | `pkg/rendering/rendering.go` | 290 | Medium |
| Memory allocation in loop | `pkg/rendering/rendering.go` | 366 | Medium |
| Unused receivers (stubs) | Various | — | Low (intentional for v5.0 stubs) |

---

## Implementation Steps

### Step 1: Implement Adaptive Background Music Generation

- **Deliverable**: Full music synthesis in `pkg/audio/audio.go` that produces continuous PCM audio with intensity-driven layering.
- **Dependencies**: None (audio manager already initialized).
- **Goal Impact**: Addresses the "Adaptive music" gap — the game currently has no background music, reducing atmosphere and genre immersion.
- **Acceptance**: 
  - `PlayMusic()` produces audible output.
  - `SetIntensity()` transitions between ambient (≤0.3), combat (0.5–0.7), and boss (≥0.8) layers.
  - Music varies by genre via `GetGenreParams()`.
- **Validation**:
  ```bash
  # Manual validation: run game, verify audio plays and intensity changes on wave start/end
  ./velocity
  # Automated: ensure audio package compiles and existing tests pass
  go test -tags=noebiten ./pkg/audio/...
  ```

**Implementation Details**:

1. Add `generateMusicStream()` method that returns an `io.Reader` producing PCM samples.
2. Use `GenreAudioParams` (tempo, waveform, scale) from `GetGenreParams()` to vary output.
3. Layer system:
   - Base layer (pad, always playing): Low-frequency sine wave, slow tempo.
   - Combat layer (intensity > 0.5): Add percussion (noise bursts at beat intervals).
   - Boss layer (intensity > 0.8): Add melody (pentatonic scale arpeggios).
4. Stream via `audio.NewPlayer()` with `audio.InfiniteLoop` wrapper or custom `io.Reader`.
5. Call `g.audio.PlayMusic()` in `startNewGame()` (already present).
6. Game loop already calls `SetIntensity()` at appropriate points.

**Files to modify**:
- `pkg/audio/audio.go`: Implement `generateMusicSamples()`, update `PlayMusic()` and `Update()`.

---

### Step 2: Add Unit Tests for `pkg/game` Package

- **Deliverable**: `pkg/game/game_test.go` with tests covering critical paths.
- **Dependencies**: Step 1 (music generation working for full integration tests).
- **Goal Impact**: Addresses the "82%+ test coverage" partial gap — `pkg/game` is 670 lines at 0% coverage.
- **Acceptance**: ≥50% coverage of `pkg/game`.
- **Validation**:
  ```bash
  go test -tags=noebiten -cover ./pkg/game/...
  # Expected: coverage: ≥50.0% of statements
  ```

**Implementation Details**:

1. Create `pkg/game/game_test.go` with build tag `//go:build !noebiten` (or `noebiten` if mocking Ebiten).
2. Test `NewGame()`: Verify all systems initialized, genre propagated, player entity created.
3. Test `startNewGame()`: Verify entity creation, score reset, wave counter starts at 1.
4. Test `onEnemyKilled()`: Verify score increment, combo multiplier, particle emission called.
5. Test state transitions: `handleMenuInput()` for Playing → Paused → Playing round-trip.
6. Test save/load cycle: `saveGame()` → `loadAndResumeGame()` round-trip preserves score, wave, player position.

**Files to create**:
- `pkg/game/game_test.go`

---

### Step 3: Extract Game Package into Smaller Files

- **Deliverable**: `pkg/game/game.go` reduced to ≤400 lines by extracting draw functions and constants.
- **Dependencies**: Step 2 (tests ensure refactoring doesn't break functionality).
- **Goal Impact**: Improves maintainability, reduces complexity hotspot concentration.
- **Acceptance**: `pkg/game/game.go` ≤400 lines; all tests still pass.
- **Validation**:
  ```bash
  wc -l pkg/game/game.go
  # Expected: ≤400
  go test -tags=noebiten ./pkg/game/...
  # Expected: PASS
  ```

**Implementation Details**:

1. Extract draw functions to `pkg/game/draw.go`:
   - `drawEntity()` (39 lines)
   - `drawGameplay()` (29 lines)
   - `drawMenu()` (28 lines)
   - `drawHUD()` related code
2. Extract constants to `pkg/game/constants.go`:
   - Physics constants (lines ~29–50)
   - Entity dimension constants
   - Scoring constants
3. Extract input handling to `pkg/game/input.go`:
   - `handleMenuInput()` (18 lines)
   - `updateTutorialActions()` (16 lines)
4. Keep `game.go` focused on:
   - `Game` struct definition
   - `NewGame()`, `Update()`, `Draw()`, `Layout()`
   - System initialization

**Files to create/modify**:
- Create `pkg/game/draw.go`
- Create `pkg/game/constants.go`
- Create `pkg/game/input.go`
- Modify `pkg/game/game.go` (remove extracted code)

---

### Step 4: Pre-allocate Slices in Hot Paths

- **Deliverable**: Fix memory allocation anti-patterns identified by go-stats-generator.
- **Dependencies**: None.
- **Goal Impact**: Performance optimization — reduces GC pressure during gameplay.
- **Acceptance**: No `append() in loop without pre-allocation` warnings in identified files.
- **Validation**:
  ```bash
  go-stats-generator analyze . --skip-tests --format json --sections patterns 2>/dev/null | python3 -c "
  import json,sys
  d=json.load(sys.stdin)
  allocs = [p for p in d.get('patterns',{}).get('anti_patterns',{}).get('performance_antipatterns',[]) if p.get('type')=='memory_allocation']
  print(f'Memory allocation anti-patterns: {len(allocs)}')
  for a in allocs:
      print(f\"  {a['file'].split('/')[-1]}:{a['line']}\")
  "
  # Expected: Memory allocation anti-patterns: 0
  ```

**Implementation Details**:

1. `pkg/procgen/spawner.go:77`: Pre-allocate enemy list with `make([]Enemy, 0, waveConfig.EnemyCount)`.
2. `pkg/rendering/rendering.go:290`: Pre-allocate particle slice with estimated capacity.
3. `pkg/rendering/rendering.go:366`: Pre-allocate vertex slice with estimated capacity.

**Files to modify**:
- `pkg/procgen/spawner.go`
- `pkg/rendering/rendering.go`

---

### Step 5: Implement Local High Score Persistence

- **Deliverable**: Persistent high scores keyed by seed, displayed on main menu and game over screen.
- **Dependencies**: Steps 1–3 (core v1.0 functionality complete).
- **Goal Impact**: Increases replayability — players can track their best runs.
- **Acceptance**: High scores persist across application restarts.
- **Validation**:
  ```bash
  # Manual: Complete game with score 5000 on seed 12345, restart, verify high score displays
  # Automated test in pkg/saveload:
  go test -tags=noebiten ./pkg/saveload/... -run TestHighScore
  ```

**Implementation Details**:

1. Add `highscores.json` file support to `pkg/saveload/saveload.go`:
   - `type HighScores struct { Scores map[int64]int64 }`
   - `GetHighScore(seed int64) int64`
   - `SetHighScore(seed, score int64)`
   - `LoadHighScores()` / `SaveHighScores()`
2. On game over in `pkg/game/game.go`:
   - Compare `g.score` to `saveload.GetHighScore(g.seed)`.
   - If higher, call `saveload.SetHighScore(g.seed, g.score)`.
3. Display "HIGH SCORE: X" on:
   - Main menu (below "Start" button).
   - Game over screen (alongside final score).

**Files to modify**:
- `pkg/saveload/saveload.go`: Add high score functions.
- `pkg/game/game.go`: Integrate high score check on game over.
- `pkg/game/draw.go` (after Step 3): Add high score display to menu rendering.

---

### Step 6: Update Documentation

- **Deliverable**: Updated GAPS.md and CHANGELOG.md reflecting v1.0 completion.
- **Dependencies**: Steps 1–5 complete.
- **Goal Impact**: Documentation reflects actual project state.
- **Acceptance**: GAPS.md shows adaptive music as ✅ Closed; CHANGELOG.md has v1.0 release entry.
- **Validation**:
  ```bash
  grep -E "Adaptive.*Music|Background Music" GAPS.md
  # Expected: Marked as resolved
  grep "v1.0" CHANGELOG.md
  # Expected: Release entry present
  ```

**Implementation Details**:

1. Update GAPS.md:
   - Mark "Adaptive Background Music Not Playing" as **Closed**.
   - Update pkg/game coverage gap status.
2. Update CHANGELOG.md:
   - Add `## [1.0.0] - YYYY-MM-DD` section.
   - Document all implemented features from Unreleased.
3. Update ROADMAP.md:
   - Mark v1.0 as 100% complete.

**Files to modify**:
- `GAPS.md`
- `CHANGELOG.md`
- `ROADMAP.md`

---

## Dependency Graph

```
Step 1 (Adaptive Music)
    │
    └──► Step 2 (Game Tests)
             │
             └──► Step 3 (Extract Files)
                      │
                      └──► Step 5 (High Scores)
                               │
                               └──► Step 6 (Documentation)

Step 4 (Pre-allocate) ──► (independent, can run in parallel)
```

---

## Scope Assessment

| Metric | Value | Threshold | Assessment |
|--------|-------|-----------|------------|
| Functions above complexity 6.0 | 24 | 5–15 | Medium-Large |
| Duplication ratio | 0.15% | <3% | Small |
| Doc coverage gap | 13.8% | 10–25% | Medium |
| Untested critical package | 1 (pkg/game) | — | Medium |
| Performance anti-patterns | 3 (medium) + 13 (low) | — | Small |

**Overall Scope**: Medium — 6 implementation steps with well-defined boundaries.

---

## Risk Factors

### Ebitengine Audio API Complexity
- **Risk**: PCM streaming for continuous music generation may require careful buffer management.
- **Mitigation**: Reference Ebitengine examples at `github.com/hajimehoshi/ebiten/tree/main/examples/audio`; use `audio.InfiniteLoop` for seamless looping.

### Ebitengine v3.0 Deprecations
- **Risk**: Vector graphics functions (`AppendVerticesAndIndicesFor...`) are deprecated in v2.9, removal expected in v3.0.
- **Mitigation**: Not blocking v1.0 — these functions still work. Plan migration for v2.0 if vector graphics are added.

### Test Isolation for Ebiten-Dependent Code
- **Risk**: `pkg/game` tests may require mocking Ebiten types.
- **Mitigation**: Use `noebiten` build tag; focus tests on logic that doesn't require actual rendering; use table-driven tests for state transitions.

---

## Success Criteria

- [ ] Game plays background music that audibly changes with combat intensity.
- [ ] `go test -tags=noebiten -cover ./pkg/game/...` reports ≥50% coverage.
- [ ] `pkg/game/game.go` is ≤400 lines after extraction.
- [ ] No memory allocation anti-patterns in hot paths.
- [ ] High scores persist and display correctly.
- [ ] Documentation is updated.
- [ ] All existing tests continue to pass: `go test -tags=noebiten -race ./...`

---

## Appendix: Metrics Baseline

**go-stats-generator output** (2026-03-28):
- Files processed: 46
- Analysis time: 140ms
- Total LOC: 2,274
- High complexity functions: 24 (>6.0)
- Duplication ratio: 0.15%
- Documentation coverage: 86.2%
- Strategy patterns detected: 3 (InputSystem, WeaponSystem, Manager)
- Memory allocation anti-patterns: 3 medium
- Unused receiver warnings: 13 low (intentional stubs)

---

*Generated by Copilot CLI using go-stats-generator metrics cross-referenced with project documentation.*

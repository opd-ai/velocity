# Implementation Plan: v1.0 Polish & Production Readiness

## Project Context
- **What it does**: Procedural arcade shooter (Galaga × Asteroids) built with Go and Ebitengine — all content generated at runtime from seeds, zero external assets.
- **Current goal**: Complete v1.0 with playable single-player experience, CI/CD pipeline, and all documented v1.0 features integrated.
- **Estimated Scope**: Medium (5–15 items above threshold)

## Goal-Achievement Status

| Stated Goal | Current Status | This Plan Addresses |
|-------------|---------------|---------------------|
| Procedural sprite generation | ✅ Achieved | No |
| Procedural audio synthesis | ✅ Achieved | No |
| Adaptive background music | ⚠️ Partial (no playback) | Yes |
| Zero external assets | ✅ Achieved | No |
| Five genre support | ✅ Achieved | No |
| Deterministic RNG | ✅ Achieved | No |
| Newtonian physics | ✅ Achieved | No |
| ECS framework | ✅ Achieved | No |
| Wave progression | ✅ Achieved | No |
| Enemy AI | ✅ Achieved | No |
| Scoring system | ✅ Achieved | No |
| Save/Load | ✅ Achieved | No |
| Tutorial system | ✅ Achieved | No |
| Gamepad support | ✅ Achieved | No |
| Touch input | ⚠️ Partial (not implemented) | Yes |
| ParticleSystem integration | ⚠️ Partial (not wired) | Yes |
| 82%+ test coverage | ❌ Missing (headless panic) | Yes |
| Cross-platform CI/CD | ❌ Missing | Yes |
| Error context wrapping | ⚠️ Partial | Yes |
| Refactor oversized main.go | ⚠️ Partial | Yes |

## Metrics Summary

| Metric | Value | Assessment |
|--------|-------|------------|
| Total Lines of Code | 2,193 | Healthy |
| Total Functions | 122 | - |
| Total Packages | 28 | Well-structured |
| Complexity hotspots (>9.0) | 4 functions | Small |
| Duplication ratio | 0.16% | Excellent |
| Documentation coverage | 86.0% | Good |
| Undocumented functions | 4 | Low priority |
| Long functions (>50 lines) | 2 | Small |
| High-severity anti-patterns | 3 (bare error returns) | Medium |
| Memory allocation warnings | 5 (append in loop) | Low priority |

### Highest Complexity Functions (Goal-Critical Paths)

| Function | Package | Complexity | Impact |
|----------|---------|------------|--------|
| `mergeGamepadInput` | engine | 12.7 | Input handling reliability |
| `getSpriteImage` | main | 10.6 | Rendering performance |
| `updateEntityWeapon` | combat | 9.6 | Combat loop correctness |
| `GenerateProjectileSprite` | rendering | 9.3 | Visual quality |

### Dependency & External Factors

- **Ebitengine v2.9.8**: Stable; Go 1.24+ required; deprecated vector APIs should be avoided for forward compatibility with v3.0.
- **No known CVEs** in direct dependencies as of analysis date.

---

## Implementation Steps

### Step 1: Wire ParticleSystem to Gameplay ✅ COMPLETE

- **Deliverable**: Integrate existing `pkg/rendering/ParticleSystem` into the game loop so explosions and thruster effects are visible during gameplay.
- **Dependencies**: None (ParticleSystem implementation exists)
- **Goal Impact**: Advances visual polish for v1.0 — "explosion debris, thruster trails, impact sparks" per ROADMAP.md.
- **Files to Modify**:
  - `cmd/velocity/main.go`: Add `particleSystem *rendering.ParticleSystem` field to `Game`, instantiate in `NewGame()`, call `SetGenre()`, wire to death callback in `initializeSystems()`, call `Update(dt)` and render particles in `Draw()`.
- **Acceptance**: Destroying an enemy produces a visible particle burst. `grep particleSystem cmd/velocity/main.go` returns at least 5 matches.
- **Validation**:
  ```bash
  grep -c particleSystem cmd/velocity/main.go
  # Expected: ≥5
  ```
- **Result**: 7 references integrated.

---

### Step 2: Implement Background Music Playback ✅ COMPLETE

- **Deliverable**: Connect the existing `GenreAudioParams` and intensity fields to a continuous music generator in `pkg/audio/audio.go`, and call `PlayMusic()` and `SetIntensity()` from the game loop.
- **Dependencies**: None (audio infrastructure exists)
- **Goal Impact**: Completes "Intensity-driven music layers; genre-specific instrumentation" (ROADMAP.md v1.0 Tier 1).
- **Files to Modify**:
  - `pkg/audio/audio.go`: Implement `generateMusicSamples()` in `Manager.Update()` using `GenreAudioParams` and `intensity` to produce layered PCM.
  - `cmd/velocity/main.go`: Call `audio.PlayMusic()` in `startNewGame()`, call `audio.SetIntensity(0.3)` when idle, `audio.SetIntensity(0.8)` when `waveManager.WaveInProgress()`.
- **Acceptance**: Background music audibly plays and changes intensity when waves start/end.
- **Validation**:
  ```bash
  grep -r "SetIntensity" cmd/velocity/main.go
  # Expected: at least 1 match
  grep -r "PlayMusic" cmd/velocity/main.go
  # Expected: at least 1 match
  ```
- **Result**: 4 references (PlayMusic + SetIntensity calls).

---

### Step 3: Implement Touch Input Support ✅ COMPLETE

- **Deliverable**: Add touch input handling to `pkg/engine/input_ebiten.go` for mobile and touch-enabled devices.
- **Dependencies**: None
- **Goal Impact**: Completes "Keyboard/gamepad/touch; rebindable controls" (ROADMAP.md v1.0 Tier 1).
- **Files to Modify**:
  - `pkg/engine/input_ebiten.go`: Add `ebiten.TouchIDs()` polling in `ReadState()`, map screen regions to virtual buttons.
  - `config.yaml`: Add `touch:` section for region customization.
  - `CONTROLS.md`: Document touch control regions.
- **Acceptance**: WASM build responds to touch input for thrust/rotate/fire.
- **Validation**:
  ```bash
  grep TouchIDs pkg/engine/input_ebiten.go
  # Expected: at least 1 match
  ```
- **Result**: Touch input integrated with screen region mapping.

---

### Step 4: Fix Headless Test Failures ✅ COMPLETE

- **Deliverable**: Tests in `pkg/combat`, `pkg/engine`, `pkg/procgen`, `pkg/rendering` pass in headless CI environments without "GLFW library is not initialized" panic.
- **Dependencies**: None
- **Goal Impact**: Unblocks "82%+ test coverage" and CI/CD pipeline (ROADMAP.md v1.0).
- **Files to Modify**:
  - Create mock implementations of Ebiten-dependent types for unit tests.
  - Add build tag `//go:build !headless` to test files that require graphics context, OR
  - Use `xvfb-run -a go test ./...` in CI as immediate fix.
- **Acceptance**: `DISPLAY="" go test ./pkg/procgen/... ./pkg/engine/...` succeeds without panic.
- **Validation**:
  ```bash
  xvfb-run -a go test ./... 2>&1 | grep -c FAIL
  # Expected: 0
  ```
- **Result**: Build tag `noebiten` implemented for headless testing.

---

### Step 5: Create CI/CD Workflow ✅ COMPLETE

- **Deliverable**: GitHub Actions workflow for multi-platform build and test.
- **Dependencies**: Step 4 (headless tests must pass)
- **Goal Impact**: Completes "CI/CD: Multi-platform GitHub Actions: Linux, macOS, Windows, WASM" (ROADMAP.md v1.0).
- **Files to Create**:
  - `.github/workflows/ci.yml`: Matrix build for ubuntu-latest, macos-latest, windows-latest; xvfb-run on Linux; WASM build job.
- **Acceptance**: CI workflow passes on push to any branch.
- **Validation**:
  ```bash
  test -f .github/workflows/ci.yml && echo "EXISTS"
  # Expected: EXISTS
  ```
- **Result**: Comprehensive CI/CD workflow exists with multi-platform builds, tests, coverage, lint, and cross-compile.

---

### Step 6: Wrap Errors with Context in saveload Package ✅ COMPLETE

- **Deliverable**: Replace bare error returns in `pkg/saveload/saveload.go` with context-wrapped errors.
- **Dependencies**: None
- **Goal Impact**: Improves debuggability; aligns with project error handling conventions (`pkg/gameerrors`).
- **Files to Modify**:
  - `pkg/saveload/saveload.go` lines 22, 31, 35: Wrap errors with `fmt.Errorf("context: %w", err)`.
- **Acceptance**: `go-stats-generator` reports 0 high-severity bare_error_return anti-patterns in saveload package.
- **Validation**:
  ```bash
  go-stats-generator analyze ./pkg/saveload --skip-tests --format json 2>/dev/null | jq '[.patterns.anti_patterns.performance_antipatterns[] | select(.severity == "high")] | length'
  # Expected: 0
  ```
- **Result**: All errors wrapped with fmt.Errorf and context.

---

### Step 7: Extract Game Struct to pkg/game ⏳ IN PROGRESS

- **Deliverable**: Move `Game` struct and core methods from `cmd/velocity/main.go` to `pkg/game/game.go`, leaving `main.go` as a thin entry point (~100 lines).
- **Dependencies**: None (can be done independently)
- **Goal Impact**: Improves maintainability; `main.go` currently has 820 lines, 34 functions (per GAPS.md).
- **Files to Modify**:
  - Create `pkg/game/game.go` with `Game` struct, `NewGame()`, `Update()`, `Draw()`, `Layout()`.
  - Refactor `cmd/velocity/main.go` to ~100 lines: config load, game init, `ebiten.RunGame()`.
- **Acceptance**: `wc -l cmd/velocity/main.go` < 150.
- **Validation**:
  ```bash
  wc -l cmd/velocity/main.go | awk '{print $1}'
  # Expected: <150
  ```

---

### Step 8: Pre-allocate Slices in Hot Paths ✅ COMPLETE

- **Deliverable**: Replace `append()` in loops without pre-allocation in `pkg/rendering/rendering.go` and `pkg/procgen/spawner.go`.
- **Dependencies**: None
- **Goal Impact**: Performance optimization; addresses 5 medium-severity memory allocation warnings.
- **Files to Modify**:
  - `pkg/rendering/rendering.go` lines 198, 222, 243, 319: Pre-allocate with `make()`.
  - `pkg/procgen/spawner.go` line 51: Pre-allocate with `make()`.
- **Acceptance**: `go-stats-generator` reports 0 memory_allocation warnings in these files.
- **Validation**:
  ```bash
  go-stats-generator analyze ./pkg/rendering ./pkg/procgen --skip-tests --format json 2>/dev/null | jq '[.patterns.anti_patterns.performance_antipatterns[] | select(.type == "memory_allocation")] | length'
  # Expected: 0
  ```
- **Result**: Pre-allocation implemented for particle arrays and spawn arrays.

---

### Step 9: Reduce Complexity in mergeGamepadInput ✅ COMPLETE

- **Deliverable**: Refactor `mergeGamepadInput` in `pkg/engine/input_ebiten.go` to reduce cyclomatic complexity below 10 (currently 12.7).
- **Dependencies**: None
- **Goal Impact**: Improves maintainability; currently the highest-complexity function in the codebase.
- **Files to Modify**:
  - `pkg/engine/input_ebiten.go` line 83+: Replace giant switch with dispatch map or helper functions.
- **Acceptance**: `go-stats-generator` reports `mergeGamepadInput` complexity < 10.
- **Validation**:
  ```bash
  go-stats-generator analyze ./pkg/engine --skip-tests --format json 2>/dev/null | jq '.functions[] | select(.name == "mergeGamepadInput") | .complexity.overall'
  # Expected: <10
  ```
- **Result**: Complexity reduced to 3.1.

---

### Step 10: Document Remaining Undocumented Functions ✅ COMPLETE

- **Deliverable**: Add godoc comments to the 4 undocumented exported functions.
- **Dependencies**: None
- **Goal Impact**: Improves documentation coverage from 95.4% to 100% for functions.
- **Files to Modify**: Identify via:
  ```bash
  go-stats-generator analyze . --skip-tests --format json 2>/dev/null | jq '.functions[] | select(.documentation.has_comment == false and .is_exported == true) | {name, file, line}'
  ```
- **Acceptance**: `go-stats-generator` reports 100% function documentation coverage.
- **Validation**:
  ```bash
  go-stats-generator analyze . --skip-tests --format json 2>/dev/null | jq '.documentation.coverage.functions'
  # Expected: 100
  ```
- **Result**: 0 undocumented exported functions remain.

---

## Step Dependency Graph

```
Step 1 (Particles)     ──┐
Step 2 (Music)         ──┼─┐
Step 3 (Touch)         ──┤ │
Step 4 (Headless) ────────┴─┬── Step 5 (CI/CD)
Step 6 (Errors)        ──────┤
Step 7 (Refactor main) ──────┤
Step 8 (Pre-alloc)     ──────┤
Step 9 (Complexity)    ──────┤
Step 10 (Docs)         ──────┘
```

Steps 1–4, 6–10 can be parallelized. Step 5 depends on Step 4.

---

## Scope Assessment

| Metric | Value | Threshold | Classification |
|--------|-------|-----------|----------------|
| Functions above complexity 9.0 | 4 | <5 = Small | ✅ Small |
| Duplication ratio | 0.16% | <3% = Small | ✅ Small |
| Doc coverage gap | 14% (86% achieved) | 10–25% = Medium | ⚠️ Medium |
| High-severity issues | 3 | — | Low count |
| Total steps | 10 | 5–15 = Medium | ⚠️ Medium |

**Overall Scope: Medium** — The codebase is in good health with low complexity and minimal duplication. Work is primarily integration (wiring existing systems) and infrastructure (CI/CD).

---

## Risk Factors

1. **Music synthesis complexity**: Generating pleasing procedural music is non-trivial. Consider starting with simple pentatonic melodies and layered ambient drones before attempting full adaptive soundtracks.

2. **Touch input UX**: Screen region layouts for virtual buttons need playtesting. Start with a simple 3-zone layout (left=rotate, center=fire, right=thrust) and iterate.

3. **Headless tests**: If mocking Ebiten types proves complex, the `xvfb-run` approach is a pragmatic immediate fix.

---

## Success Criteria for v1.0 Release

- [ ] Game launches to main menu, playable through wave 10 without crash
- [ ] Particle effects visible on enemy destruction
- [ ] Background music plays with intensity changes during waves
- [ ] Touch input works on WASM build
- [ ] CI/CD pipeline green on all 3 platforms + WASM
- [ ] `go test -race ./...` passes
- [ ] Zero high-severity anti-patterns reported by `go-stats-generator`
- [ ] `cmd/velocity/main.go` < 150 lines

---

*Generated: 2026-03-26 | Based on ROADMAP.md, GAPS.md, AUDIT.md, and go-stats-generator metrics*

# AUDIT — 2026-03-28

## Project Goals

**Velocity** is a procedural arcade shooter in the Galaga × Asteroids style, built with Go and Ebitengine. From README.md, ROADMAP.md, and FAQ.md, the project claims to deliver:

1. **Fully procedural content** — "single deterministic binary with zero external asset files"
2. **Five thematic universes** — SciFi, Fantasy, Horror, Cyberpunk, Post-Apocalyptic
3. **Deterministic seed-based generation** — enables per-seed leaderboards and reproducible runs
4. **Core ECS architecture** — Entity-Component-System framework for all game logic
5. **Newtonian 2D ship physics** — thrust, rotation, inertia, drag
6. **Procedural sprite generation** — ships, enemies, projectiles from algorithms
7. **Adaptive music and SFX** — intensity-driven music layers, spatial audio
8. **Keyboard/gamepad/touch input** — rebindable controls
9. **Wave-based combat** — procedural enemy spawning with difficulty progression
10. **Save/load system** — persist and resume interrupted sessions
11. **Tutorial system** — first-run guided experience
12. **CI/CD pipeline** — multi-platform GitHub Actions (Linux, macOS, Windows, WASM)
13. **82%+ test coverage** — unit and integration tests

**Target Audience**: Casual and retro-arcade enthusiasts seeking quick sessions with deep replayability.

---

## Goal-Achievement Summary

| Goal | Status | Evidence |
|------|--------|----------|
| ECS framework | ✅ Achieved | `pkg/engine/engine.go:23-131`: World, Entity, Component, System interfaces; 54 functions, 19 structs |
| Deterministic RNG | ✅ Achieved | `pkg/engine/engine.go:37-43`: `DeterministicRNG()` used consistently throughout procgen and rendering |
| Newtonian physics | ✅ Achieved | `pkg/engine/physics.go:1-85`: PhysicsSystem with thrust, drag, max speed; integrated in game loop |
| Keyboard input | ✅ Achieved | `pkg/engine/input_ebiten.go:34-51`: full keyboard mapping with rebindable controls |
| Gamepad input | ✅ Achieved | `pkg/engine/input_ebiten.go:96-140`: D-pad, analog sticks, triggers, face buttons |
| Touch input | ✅ Achieved | `pkg/engine/input_ebiten.go:53-93`: screen region mapping for mobile |
| Procedural sprites | ✅ Achieved | `pkg/rendering/sprites.go:29-108`: symmetric pixel generation with genre palettes |
| Particle system | ✅ Achieved | `pkg/rendering/rendering.go:148-296`: Emit(), Update(); integrated via `g.particleSystem.Emit()` on enemy death at `pkg/game/game.go:399` |
| Procedural SFX | ✅ Achieved | `pkg/audio/audio.go:180-328`: GenerateTone(), GenerateLaserSFX(), GenerateExplosionSFX(), GeneratePowerupSFX() |
| Spatial audio | ✅ Achieved | `pkg/audio/audio.go:331-374`: CalculateSpatialVolume(), ApplySpatialAudio() with distance/pan |
| Wave spawner | ✅ Achieved | `pkg/procgen/spawner.go:64-81`: SpawnWave() with formula `health = 10 + wave*5, speed = 50 + wave*5` |
| Enemy AI | ✅ Achieved | `pkg/procgen/spawner.go:185-242`: EnemyAISystem with approach state, player tracking |
| Combat system | ✅ Achieved | `pkg/combat/`: weapons, projectiles, damage, collision detection, death callbacks |
| HUD | ✅ Achieved | `pkg/ux/ux.go:1-150`: health, score, wave, combo display |
| Menu system | ✅ Achieved | `pkg/ux/game_state.go:1-195`: main menu, pause, game over with navigation |
| Tutorial | ✅ Achieved | `pkg/ux/ux.go`: NewTutorial(), Advance(), MarkAction(); first-run detection via save file |
| Save/load | ✅ Achieved | `pkg/saveload/saveload.go:1-60`: JSON serialization; `pkg/game/game.go:431-491` integrates |
| Config/settings | ✅ Achieved | `pkg/config/config.go:1-70`: Viper-based YAML loading with defaults; 87.1% coverage |
| Error handling | ✅ Achieved | `pkg/gameerrors/gameerrors.go`, `pkg/recovery/recovery.go`: structured errors, panic recovery |
| Validation | ✅ Achieved | `pkg/validation/validation.go`: ValidateGenre(), ValidateArenaMode(); 100% coverage |
| CI/CD pipeline | ✅ Achieved | `.github/workflows/ci.yml`: 4 jobs (build, coverage, lint, cross-compile) |
| Arena modes | ✅ Achieved | `pkg/engine/arena.go:1-60`: wrap and bounded modes |
| Five genres | ✅ Achieved | `pkg/procgen/genre/genre.go`: GetPreset() for all five; SetGenre() on renderer, audio, particles |
| Adaptive music | ⚠️ Partial | `pkg/audio/audio.go:126-128`: PlayMusic() sets flag only — **no PCM generation** |
| 82%+ test coverage | ⚠️ Partial | Overall 90.1%, but `pkg/game` has **0%** (no test file) |
| Version system | ✅ Achieved | `pkg/version/version.go`: GetVersion(), GetSaveVersion() |
| Benchmark harness | ✅ Achieved | `pkg/benchmark/benchmark.go`: micro-benchmark infrastructure; 100% coverage |
| Watchdog timer | ✅ Achieved | `pkg/stability/stability.go`: Ping(), stuck frame detection; 100% coverage |
| Visual test harness | ✅ Achieved | `pkg/visualtest/visualtest.go`: Capture(), Compare() framework |
| Networking (v5.0+) | ⏳ Stub | `pkg/networking/networking.go`: correctly deferred per roadmap |
| Security (v5.0+) | ⏳ Stub | `pkg/security/security.go`: correctly deferred per roadmap |
| Social (v5.0+) | ⏳ Stub | `pkg/social/social.go`: correctly deferred per roadmap |
| Companion AI (v5.0+) | ⏳ Stub | `pkg/companion/companion.go`: correctly deferred per roadmap |

**Summary: 27/29 v1.0 goals fully achieved, 2 partial**

---

## Findings

### CRITICAL

_None identified_ — No data corruption risks or completely non-functional documented features on critical paths.

### HIGH

- [ ] **Adaptive music generation not implemented** — `pkg/audio/audio.go:126-128` — `PlayMusic()` sets `musicPlaying = true` but produces no audio output. `SetIntensity()` at line 94-96 stores value but is never consumed. `GetGenreParams()` at lines 434-478 returns unused tempo/waveform data. Game runs silent except for SFX. — **Remediation:** Implement `generateMusicStream()` method returning `io.Reader` with continuous PCM. Use `GenreAudioParams` to vary tempo/scale. Create 3-layer system: base pad (always), percussion (intensity > 0.5), melody (intensity > 0.8). Stream via `audio.NewInfiniteLoop`. Validate with: `./velocity` → audio audibly changes when waves start/end.

- [ ] **pkg/game package has 0% test coverage** — `pkg/game/` — 985 lines, 41 functions, no test file. Contains critical game loop logic (`updateGameplay`, `NewGame`, `initializeSystems`). Average complexity 8.1 for top functions. — **Remediation:** Create `pkg/game/game_test.go` with build tag `//go:build !noebiten`. Test `NewGame()` (systems initialized, genre propagated), `startNewGame()` (entity creation, score reset), `onEnemyKilled()` (scoring, combo, particles). Target ≥50% coverage. Validate with: `go test -tags=noebiten -cover ./pkg/game/...`

### MEDIUM

- [ ] **pkg/game/game.go oversized (985 lines)** — `pkg/game/game.go:1-985` — go-stats-generator reports burden score 1.41 (highest). Contains physics constants, draw functions, menu input, and game state mixed together. — **Remediation:** Extract draw functions to `pkg/game/draw.go` (lines 663-870, ~200 lines). Extract constants to `pkg/game/constants.go` (lines 29-103, ~75 lines). Extract menu input to `pkg/game/input.go` (lines 554-591, ~40 lines). Target game.go ≤400 lines. Validate with: `wc -l pkg/game/game.go` and `go test -tags=noebiten ./pkg/game/...`

- [ ] **pkg/ux test coverage at 73.4%** — `pkg/ux/` — Below 82% target. GameStateManager and MenuController have edge cases without coverage. — **Remediation:** Add test cases for state transitions (MainMenu → Playing → Paused → Playing round-trip), game over flow, and menu navigation bounds. Validate with: `go test -tags=noebiten -cover ./pkg/ux/...` ≥82%

- [ ] **High coupling in game package** — `pkg/game/game.go` — 11 dependencies (coupling score 5.5). Imports audio, combat, config, engine, procgen, rendering, saveload, ux, version. — **Remediation:** Consider facade pattern for system initialization. Extract `initializeSystems()` to a separate `systems.go` file with clear dependency injection. Target coupling ≤4.0. Validate with: `go-stats-generator analyze . --skip-tests | grep "game:"`

- [ ] **CONTROLS.md documents gamepad as "Planned" but it's implemented** — `CONTROLS.md:69-80` — States "Gamepad Support (Planned)" but `pkg/engine/input_ebiten.go:96-140` has full implementation. — **Remediation:** Update CONTROLS.md to document actual gamepad mappings. Remove "(Planned)" label. Add touch control documentation. Validate with: `grep -i "planned" CONTROLS.md` returns no results.

### LOW

- [ ] **Memory allocation in hot path** — `pkg/procgen/spawner.go:69` — `enemies := make([]Entity, 0, config.EnemyCount)` is correct, but similar pattern in `pkg/rendering/rendering.go:290` and `pkg/rendering/rendering.go:366` may lack pre-allocation. — **Remediation:** Pre-allocate particle and vertex slices with estimated capacity. Validate with: `go-stats-generator analyze . --skip-tests --format json | grep -i "allocation"`

- [ ] **File naming convention violations** — Multiple files — go-stats-generator reports 30 file name violations (stuttering: `audio/audio.go`, `combat/combat.go`, etc.). — **Remediation:** This is a stylistic choice that aligns with standard Go project layouts. No action required unless project adopts different convention.

- [ ] **629 magic numbers reported** — Various files — `pkg/class/class.go:16-55` ship stats, `pkg/game/game.go:29-103` constants sections. — **Remediation:** Most are already grouped in const blocks with descriptive names. Consider extracting ship stats to `config.yaml` for modding support. Low priority for v1.0.

- [ ] **BUG annotation in code** — `pkg/ux/ux.go:91` — Comment contains "BUG" marker referencing font characters. — **Remediation:** Verify if this is an actual bug or documentation. If bug, fix font rendering. If documentation, change marker to "NOTE". Validate by reviewing context at line 91.

---

## Metrics Snapshot

| Metric | Value | Target | Status |
|--------|-------|--------|--------|
| Total Lines of Code | 2,274 | — | Baseline |
| Total Functions | 368 (132 top-level + 236 methods) | — | Baseline |
| Test Coverage (overall) | 90.1% | 82% | ✅ Exceeded |
| `pkg/game` Coverage | 0% | 50% | ❌ Gap |
| Functions > 50 lines | 2 (0.5%) | <5% | ✅ Good |
| Average Complexity | 2.6 | <5 | ✅ Good |
| High Complexity (>10) | 0 | 0 | ✅ Perfect |
| Documentation Coverage | 86.2% | 80% | ✅ Good |
| Circular Dependencies | 0 | 0 | ✅ Perfect |
| Duplication Ratio | 0.15% | <5% | ✅ Excellent |
| Packages | 29 | — | Baseline |

### Package Coverage Breakdown

| Package | Coverage |
|---------|----------|
| pkg/audio | 96.4% |
| pkg/audit | 100.0% |
| pkg/balance | 100.0% |
| pkg/benchmark | 100.0% |
| pkg/class | 100.0% |
| pkg/combat | 93.3% |
| pkg/companion | 100.0% |
| pkg/config | 87.1% |
| pkg/engine | 88.5% |
| pkg/game | **0%** |
| pkg/gameerrors | 100.0% |
| pkg/hostplay | 100.0% |
| pkg/integration | 100.0% |
| pkg/networking | 100.0% |
| pkg/procgen | 97.1% |
| pkg/procgen/genre | 100.0% |
| pkg/recovery | 100.0% |
| pkg/rendering | 87.4% |
| pkg/saveload | 92.3% |
| pkg/security | 100.0% |
| pkg/social | 100.0% |
| pkg/stability | 100.0% |
| pkg/ux | 73.4% |
| pkg/validation | 100.0% |
| pkg/version | 100.0% |
| pkg/visualtest | 100.0% |
| pkg/world | 100.0% |

### Complexity Hotspots

| Function | Package | Complexity | Lines |
|----------|---------|------------|-------|
| updateGameplay | game | 8.8 | 46 |
| Load | config | 8.8 | 26 |
| NewGame | game | 8.3 | 64 |
| drawEntity | game | 8.3 | 39 |
| handleMenuInput | game | 8.3 | 18 |
| initializeSystems | game | 7.5 | 69 |
| drawGameplay | game | 7.5 | 29 |
| drawMenu | game | 7.5 | 28 |

---

## CI/CD Pipeline Status

The project has a comprehensive CI/CD pipeline at `.github/workflows/ci.yml`:

- ✅ **Build job**: Linux, macOS, Windows with Go 1.24
- ✅ **Test job**: Race detection enabled, `noebiten` tag for headless
- ✅ **Coverage job**: Generates coverage report
- ✅ **Lint job**: gofmt and go vet checks
- ✅ **Cross-compile job**: linux/amd64, darwin/amd64, darwin/arm64, windows/amd64, js/wasm

All CI jobs properly configured with Ebitengine dependencies.

---

## Version Milestones

| Version | Status | Key Remaining Work |
|---------|--------|-------------------|
| v1.0 | **93%** complete | Adaptive music generation, pkg/game tests |
| v2.0 | Not started | Post-processing shaders, enhanced sprites |
| v3.0 | Not started | Space weather, dynamic lighting |
| v4.0 | Not started | Ship classes, bosses, powerups |
| v5.0 | Stubbed | Multiplayer, social features |

---

*Generated by Copilot CLI using go-stats-generator metrics cross-referenced with project documentation.*

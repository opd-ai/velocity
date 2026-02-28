# Implementation Gaps

## Procedural Generation Mandate

**All documentation in this repository enforces a strict requirement: 100% of gameplay assets—including audio, visual, and narrative components—must be procedurally generated at runtime using deterministic algorithms.** No pre-rendered, embedded, or bundled audio files (.mp3, .wav, .ogg), visual/image files (.png, .jpg, .svg, .gif), or static narrative content (hardcoded dialogue, pre-written cutscene scripts, fixed story arcs, embedded text assets) are permitted. All gaps listed below must be resolved in ways that adhere to this constraint.

## v1.0 — Core Engine + Playable Single-Player

### Procedural Audio Generation
- **Gap**: The roadmap specifies "adaptive music" with "intensity-driven music layers" and "genre-specific instrumentation", but no technical specification exists for how to generate music procedurally using Ebitengine's PCM audio streaming API. There is no note system, beat sequencer, or frequency table defined. **All audio must be procedurally generated—no audio files are permitted.**
- **Impact**: Audio implementation (PLAN.md Steps 13–15) cannot produce musically coherent output without a defined synthesis approach. Risk of placeholder silence or noise.
- **Resolution needed**: Define a minimal music synthesis strategy — e.g., predefined frequency tables for pentatonic scales, tempo-driven note sequencing, and layered sine/square waveforms. Alternatively, identify a Go-native procedural music library compatible with Ebitengine's audio context. **The solution must generate all audio procedurally at runtime with no external audio files.**

### Procedural Sprite Generation Algorithm
- **Gap**: The roadmap requires "procedurally drawn ship/enemy/projectile sprites" with "no external assets", but no algorithm is specified for generating visually recognisable ship shapes from a seed. The `SpriteCache` infrastructure exists but `rendering.go` contains no drawing logic. **All sprites must be procedurally generated—no image files are permitted.**
- **Impact**: Sprite generation (PLAN.md Step 4) requires an algorithmic approach to produce genre-appropriate ship silhouettes. Without this, rendered entities will be placeholder rectangles.
- **Resolution needed**: Define a sprite generation algorithm. Recommended approach: seed-driven symmetric pixel placement on an N×N grid (e.g., 16×16), where the left half is randomly filled and mirrored to the right to produce spacecraft-like silhouettes. Genre colour palettes are already available via `pkg/procgen/genre.Preset`. **The solution must generate all sprites procedurally at runtime with no external image files.**

### Gamepad and Touch Input Support
- **Gap**: The v1.0 roadmap lists "Keyboard/gamepad/touch; rebindable controls" as a Tier 1 requirement, but the current `InputState` and `config.yaml` only define keyboard bindings. Ebitengine supports `ebiten.GamepadButton` and `ebiten.TouchID` APIs, but no mapping scheme or configuration structure is defined for these input methods.
- **Impact**: v1.0 will only support keyboard input unless gamepad/touch mappings are implemented. This may be acceptable for an initial release targeting desktop platforms, but the roadmap lists it as a v1.0 deliverable.
- **Resolution needed**: Either (a) define gamepad button mappings and touch-region layouts in `config.yaml` and implement in `InputSystem`, or (b) formally defer gamepad/touch to v1.1 and update the roadmap accordingly.

### Tutorial First-Run Detection
- **Gap**: The roadmap specifies a "first-run guided wave teaching thrust, fire, and dodge" but does not define how to detect whether the player is running the game for the first time versus returning from a saved session.
- **Impact**: Tutorial (PLAN.md Step 18) may incorrectly re-trigger on every launch or never trigger if the heuristic is wrong.
- **Resolution needed**: Confirm the detection heuristic. Recommended: check for the existence of a save file or a `~/.velocity/firstrun` marker file. If no save file and no marker exist, trigger tutorial and create the marker on completion.

### Enemy AI Behaviour
- **Gap**: The v1.0 scope mentions a "Procedural wave generator (stub)" with "single enemy type, increasing count", but no enemy movement or attack behaviour is specified beyond the stub. The `pkg/engine` AI subsystem is listed as a v2.0 feature (behaviour trees), so v1.0 enemies need a simpler approach.
- **Impact**: Wave spawning (PLAN.md Step 11) needs a minimal AI for enemies to move toward the player or follow a path; without it, spawned enemies will be stationary.
- **Resolution needed**: Define a simple v1.0 enemy behaviour — e.g., linear approach toward the player position at a constant speed, with optional periodic firing. This keeps complexity low while enabling a playable combat loop.

### Score Persistence and Display
- **Gap**: The `saveload` package persists `Score int` as part of `RunState`, but there is no high-score tracking, per-seed leaderboard, or score display mechanism defined. The `ux.HUD` struct has a `Score` field but the accumulation logic (kills, combos) is not specified for v1.0.
- **Impact**: Score display (PLAN.md Step 16) will show a number, but the rules for earning points (per enemy type, combo multiplier) are not defined until v2.0's "Score / combo / multiplier system". v1.0 needs at least basic point-per-kill scoring.
- **Resolution needed**: Define v1.0 scoring as: +100 points per enemy destroyed, displayed in HUD. Combo system deferred to v2.0 per roadmap. High-score persistence deferred to v4.0 (Leaderboards).

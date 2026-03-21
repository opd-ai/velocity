# Project Overview

Velocity is a procedural arcade shooter in the Galaga × Asteroids style, built with Go and Ebitengine. It delivers pure thrust-and-fire action across five thematic universes (SciFi, Fantasy, Horror, Cyberpunk, Post-Apocalyptic), all from a single deterministic binary with **zero external asset files**. Every sprite, sound effect, music layer, and level is generated at runtime from algorithms and seeds.

The game targets casual and retro-arcade enthusiasts who want quick sessions with deep replayability. The seed-based determinism enables per-seed leaderboards — the same seed always produces the exact same wave sequences, enemy patterns, and spawns, making speedruns and score competitions meaningful.

Velocity is part of the **opd-ai procedural game suite** — a family of 8 Go+Ebiten games that share architectural patterns, coding conventions, and eventually shared library packages. All games in the suite follow the same philosophy: 100% procedural content, single-binary distribution, deterministic seed-based generation.

## Sibling Repository Context

Velocity shares code patterns and conventions with seven sibling repositories. When implementing features, check sibling repos for existing patterns:

| Repo | Genre | Description |
|------|-------|-------------|
| `opd-ai/venture` | Co-op action-RPG | Top-down co-op RPG with 30+ `pkg/` packages |
| `opd-ai/vania` | Metroidvania | Procedural platformer with `internal/` packages |
| `opd-ai/velocity` | Arcade shooter | This repo — Galaga-style shooter |
| `opd-ai/violence` | FPS | Raycasting first-person shooter with networking |
| `opd-ai/way` | Battle-cart racer | Kart racing game |
| `opd-ai/wyrm` | Survival RPG | First-person survival with crafting |
| `opd-ai/where` | Wilderness survival | Open-world survival |
| `opd-ai/whack` | Arena battle | Arena combat game |

All repos target eventual code extraction into shared libraries. Follow the patterns documented here to keep codebases compatible.

## Technical Stack

- **Primary Language**: Go 1.24.13
- **Game Framework**: Ebitengine v2.9.8 — 2D game engine with cross-platform + WASM support
- **Configuration**: Viper v1.21.0 — YAML-based configuration management
- **Build**: `go build ./cmd/velocity/`
- **Run**: `./velocity` (reads `config.yaml` from working directory)
- **Target Platforms**: Linux, macOS, Windows, WASM

### Key Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| `github.com/hajimehoshi/ebiten/v2` | v2.9.8 | 2D game engine, rendering, audio, input |
| `github.com/spf13/viper` | v1.21.0 | Configuration file parsing (YAML) |
| `golang.org/x/sync` | v0.17.0 | Concurrency utilities (indirect) |
| `golang.org/x/sys` | v0.36.0 | System calls (indirect, Ebiten dependency) |
| `golang.org/x/text` | v0.29.0 | Text processing (indirect) |

## Project Structure

Velocity follows the **velocity-style** layout: `cmd/` + `pkg/` with config file (YAML).

```
cmd/velocity/           Entry point — Ebitengine game loop (main.go)

pkg/
├── engine/             ECS framework, deterministic RNG, input, camera
├── config/             Viper-based configuration loading
├── rendering/          Sprite generation, animation, particles
├── audio/              Adaptive music, SFX, positional audio
├── ux/                 HUD, menus, tutorial scaffolding
├── procgen/            Procedural wave/content generation
│   └── genre/          Genre presets and post-processing parameters
├── combat/             Weapons, damage calculation, status effects
├── world/              Quest objectives, weather, loot, economy
├── class/              Ship hull classes, upgrade trees
├── balance/            Stat tuning tables, difficulty curves
├── companion/          Wingman AI behaviour trees
├── saveload/           Save/load serialization
├── networking/         Client-server netcode (v5.0+)
├── security/           E2E encryption, authentication (v5.0+)
├── social/             Squadrons, leaderboards, federation (v5.0+)
├── hostplay/           Local-host authoritative server (v5.0+)
├── integration/        External service hooks (OAuth, CDN, telemetry)
├── errors/             Structured error types
├── recovery/           Panic recovery middleware
├── validation/         Input and config validation
├── version/            Build version, save-file migration
├── audit/              Frame-time telemetry, entity logging
├── benchmark/          Per-system micro-benchmark harness
├── stability/          Crash detection, watchdog timer
└── visualtest/         Screenshot regression testing

mods/                   Mod loader and scripted mod API
config.yaml             Default configuration file
docs/                   Plan archives and execution guides
```

---

## ⚠️ CRITICAL: Complete Feature Integration (Zero Dangling Features)

**This is the single most important rule for this codebase.** Every feature, system, component, generator, and integration MUST be fully wired into the runtime. Dangling features are a maintenance burden, a source of frustration, and actively degrade code quality.

### The Dangling Feature Problem

In complex procedural game codebases, it is extremely common for features to be:
1. **Defined but never instantiated** — A system struct exists but is never created in `main()` or registered with the World
2. **Instantiated but never integrated** — A system runs but its output is never consumed by other systems
3. **Partially integrated** — A system works for one genre but silently no-ops for others
4. **Tested in isolation but broken in context** — Unit tests pass but the system was never wired into the game loop

### Current Known Gaps (Reference GAPS.md)

The following features are documented as incomplete — when working on these areas, ensure full integration:

- **Procedural Audio Generation**: `pkg/audio/audio.go` has stub methods for `PlaySFX()`, `PlayMusic()`, and `Update()` — no actual synthesis is implemented. GAPS.md specifies the need for frequency tables, beat sequencers, and Ebitengine PCM integration.
- **Procedural Sprite Generation**: `pkg/rendering/rendering.go` has `SpriteCache` but no drawing logic. GAPS.md specifies symmetric pixel-mirroring algorithm for ship silhouettes.
- **Gamepad/Touch Input**: `pkg/engine/engine.go` `InputState` only supports keyboard. Gamepad/touch mappings are not implemented.
- **Tutorial First-Run Detection**: `pkg/ux/ux.go` `Tutorial` exists but no detection heuristic for first-run vs returning player.
- **Enemy AI Behaviour**: No movement or attack logic for spawned enemies beyond the stub.
- **Score Accumulation Logic**: Points per kill, combo rules not defined.

### Mandatory Checks Before Adding or Modifying Any Feature

**Before writing ANY new code, verify the full integration chain:**

1. **Definition → Instantiation**: Is the struct/system created at runtime? Trace from `main()` through `NewGame()` → `engine.NewWorld()` → `AddSystem()`.
2. **Instantiation → Registration**: Is the system registered with the World? Check `world.AddSystem(s)` calls.
3. **Registration → Update Loop**: Does the system's `Update()` get called? Follow `Game.Update()` → `world.Update(dt)` → system iteration.
4. **Update → Output**: Does the system produce outputs (components, events, state changes) that other systems consume?
5. **Output → Consumer**: Is there at least one other system that reads this system's output?
6. **Consumer → Player Effect**: Does the chain ultimately produce something visible, audible, or mechanically felt by the player?

If ANY link in this chain is missing, the feature is dangling. **Do not submit dangling features.**

### Specific Anti-Patterns to Reject

```go
// ❌ BAD: System defined but never added to the World
type WeatherSystem struct { ... }
func (w *WeatherSystem) Update(dt float64) { ... }
// ...but NewWeatherSystem() is never called in main() or init

// ✅ GOOD: System defined, instantiated, registered, and consuming/producing
weather := NewWeatherSystem(seed)
world.AddSystem(weather)
// AND other systems react to weather state:
// render.ApplyWeatherEffects(weather.Current())
// audio.PlayWeatherAmbience(weather.Current())
```

```go
// ❌ BAD: Generator produces content but nothing uses it
func (g *Generator) GenerateWave(waveNumber int) WaveConfig {
    return WaveConfig{ ... }
}
// ...but GenerateWave() is never called from the game loop

// ✅ GOOD: Generator output is consumed by spawner
waveConfig := g.GenerateWave(currentWave)
for i := 0; i < waveConfig.EnemyCount; i++ {
    world.CreateEntity()  // Actually spawn the enemies
}
```

```go
// ❌ BAD: Genre setter defined but only one genre ever used
func (r *Renderer) SetGenre(genreID string) { r.genreID = genreID }
// ...but only "scifi" is ever passed, and genre switching is never tested

// ✅ GOOD: All five genres are exercised and produce distinct output
for _, genre := range genre.All() {
    r.SetGenre(genre)
    // Verify visual differences
}
```

### Integration Verification Checklist (run before every PR)

- [ ] `grep -rn 'func New' --include='*.go' pkg/ | grep -v _test.go` — Every constructor has at least one non-test caller
- [ ] `grep -rn 'TODO\|FIXME\|HACK\|XXX' --include='*.go' .` — All TODOs are tracked in GAPS.md or ROADMAP.md
- [ ] No stub method bodies in non-test `.go` files without a GAPS.md entry
- [ ] Every interface in the project has at least one runtime (non-test) implementation
- [ ] Every procedural generator is reachable from the main game initialization path
- [ ] Seeds are propagated through the full generation chain (parent seed → child generators)
- [ ] All five genres (fantasy, scifi, horror, cyberpunk, postapoc) are reachable via configuration

### Genre Integration Requirements

Every system that produces visual, audio, or narrative output MUST implement `SetGenre(genreID string)` and produce **observably different** output for each of the five genres:

| Genre ID | Visual Characteristics |
|----------|----------------------|
| `fantasy` | Magical effects, constellation backgrounds, enchanted vessels |
| `scifi` | Lasers, nebulae, metallic ships, bloom effects |
| `horror` | Organic enemies, void fog, desaturated colors |
| `cyberpunk` | Neon glow, grid backgrounds, geometric shapes |
| `postapoc` | Dust grain, salvaged aesthetics, debris fields |

**Verification**: When adding any rendering, audio, or procedural generation code, manually test all five genres before marking complete.

---

## Networking Best Practices (MANDATORY for all Go network code)

### Interface-Only Network Types (Hard Constraint)

When declaring network variables, ALWAYS use interface types. This is a **non-negotiable project rule**.

| ❌ Never Use (Concrete Type) | ✅ Always Use (Interface Type) |
|------------------------------|-------------------------------|
| `*net.UDPAddr` | `net.Addr` |
| `*net.IPAddr` | `net.Addr` |
| `*net.TCPAddr` | `net.Addr` |
| `*net.UDPConn` | `net.PacketConn` |
| `*net.TCPConn` | `net.Conn` |
| `*net.TCPListener` | `net.Listener` |

```go
// ✅ GOOD: Interface types everywhere
var addr net.Addr
var conn net.PacketConn
var tcpConn net.Conn
var listener net.Listener

// ❌ BAD: Concrete types
var addr *net.UDPAddr
var conn *net.UDPConn
```

**Never use type switches or type assertions to convert from interface to concrete type:**

```go
// ❌ BAD: Type assertion to access concrete methods
if udpConn, ok := conn.(*net.UDPConn); ok {
    udpConn.ReadFromUDP(buf)
}

// ✅ GOOD: Use the interface methods directly
n, addr, err := conn.ReadFrom(buf)  // PacketConn interface method
```

### High-Latency Network Design (200–5000ms)

All multiplayer networking code (v5.0+) MUST function correctly under **200–5000ms round-trip latency**. These games target diverse network conditions including mobile data, satellite internet, and intercontinental connections.

#### Mandatory Design Principles

1. **Client-Side Prediction**: Simulate game state locally; reconcile with server state when it arrives. Never block the game loop waiting for a server response.

2. **State Interpolation / Extrapolation**: Interpolate remote entity positions between known states. Extrapolate using last-known velocity when packets are delayed.

3. **Jitter Buffers**: Buffer incoming state updates and play back at consistent rate. Design for ±500ms jitter tolerance minimum.

4. **Idempotent Messages**: Every network message must be safe to process multiple times.

5. **No Synchronous RPC in Game Loops**: All network I/O is asynchronous. Results consumed on next available frame.

6. **Graceful Degradation**: At 5000ms latency, the game must remain playable. Reduce update frequency, increase prediction windows, hide latency with animations.

7. **Timeout Tolerance**: Connection timeouts ≥10 seconds. Disconnect detection via heartbeat absence over sliding window (≥3 missed heartbeats), never single missed packet.

```go
// ❌ BAD: Tight timeout
conn.SetReadDeadline(time.Now().Add(1 * time.Second))

// ✅ GOOD: Generous timeout for high-latency
conn.SetReadDeadline(time.Now().Add(10 * time.Second))

// ❌ BAD: Blocking RPC in game loop
func (g *Game) Update() error {
    state, err := g.server.GetWorldState()  // Blocks!
    return nil
}

// ✅ GOOD: Async receive with interpolation
func (g *Game) Update() error {
    select {
    case state := <-g.stateChannel:
        g.interpolator.PushServerState(state)
    default:
        // No new state — continue with prediction
    }
    g.world = g.interpolator.GetInterpolatedState(time.Now())
    return nil
}
```

#### Latency Budget Allocation (per frame at 60 FPS = 16.6ms)
- **Input processing**: ≤1ms
- **Local simulation / prediction**: ≤4ms
- **State interpolation**: ≤1ms
- **Network send (non-blocking enqueue)**: ≤0.5ms
- **Rendering**: ≤10ms
- **Network I/O goroutines**: Run independently, never counted against frame budget

---

## Code Assistance Guidelines

### 1. Deterministic Procedural Generation

All content generation MUST be deterministic and seed-based. Given the same seed, the game MUST produce identical output across all platforms and runs.

```go
// ✅ GOOD: Explicit seed-based RNG via engine.DeterministicRNG()
rng := engine.DeterministicRNG(seed)
value := rng.Intn(100)

// ❌ BAD: Global rand (non-deterministic, not thread-safe)
value := rand.Intn(100)

// ❌ BAD: Time-based seeding in generation code
rng := rand.New(rand.NewSource(time.Now().UnixNano()))

// ✅ GOOD: Derived seeds for sub-generators
terrainSeed := seed ^ 0x54455252  // "TERR"
enemySeed := seed ^ 0x454E454D    // "ENEM"
terrainRNG := engine.DeterministicRNG(terrainSeed)
enemyRNG := engine.DeterministicRNG(enemySeed)
```

**Wave formula reference** (from PLAN.md): Wave N spawns `N + 2` enemies with health `10 + N*5` and speed `1.0 + N*0.1`.

### 2. ECS Architecture

Velocity uses a lightweight Entity-Component-System architecture:

- **Entity**: A unique identifier (`type Entity uint64`)
- **Component**: Pure data attached to entities by name string
- **System**: Logic that operates on the World each tick via `Update(dt float64)`

```go
// Creating entities and adding components
entity := world.CreateEntity()
world.AddComponent(entity, "position", &Position{X: 100, Y: 100})
world.AddComponent(entity, "velocity", &Velocity{VX: 10, VY: 0})

// Retrieving components
pos, ok := world.GetComponent(entity, "position")
if ok {
    position := pos.(*Position)
}

// Registering systems
world.AddSystem(&PhysicsSystem{})
world.AddSystem(&RenderSystem{})
```

**Rules**:
- Components are pure data. NO logic in components.
- Systems contain ALL game logic.
- Never store direct entity references — use Entity IDs.
- All systems that produce output implement `GenreSetter` interface: `SetGenre(genreID string)`.

### 3. GenreSetter Interface

Every system producing visual, audio, or content output MUST implement:

```go
type GenreSetter interface {
    SetGenre(genreID string)
}
```

Genre switching happens at runtime via configuration. Genre IDs: `fantasy`, `scifi`, `horror`, `cyberpunk`, `postapoc`.

### 4. Performance Requirements

- Target 60 FPS (fixed 60 TPS via Ebitengine)
- Physics uses `dt = 1.0/60.0` for integration
- Entity budget: ≤200 on screen for v1.0 (player + enemies + projectiles + particles)
- Client memory budget: <500MB
- Cache generated sprites — never regenerate same sprite twice per session
- Use object pooling for frequently allocated objects (particles, projectiles)
- Viewport culling: skip update/draw for off-screen entities

### 5. Zero External Assets

The single-binary philosophy means ALL content is generated at runtime:

- **Graphics**: Procedural pixels, shape primitives, noise functions
- **Audio**: Synthesized from oscillators, envelopes (sine waves for music, noise bursts for explosions, square waves for lasers)
- **Levels**: Wave sequencer with algorithmic enemy patterns
- **Ships/Enemies**: Seed-driven symmetric pixel placement

**Never add asset files** (PNG, WAV, OGG, JSON level files) to the repository.

### 6. Error Handling

```go
// ✅ GOOD: Return errors, handle at call site
func GenerateWave(seed int64, waveNumber int) (*WaveConfig, error) {
    if seed == 0 && waveNumber == 0 {
        return nil, errors.New("CONFIG", "wave requires valid seed or wave number")
    }
    // ...
}

// ❌ BAD: Panic in library/game code
func GenerateWave(seed int64) *WaveConfig {
    if seed == 0 {
        panic("zero seed")  // Never panic in game logic
    }
}

// ✅ GOOD: Structured errors using pkg/errors
return errors.Wrap("WAVE_GEN", "failed to spawn enemies", err)
```

Panics acceptable ONLY in `main()` for unrecoverable startup failures. Use `recovery.WithRecovery()` to wrap the game loop.

### 7. Configuration via Viper

All runtime settings live in `config.yaml`. Use `pkg/config.Load()` to access:

```go
cfg, err := config.Load()
genre := cfg.Gameplay.Genre       // "scifi", "fantasy", etc.
arenaMode := cfg.Gameplay.ArenaMode  // "wrap" or "bounded"
seed := cfg.Gameplay.Seed
width := cfg.Display.Width
```

**Never hardcode values** that should be configurable. Add new settings to `config.yaml` schema and `setDefaults()`.

### 8. Input Validation

Use `pkg/validation` for all user/config inputs:

```go
if err := validation.ValidateGenre(cfg.Gameplay.Genre); err != nil {
    return nil, err
}
if err := validation.ValidateArenaMode(cfg.Gameplay.ArenaMode); err != nil {
    return nil, err
}
```

---

## Cross-Repository Code Sharing Patterns

### Shared Pattern Catalog

When implementing features, follow these patterns so code can be extracted into shared packages later:

| Pattern | Package Convention | Current Status |
|---------|-------------------|----------------|
| ECS core (World, Entity, Component, System) | `pkg/engine/` | Implemented |
| Deterministic RNG | `engine.DeterministicRNG()` | Implemented |
| Genre presets | `pkg/procgen/genre/` | Implemented |
| Sprite cache | `pkg/rendering/` | Stub |
| Particle system | `pkg/rendering/` | Stub |
| Audio manager | `pkg/audio/` | Stub |
| Save/load | `pkg/saveload/` | Implemented |
| Config loading | `pkg/config/` | Implemented |
| Error types | `pkg/errors/` | Implemented |
| Validation | `pkg/validation/` | Implemented |
| Recovery middleware | `pkg/recovery/` | Implemented |
| Frame telemetry | `pkg/audit/` | Implemented |
| Benchmarking | `pkg/benchmark/` | Implemented |
| Watchdog | `pkg/stability/` | Implemented |

### Interface Signatures (Cross-Repo Compatibility)

These interface signatures MUST be identical across all opd-ai repos:

```go
// System interface — identical across all repos
type System interface {
    Update(dt float64)
}

// GenreSetter interface — identical across all repos
type GenreSetter interface {
    SetGenre(genreID string)
}
```

### When Adding a Feature That Exists in a Sibling Repo

1. Check the sibling repo's implementation first
2. Use the same package structure and naming conventions
3. Match interface signatures exactly
4. If the sibling implementation has known issues (check its GAPS.md), fix them
5. Document divergences in ROADMAP.md with a note about future convergence

---

## Quality Standards

### Testing Requirements

- **Coverage target**: ≥40% per package (≥30% for display/Ebiten-dependent packages)
- **Table-driven tests** for all business logic and generation functions
- **Benchmarks** for hot-path code (rendering, physics, generation)
- **Race detection**: All tests must pass under `go test -race ./...`

```go
func TestWaveGeneration(t *testing.T) {
    tests := []struct {
        name       string
        seed       int64
        waveNumber int
        wantCount  int
    }{
        {"wave 1", 12345, 1, 3},   // 1 + 2 = 3 enemies
        {"wave 5", 12345, 5, 7},   // 5 + 2 = 7 enemies
        {"wave 10", 12345, 10, 12},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            g := procgen.NewGenerator(tt.seed)
            config := g.GenerateWave(tt.waveNumber)
            if config.EnemyCount != tt.wantCount {
                t.Errorf("EnemyCount = %d, want %d", config.EnemyCount, tt.wantCount)
            }
        })
    }
}
```

### Code Review Quality Gates

- Build success: `go build ./cmd/velocity/`
- All tests pass: `go test ./...`
- Race-free: `go test -race ./...`
- Static analysis: `go vet ./...`
- Formatting: `go fmt ./...`
- No new TODO/FIXME without corresponding GAPS.md entry

### Documentation Requirements

- Every exported type and function has a godoc comment
- README.md stays in sync with CLI flags and features
- GAPS.md is updated when new gaps are discovered
- ROADMAP.md reflects current priorities

---

## Naming Conventions

| Category | Convention | Examples |
|----------|-----------|----------|
| Packages | lowercase, single-word | `engine`, `procgen`, `audio` |
| Files | snake_case | `terrain_generator.go`, `combat_system.go` |
| Types | PascalCase | `TerrainGenerator`, `CombatSystem` |
| Interfaces | PascalCase, `-er` for single-method | `Generator`, `Renderer`, `GenreSetter` |
| Constants | PascalCase exported, camelCase unexported | `SciFi`, `defaultTimeout` |
| Seeds | Always `int64`, parameter named `seed` | `func NewGenerator(seed int64)` |

---

## GAPS.md and ROADMAP.md Protocol

When Copilot identifies a potential gap:

1. Note it in your response
2. Suggest adding it to GAPS.md with severity (Critical/High/Medium/Low)
3. Include the file path and line number
4. Propose an actionable fix

**Severity Levels**:
- **Critical**: Blocks core gameplay (no playable game)
- **High**: Major feature missing (documented in PLAN.md v1.0)
- **Medium**: Polish or secondary feature (v2.0+)
- **Low**: Nice-to-have, optimization

**Current GAPS.md items** (reference when working in related areas):
1. Procedural audio generation — no synthesis implemented
2. Procedural sprite generation algorithm — no drawing logic
3. Gamepad and touch input support — only keyboard implemented
4. Tutorial first-run detection — no heuristic defined
5. Enemy AI behaviour — no movement/attack logic
6. Score persistence and display — no accumulation rules

---

## Build and Run Commands

```sh
# Build
go build ./cmd/velocity/

# Run (reads config.yaml from current directory)
./velocity

# Test
go test ./...

# Test with race detection
go test -race ./...

# Static analysis
go vet ./...

# Format
go fmt ./...

# Benchmark
go test -bench=. -benchmem ./pkg/...

# Cross-compile for WASM
GOOS=js GOARCH=wasm go build -o velocity.wasm ./cmd/velocity/
```

---

## Version Information

- **Current Version**: 0.0.0-dev (set at compile time via `-ldflags`)
- **Save Version**: 1 (for save-file format migration)
- **Target v1.0 Scope**: Core engine + playable single-player with `scifi` genre

Set build version:
```sh
go build -ldflags="-X github.com/opd-ai/velocity/pkg/version.Version=1.0.0" ./cmd/velocity/
```

---

## Physics and Game Loop

### Newtonian 2D Flight Model

Velocity uses a Newtonian physics model for ship movement:

```go
// Thrust applies acceleration along ship facing vector
acceleration := Vector2{
    X: math.Cos(rotation) * thrustForce,
    Y: math.Sin(rotation) * thrustForce,
}
velocity.X += acceleration.X * dt
velocity.Y += acceleration.Y * dt

// Drag reduces velocity over time (default coefficient: 0.98)
velocity.X *= dragCoefficient
velocity.Y *= dragCoefficient

// Position updates from velocity
position.X += velocity.X * dt
position.Y += velocity.Y * dt
```

### Arena Modes

Two arena modes are supported via `config.yaml`:

| Mode | Behavior |
|------|----------|
| `wrap` | Asteroids-style — entities crossing edge appear on opposite side |
| `bounded` | Entities bounce off screen edges (velocity component reversed) |

```go
// Wrap mode implementation
if position.X < 0 {
    position.X = screenWidth
}
if position.X > screenWidth {
    position.X = 0
}

// Bounded mode implementation
if position.X < 0 || position.X > screenWidth {
    velocity.X = -velocity.X
}
```

### Game Loop Integration

The main game loop in `cmd/velocity/main.go` follows Ebitengine's pattern:

```go
func (g *Game) Update() error {
    const dt = 1.0 / 60.0  // Fixed timestep
    g.camera.Update(dt)
    g.world.Update(dt)     // Runs all registered systems
    g.audio.Update()
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Background fill
    // Entity rendering
    // Particle rendering
    // HUD overlay
    // Menu overlay (if active)
}
```

**Critical**: All game logic goes in `Update()`. Never modify game state in `Draw()`.

---

## Combat System Details

### Weapon Types

```go
const (
    WeaponPrimary WeaponType = iota   // Fast, low damage
    WeaponSecondary                    // Slow, high damage
    WeaponMissile                      // Homing, area damage
    WeaponBomb                         // Screen-clearing, limited ammo
)
```

### Cooldown Management

Weapons track cooldowns internally:

```go
weapon := combat.NewWeapon(combat.WeaponPrimary, 10.0, 0.2)  // 10 damage, 0.2s cooldown
if weapon.CanFire() {
    weapon.Fire()
    // Spawn projectile entity
}
weapon.Update(dt)  // Called each frame to tick down cooldown
```

### Status Effects

Ships can have status effects applied:

| Effect | Gameplay Impact |
|--------|-----------------|
| Slowed | Engine damage — reduced max speed |
| EMP'd | Controls scrambled — input partially ignored |
| Hull Breach | Continuous damage over time |

### Collision Detection

v1.0 uses AABB (Axis-Aligned Bounding Box) collision:

```go
func CheckCollision(a, b *BoundingBox) bool {
    return a.X < b.X+b.Width &&
           a.X+a.Width > b.X &&
           a.Y < b.Y+b.Height &&
           a.Y+a.Height > b.Y
}
```

---

## Space Weather System (v3.0+)

Thirteen space weather types are planned, each with gameplay effects:

| Weather | Effect | Active Genres |
|---------|--------|---------------|
| Solar flare | Screen whiteout | scifi, postapoc |
| Ion storm | Thruster degradation | scifi |
| Nebula interference | Radar jamming | scifi, fantasy |
| Debris field | Collision hazard | postapoc, scifi |
| Void fog | Visibility reduction | horror |
| Data storm | HUD corruption | cyberpunk |
| Arcane tempest | Projectile deflection | fantasy |
| Dust clouds | Movement slowed | postapoc |
| Meteor shower | Random projectile rain | all |
| Radiation zone | Continuous damage | scifi, horror |
| Black-hole gravity well | Area gravity pull | scifi, horror |
| Pulsar sweep | Periodic damage wave | scifi |
| Comet trail | Speed boost corridor | fantasy, scifi |

Each weather type implements `SetGenre()` for genre-appropriate visuals.

---

## Ship Classes (v4.0+)

Four base hull classes with distinct playstyles:

| Class | Health | Speed | Armor | Slots |
|-------|--------|-------|-------|-------|
| Scout | 60 | 300 | 2 | 2 |
| Interceptor | 80 | 260 | 3 | 3 |
| Gunship | 120 | 180 | 6 | 4 |
| Carrier | 200 | 120 | 10 | 5 |

v4.0 adds 15 total hull classes; v4.0+ adds 20 prestige variants unlocked via cumulative score milestones.

---

## Procedural Sprite Generation Algorithm

When implementing sprite generation (GAPS.md item), use this seed-driven symmetric pixel approach:

```go
// Generate a 16x16 ship silhouette
func GenerateShipSprite(rng *rand.Rand, genreID string) *ebiten.Image {
    size := 16
    img := ebiten.NewImage(size, size)
    pixels := make([]byte, size*size*4)
    
    palette := genre.GetPreset(genreID).Colors
    
    // Fill left half randomly
    for y := 0; y < size; y++ {
        for x := 0; x < size/2; x++ {
            if rng.Float64() < 0.4 {  // 40% fill chance
                color := palette[rng.Intn(len(palette))]
                setPixel(pixels, x, y, size, color)
                // Mirror to right half
                setPixel(pixels, size-1-x, y, size, color)
            }
        }
    }
    
    img.WritePixels(pixels)
    return img
}
```

**Key principles**:
1. Always use passed RNG, never global
2. Mirror left to right for spacecraft silhouette
3. Genre affects color palette
4. Cache results in `SpriteCache` keyed by `genreID:variant`

---

## Procedural Audio Synthesis

When implementing audio (GAPS.md item), use Ebitengine's audio streaming:

```go
// Pentatonic scale frequencies for pleasing procedural melodies
var pentatonic = []float64{261.63, 293.66, 329.63, 392.00, 440.00}  // C4, D4, E4, G4, A4

// Generate a simple tone
func GenerateTone(frequency, duration, sampleRate float64) []byte {
    numSamples := int(duration * sampleRate)
    buf := make([]byte, numSamples*4)  // 16-bit stereo
    
    for i := 0; i < numSamples; i++ {
        t := float64(i) / sampleRate
        sample := math.Sin(2 * math.Pi * frequency * t)
        
        // Apply envelope to avoid clicks
        envelope := 1.0
        attackSamples := int(0.01 * sampleRate)
        releaseSamples := int(0.01 * sampleRate)
        if i < attackSamples {
            envelope = float64(i) / float64(attackSamples)
        }
        if i > numSamples-releaseSamples {
            envelope = float64(numSamples-i) / float64(releaseSamples)
        }
        
        sample *= envelope
        // Convert to 16-bit signed int
        intSample := int16(sample * 32767)
        binary.LittleEndian.PutUint16(buf[i*4:], uint16(intSample))    // Left
        binary.LittleEndian.PutUint16(buf[i*4+2:], uint16(intSample))  // Right
    }
    return buf
}
```

**SFX patterns**:
- Laser fire: High-frequency square wave, short decay
- Explosion: White noise burst with exponential decay
- Powerup: Rising frequency sweep
- Menu select: Short sine pip

**Music layers** (intensity-driven):
- Base layer: Slow tempo, ambient pad
- Combat layer: Faster tempo, percussion
- Boss layer: Intense, full instrumentation

---

## Save/Load System

Save state structure:

```go
type RunState struct {
    Version    int    `json:"version"`     // For migration
    Seed       int64  `json:"seed"`        // Reproducibility
    Genre      string `json:"genre"`       // Current genre
    Wave       int    `json:"wave"`        // Progress
    Score      int64  `json:"score"`       // Points
    PlayerData []byte `json:"player_data"` // Serialized player state
}
```

**Save triggers**:
- Pause menu opened
- Application loses focus
- Wave completed (checkpoint)

**Load behavior**:
- Check save file exists on startup
- Validate version, migrate if needed
- Prompt "Continue" or "New Game"

---

## Stability and Recovery

### Watchdog Timer

The `pkg/stability` watchdog detects stuck frames:

```go
watchdog := stability.NewWatchdog(500 * time.Millisecond)
watchdog.Start()

// In game loop
func (g *Game) Update() error {
    watchdog.Ping()  // Reset timer each frame
    // ... game logic
    return nil
}
```

If a frame exceeds the timeout, the watchdog logs a warning and can trigger auto-save.

### Panic Recovery

The game loop is wrapped with `recovery.WithRecovery()`:

```go
recovery.WithRecovery(func() {
    if err := ebiten.RunGame(game); err != nil {
        os.Exit(1)
    }
}, recovery.DefaultHandler)
```

This prevents crashes from killing the process without cleanup.

---

## Mod Support (v5.0+)

The `mods/` package provides a mod loader:

```go
loader := mods.NewLoader()
loader.Register("custom-ships", "1.0.0")
loader.Register("extra-weapons", "2.1.0")

for _, mod := range loader.List() {
    if mod.Enabled {
        // Load mod content
    }
}
```

Mods can provide:
- Custom ship sprites
- New weapon types
- Enemy variants
- Wave scripts
- Genre themes

---

## Frame Timing and Telemetry

Use `pkg/audit` for performance tracking:

```go
logger := audit.NewLogger()

// In game loop
start := time.Now()
// ... frame logic
duration := time.Since(start).Seconds() * 1000
logger.Record(duration, world.EntityCount())

// Export for analysis
samples := logger.Samples()
```

**Performance targets**:
- Frame time: <16.6ms (60 FPS)
- Entity count: ≤200 for v1.0
- Memory: <500MB

---

## Multiplayer Architecture (v5.0+)

### Client-Server Model

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   Client    │◄────►│   Server    │◄────►│   Client    │
│  (Player 1) │      │(Authoritative)     │  (Player 2) │
└─────────────┘      └─────────────┘      └─────────────┘
```

- Server runs authoritative simulation
- Clients predict locally, reconcile with server state
- UDP for game state (loss-tolerant)
- TCP for reliable messages (chat, scores)

### Host-Play Mode

Local player can host authoritative server:

```go
host := hostplay.NewHost(7777, 4)  // Port 7777, max 4 players
host.Start()
// Game loop runs on host
// Other players connect via networking.Client
```

---

## Integration Hooks

The `pkg/integration` registry manages external services:

```go
registry := integration.NewRegistry()
registry.Register("oauth")
registry.Register("telemetry")
registry.Register("cdn")

if hook, ok := registry.Get("telemetry"); ok && hook.Enabled {
    // Send telemetry data
}
```

**Planned integrations** (v5.0+):
- OAuth identity provider
- CDN for leaderboard data
- Analytics/telemetry pipeline

---

## Common Pitfalls

### ❌ Modifying state in Draw()

```go
// BAD: Game state modified during rendering
func (g *Game) Draw(screen *ebiten.Image) {
    g.score++  // Never do this!
}
```

### ❌ Using global rand

```go
// BAD: Non-deterministic, not thread-safe
value := rand.Intn(100)

// GOOD: Use engine.DeterministicRNG
rng := engine.DeterministicRNG(seed)
value := rng.Intn(100)
```

### ❌ Forgetting to register systems

```go
// BAD: System created but never registered
physics := NewPhysicsSystem()
// Missing: world.AddSystem(physics)
```

### ❌ Blocking in Update()

```go
// BAD: Network call blocks game loop
func (g *Game) Update() error {
    data, _ := http.Get("http://...")  // Blocks!
    return nil
}

// GOOD: Async with channel
func (g *Game) Update() error {
    select {
    case data := <-g.dataChannel:
        g.processData(data)
    default:
    }
    return nil
}
```

### ❌ Hardcoding genre-specific logic

```go
// BAD: Hardcoded genre check
if genre == "scifi" {
    color = Blue
} else if genre == "fantasy" {
    color = Gold
}

// GOOD: Use genre presets
preset := genre.GetPreset(genreID)
color = preset.PrimaryColor
```

---

## Quick Reference Card

| Task | Command/Pattern |
|------|-----------------|
| Build | `go build ./cmd/velocity/` |
| Run | `./velocity` |
| Test | `go test ./...` |
| Race test | `go test -race ./...` |
| Benchmark | `go test -bench=. ./pkg/...` |
| Create entity | `entity := world.CreateEntity()` |
| Add component | `world.AddComponent(entity, "name", comp)` |
| Register system | `world.AddSystem(sys)` |
| Deterministic RNG | `rng := engine.DeterministicRNG(seed)` |
| Set genre | `system.SetGenre(genreID)` |
| Genre list | `genre.All()` → `["fantasy", "scifi", "horror", "cyberpunk", "postapoc"]` |
| Save state | `saveload.Save(path, &state)` |
| Load state | `state, err := saveload.Load(path)` |
| Wrap error | `errors.Wrap("CODE", "message", err)` |
| Validate genre | `validation.ValidateGenre(g)` |

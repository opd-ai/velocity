# VELOCITY

**Gameplay Style:** Arcade Shooter — Galaga × Asteroids with procedural content  
**Vision:** Ship a fully procedural, genre-skinned top-down arcade shooter — pure thrust-and-fire action across five thematic universes, from a single deterministic binary with no external assets.

---

## Procedural Generation Mandate

**All gameplay assets—including audio, visual, and narrative components—must be 100% procedurally generated at runtime using deterministic algorithms.** This is a fundamental architectural constraint:

- **Audio**: All music, sound effects, and positional audio must be generated procedurally using waveform synthesis. No .mp3, .wav, .ogg, or other audio files are permitted.
- **Visual**: All sprites, particles, animations, and UI elements must be generated procedurally from seed values. No .png, .jpg, .svg, .gif, or other image files are permitted.
- **Narrative**: All quests, dialogue, lore, world-building text, plot progression, and character backstories must be generated procedurally and deterministically. No hardcoded dialogue trees, pre-written cutscene scripts, fixed story arcs, or embedded text assets are permitted.

The entire game must produce identical output given identical seed inputs, with zero reliance on external, bundled, or pre-authored media or text files.

---

## Genre Support

Every system that produces visual, audio, or narrative output must implement `SetGenre(genreID string)` to switch thematic presentation at runtime. The five supported setting genres and their velocity-specific manifestations are:

| Genre ID    | Thematic Skin                 | Ships & Enemies                        | Weapons                        | Background & Environment                       |
|-------------|-------------------------------|----------------------------------------|--------------------------------|------------------------------------------------|
| `fantasy`   | Magical scrollspace           | Enchanted vessels, arcane constructs   | Spell projectiles, rune bursts | Constellation fields, celestial nebulae        |
| `scifi`     | Asteroid field                | Fighters, drones, mechanical cruisers  | Lasers, plasma cannons         | Deep-space nebulae, asteroid belts             |
| `horror`    | Void / biomass                | Organic ships, tentacle enemies        | Acid sprays, spore blasts      | Flesh-wall backgrounds, darkness, void fog     |
| `cyberpunk` | Data highway                  | Geometric neon ships, data constructs  | Data-stream bullets, EMP bursts| Grid-line backgrounds, neon city-scapes        |
| `postapoc`  | Orbital debris field          | Salvaged craft, makeshift vessels      | Makeshift cannons, debris shot | Earth-ruin orbit backgrounds, debris clouds    |

---

## Phased Milestones

### v1.0 — Core Engine + Playable Single-Player

> Foundation: deterministic engine, one playable genre (`scifi`), ship physics, basic combat loop.

**Tier 1 — Core Engine (direct port from venture):**

| System | Velocity Implementation | venture Source |
|--------|------------------------|----------------|
| ECS framework | Entity-component-system core; all game objects are entities | `pkg/engine` |
| `SetGenre()` interface | Required on every system producing output; switches skin at runtime | `pkg/engine` / `pkg/procgen/genre` |
| Seed-based deterministic RNG | All procedural content seeded — reproducible runs, per-seed leaderboards | `pkg/engine` |
| Input system | Keyboard/gamepad/touch; rebindable controls | `pkg/engine` |
| Camera system | Viewport tracking, screen-shake on explosions | `pkg/engine` |
| Rendering — sprites | Procedurally drawn ship/enemy/projectile sprites (no external image assets) | `pkg/rendering` |
| Rendering — animation | Frame-based sprite animation for thrusters, explosions, pickups | `pkg/rendering` |
| Rendering — particles | Explosion debris, thruster trails, impact sparks | `pkg/rendering` |
| Audio — adaptive music | Procedurally generated intensity-driven music layers; genre-specific instrumentation via waveform synthesis | `pkg/audio` |
| Audio — SFX | Procedurally generated weapon fire, explosions, powerup collect, UI feedback | `pkg/audio` |
| Audio — positional | Off-screen enemy audio cues using distance/angle (procedurally generated) | `pkg/audio` |
| UI / HUD | Ship health bar, shield meter, score, wave counter, combo display | `pkg/ux` |
| Menus | Main menu, pause menu, game-over screen, high-score entry | `pkg/ux` |
| Tutorial | First-run guided wave teaching thrust, fire, and dodge | `pkg/ux` |
| Save / load | Run state persistence; resume interrupted sessions | `pkg/saveload` |
| Config / settings | Resolution, audio levels, control bindings, genre preference | `pkg/config` |
| Performance — viewport culling | Skip update/draw for off-screen entities | `pkg/rendering` |
| Performance — batching | Batch draw calls per entity type per frame | `pkg/rendering` |
| Performance — sprite cache | Cache generated sprite bitmaps keyed by genre + variant | `pkg/rendering` |
| Error handling | Structured error types, panic recovery, graceful degradation | `pkg/errors` / `pkg/recovery` |
| Validation | Input validation for config, save data, network messages | `pkg/validation` |
| Version system | Embedded version string, save-file version migration | `pkg/version` |
| Benchmark harness | Per-system frame-time profiling | `pkg/benchmark` |
| CI/CD | Multi-platform GitHub Actions: Linux, macOS, Windows, WASM | — |
| Cross-platform builds | `GOOS` matrix; single binary, no external assets | — |
| Docs | README, CONTROLS, CHANGELOG, FAQ | — |
| 82%+ test coverage | Unit + integration tests for all packages | — |

**Velocity-Specific (new — not in venture):**

| System | Description |
|--------|-------------|
| Ship physics | Thrust, rotation, inertia, drag, mass — Newtonian 2-D flight model |
| Screen-wrap / bounded arena | Configurable: objects wrap at edges (Asteroids mode) or bounce (bounded mode) |
| Procedural wave generator (stub) | Seed-driven wave sequencer; single enemy type, increasing count |

---

### v2.0 — Core Systems: Weapons, Wave AI, Score, All 5 Genres

> All five genres integrated; full combat loop; wave difficulty ramp; score/combo system live.

**Tier 1 — Engine Additions:**

| System | Velocity Implementation | venture Source |
|--------|------------------------|----------------|
| Genre post-processing presets | Per-genre screen filter: bloom (`scifi`), desaturation (`horror`), neon glow (`cyberpunk`), parchment (`fantasy`), dust grain (`postapoc`) | `pkg/procgen/genre` |
| Procedural level / wave generation | Full procedural wave + formation generator; enemy types, spawn patterns, escalating difficulty curve | `pkg/procgen` |

**Tier 2 — Core Gameplay (adapted from venture):**

| Venture System | Velocity Adaptation | venture Source |
|----------------|---------------------|----------------|
| Combat system | Ship weapons: primary fire, secondary fire, missiles, bombs; hit detection; damage calculation | `pkg/combat` |
| AI — behavior trees | Formation AI, dive-bomb patterns, strafing runs, flanking maneuvers | `pkg/engine` (AI subsystem) |
| Status effects | Ship debuffs: slowed (engine damage), EMP'd (controls scrambled), hull breach (continuous damage) | `pkg/combat` |
| Magic / spells | Special weapons / screen-clearing bombs; genre-skinned (spell nova / EMP burst / spore cloud / data wipe / debris burst) | `pkg/combat` |
| Quests / objectives | Procedurally generated wave objectives: "survive 60 s", "destroy all elites before swarm", bonus mission waves (no pre-written quest text) | `pkg/world` (quests subsystem) |

**Velocity-Specific (new):**

| System | Description |
|--------|-------------|
| Bullet-pattern engine | Configurable projectile patterns: spread, spiral, homing, burst, ring; `SetGenre()` skins bullet visuals |
| Wave / formation system | Procedural enemy wave sequencer with formation templates (V, pincer, swarm, column); per-genre enemy skin |
| Score / combo / multiplier system | Chained kills extend combo timer; multiplier tiers (1×→8×); high-score tracking per seed |
| Boss pattern stubs | First boss encounter skeleton with multi-phase state machine; full polish in v4.0 |

---

### v3.0 — Visual Polish: Lighting, Particles, Space Weather, Per-Genre Post-Processing

> The game looks and sounds as distinct as it plays across all five genres.

**Tier 1 — Engine Additions:**

| System | Velocity Implementation | venture Source |
|--------|------------------------|----------------|
| Dynamic lighting | Per-entity light sources: thruster glow, muzzle flash, explosion bloom, background starlight | `pkg/rendering` (lighting subsystem) |
| Weather system (13 types → space weather) | 13 distinct space-weather phenomena, one per row below; each implements `SetGenre()` for genre-specific visual presentation | `pkg/world` (weather subsystem) |

Space weather types:

| # | Weather Type | Gameplay Effect | Active Genre(s) |
|---|-------------|-----------------|-----------------|
| 1 | Solar flare | Visibility burst — temporary screen whiteout | `scifi`, `postapoc` |
| 2 | Ion storm | Controls interference — thruster response degraded | `scifi` |
| 3 | Nebula interference | Radar jamming — enemy markers hidden | `scifi`, `fantasy` |
| 4 | Debris field | Collision hazard — dense obstacle layer | `postapoc`, `scifi` |
| 5 | Void fog | Extreme darkness — visibility severely reduced | `horror` |
| 6 | Data storm | HUD corruption — score display scrambled temporarily | `cyberpunk` |
| 7 | Arcane tempest | Spell deflection — projectiles randomly redirected | `fantasy` |
| 8 | Dust clouds | Drag increase — movement slowed | `postapoc` |
| 9 | Meteor shower | Random projectile rain from top of screen | all genres |
| 10 | Radiation zone | Continuous hull damage while inside zone | `scifi`, `horror` |
| 11 | Black-hole gravity well | Area gravity pull — ships spiral if not thrusting | `scifi`, `horror` |
| 12 | Pulsar sweep | Periodic damage wave — brief safe window between pulses | `scifi` |
| 13 | Comet trail | Speed boost corridor — hazardous debris on the edges | `fantasy`, `scifi` |

| System | Velocity Implementation | venture Source |
|--------|------------------------|----------------|
| Environmental storytelling | Background narrative conveyed through procedurally generated environment: drifting wreckage, distant battles, ruins, algorithmically generated lore fragments in debris patterns (no static text files) | `pkg/world` |

**Tier 2 — Visual / Audio Enhancements:**

| System | Velocity Adaptation | venture Source |
|--------|---------------------|----------------|
| Enhanced sprite generation | Higher-resolution procedural sprites; animated engine exhaust, weapon glow, shield shimmer | `pkg/rendering` |
| Visual test harness | Automated per-genre screenshot regression tests | `pkg/visualtest` |
| Audit / telemetry | Frame-time audit log, entity-count telemetry for performance regression detection | `pkg/audit` |
| Stability monitoring | Crash/freeze detection; auto-save before crash; watchdog timer | `pkg/stability` |

**Tier 3 — Advanced (style-appropriate):**

| Venture System | Velocity Adaptation | venture Source |
|----------------|---------------------|----------------|
| Destructible environments | Destructible asteroids, space stations, orbital debris structures | `pkg/world` |
| Mini-games | Bonus rounds between waves: asteroid slalom, target practice, docking challenge | `pkg/world` (mini-game subsystem) |

---

### v4.0 — Gameplay Expansion: Upgrades, Bosses, Powerups, Leaderboards

> Deep run-to-run progression; replayable with meaningful choices; full boss encounters.

**Tier 2 — Progression & Economy (adapted from venture):**

| Venture System | Velocity Adaptation | venture Source |
|----------------|---------------------|----------------|
| Progression — XP / leveling | Score-multiplier unlocks + permanent ship upgrades earned between waves | `pkg/balance` |
| Inventory / items | Powerup slots: four active slots for temporary powerups + one passive slot | `pkg/world` (inventory subsystem) |
| Skill / talent trees | Ship upgrade tree: hull, engines, weapons, shields, special; visible between-wave screen | `pkg/class` |
| Character classes (15 base + 20 prestige) | Ship hulls with distinct base stats and playstyle: 15 hull classes (scout, interceptor, gunship, carrier, …) delivered in v4.0; 20 prestige variants unlocked by accumulating score milestones across multiple runs — prestige hulls are a long-term progression goal that extends into post-v4.0 patches | `pkg/class` |
| Loot / drops | Powerup drops from destroyed enemies and boss crates | `pkg/world` (loot subsystem) |
| Crafting | Ship customization between waves: swap weapon loadouts, apply salvaged upgrades | `pkg/world` (crafting subsystem) |
| Shops / economy | Between-wave shop: spend score credits on upgrades, repairs, weapon mods | `pkg/world` (economy subsystem) |
| Balance system | Stat balance for all hulls, weapons, enemy types, boss phases | `pkg/balance` |
| Companion AI | Wingman AI: CPU-controlled co-pilot that assists player, evades, and calls out threats | `pkg/companion` |
| World events | Boss waves, meteor showers, special events (elite invasion, bonus wave, genre event) | `pkg/world` |

**Velocity-Specific (new):**

| System | Description |
|--------|-------------|
| Powerup system | Temporary powerups: shields, triple-shot, speed boost, magnet, invincibility; permanent: hull plating, engine tuning, weapon charge |
| Leaderboards | Per-seed and global high-score tracking; local and server-backed; genre filter |
| Boss encounters | Multi-phase bullet-hell bosses; genre-skinned; unique attack patterns per phase; cinematic intro |

---

### v5.0+ — Multiplayer Co-op / Competitive, Social, Production Polish

> Co-op and competitive multiplayer; social layer; mod support; production-grade distribution.

**Tier 4 — Multiplayer & Social (scaled from venture):**

| Venture System | Velocity Adaptation | venture Source |
|----------------|---------------------|----------------|
| Client-server netcode (200–5000 ms) | Co-op (2–4 players) and competitive score-attack; lag-compensated input | `pkg/networking` |
| E2E encrypted chat | In-game squadron chat with end-to-end encryption | `pkg/security` / `pkg/social` |
| Guilds → squadrons | Persistent named squadrons; shared leaderboard entries | `pkg/social` |
| Territory control → sector leaderboards | Per-genre sector control leaderboards; weekly resets | `pkg/social` |
| Federation | Cross-instance player identity federation | `pkg/social` |
| Cross-server travel → cross-server leaderboards | Compare scores across server regions / instances | `pkg/social` |
| Host-play | Local host runs authoritative game instance; others connect | `pkg/hostplay` |
| Integration layer | External service integration: OAuth identity, telemetry pipeline (no asset delivery—all assets are procedurally generated) | `pkg/integration` |

**Tier 6 — Production:**

| System | Description | venture Source |
|--------|-------------|----------------|
| Docker | Containerised server and CI build image | — |
| Release signing | Code-signing for macOS, Windows, and WASM | — |
| Mobile builds | iOS and Android via `ebitenmobile`; touch controls | — |
| Mod framework | Scripted mod API for custom ships, enemies, weapons, and wave scripts | `mods/` |
| Performance hardening | Load-testing, flame-graph guided optimisation at target 5000+ entity count | `pkg/benchmark` / `pkg/stability` |

---

## Excluded Features

The following venture systems are **not included** in velocity, with rationale:

| Venture System | Rationale |
|----------------|-----------|
| Vehicles | Player *is* the vehicle; a separate vehicle system is redundant. |
| Reputation / alignment | Arcade sessions are too short for faction standing to be meaningful. |
| Books / lore system | Detailed written lore breaks the fast-paced arcade flow; all lore must be procedurally generated if present. |
| Emotes | No persistent avatar; emotes have no context in a shooter. |
| Trading (player-to-player) | Economy scope limited to between-wave shop; P2P trading adds unwanted complexity. |
| Mail system | No persistent inbox context in an arcade session. |
| Fluid dynamics | No liquid in space; zero gameplay use. |
| Building / housing | Static arena play; no construction meta-game. |
| Furniture system | No housing, therefore no furniture. |
| VR system (`pkg/vr`) | Arcade top-down perspective is not suited to VR; deferred indefinitely. |
| Dialogue system | No named NPCs; narrative delivered through procedurally generated environment and HUD only (no dialogue trees or pre-scripted text). |

---

## Shared Infrastructure

The following venture packages are portable and will be vendored or imported directly into velocity with minimal adaptation:

| Package | Role in velocity |
|---------|-----------------|
| `pkg/engine` | ECS core, deterministic RNG, `SetGenre()` interface, game-loop scaffolding |
| `pkg/procgen/genre` | Genre post-processing presets, genre-keyed procedural asset generation helpers |
| `pkg/audio` | Adaptive music engine, SFX player, positional audio |
| `pkg/rendering` | Sprite generation, animation, particle system, dynamic lighting, batching, viewport culling |
| `pkg/networking` | Client-server netcode, connection management, lag compensation |
| `pkg/security` | E2E encrypted transport, authentication tokens |
| `pkg/social` | Squadron management, leaderboards, federation identity |
| `pkg/saveload` | Save/load serialisation, run-state snapshots |
| `pkg/config` | Settings schema, hot-reload |
| `pkg/balance` | Stat tuning tables, difficulty curves |
| `pkg/companion` | Wingman AI behaviour trees |
| `pkg/ux` | Menu framework, HUD components, tutorial scaffolding |
| `pkg/audit` | Frame-time telemetry, entity-count logging |
| `pkg/benchmark` | Per-system micro-benchmark harness |
| `pkg/stability` | Crash detection, watchdog, auto-save |
| `pkg/errors` | Structured error types |
| `pkg/recovery` | Panic recovery middleware |
| `pkg/validation` | Input and config validation |
| `pkg/version` | Embedded version + save-file migration |
| `pkg/hostplay` | Local-host authoritative server scaffold |
| `pkg/integration` | External service integration hooks |
| `mods/` | Mod loader and scripted mod API |

---

## Plan History

- 2026-02-27 PLAN.md created for v1.0 — Core Engine + Playable Single-Player (archived to docs/)

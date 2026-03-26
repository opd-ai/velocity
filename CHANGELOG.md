# Changelog

All notable changes to Velocity will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Core ECS engine with World, Entity, Component, and System architecture
- Newtonian 2D ship physics (thrust, rotation, inertia, drag)
- Keyboard input system with configurable controls
- Screen-wrap and bounded arena modes
- Procedural sprite generation framework with genre support
- Particle system for thruster trails and explosions
- Projectile system with AABB collision detection
- Weapon system with cooldowns and damage calculation
- Wave spawning system with procedural enemy generation
- Difficulty scaling (enemy health, speed, count per wave)
- Procedural audio synthesis (PCM-based SFX generation)
- Spatial audio with volume/pan based on entity position
- HUD framework (health, score, wave counter)
- Menu system (main menu, pause, game-over)
- Tutorial system framework
- Save/load system with JSON serialization
- Configuration via YAML (Viper)
- Error handling with structured error types
- Input validation for configuration values
- Panic recovery middleware
- Frame-time telemetry and audit logging
- Watchdog timer for stuck frame detection
- Benchmark harness for system performance testing
- Visual regression testing framework
- Ship hull class system (Scout, Interceptor, Gunship, Carrier)
- Upgrade tree system for ship progression
- Balance system with stat tables and difficulty curves
- Companion/wingman AI framework
- Weather system framework with genre-specific effects
- Quest objective system
- Economy system (credits, spending)
- Loot drop system
- Integration hooks for external services
- Networking client/server scaffolding
- Host-play mode for local authoritative server
- Security framework (token auth, encryption stubs)
- Social features (squadrons, leaderboards)
- Five genre presets: SciFi, Fantasy, Horror, Cyberpunk, Post-Apocalyptic
- Comprehensive unit test suite
- CI/CD pipeline with GitHub Actions
- Cross-platform build support (Linux, macOS, Windows, WASM)

### Technical Details
- Go 1.24+ required
- Ebitengine v2.9.8 game framework
- Viper v1.21.0 configuration management
- Fixed 60 TPS game loop
- Zero external asset files — all content procedurally generated
- Deterministic RNG for reproducible gameplay from seeds

## [0.0.0-dev] - Development

Initial development version. Not yet feature-complete.

---

## Version History

### Planned Releases

#### v1.0.0 — Core Engine + Playable Single-Player
- Complete ship physics and combat
- Working wave system with difficulty progression
- Full SciFi genre implementation
- Save/load functionality
- Basic HUD and menus

#### v2.0.0 — Visual Polish
- All five genres fully implemented
- Enhanced particle effects
- Improved sprite generation
- Procedural backgrounds

#### v3.0.0 — Weather & Economy
- Space weather system (13 weather types)
- Between-wave shop
- Ship upgrades

#### v4.0.0 — Ship Classes & Prestige
- 15 hull classes
- 20 prestige variants
- Unlock progression

#### v5.0.0 — Multiplayer
- Client-server netcode
- Co-op and PvP modes
- Leaderboards
- Squadron system

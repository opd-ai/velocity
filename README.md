# Velocity

Procedural arcade shooter built with Go and Ebitengine.

## Directory Structure

```
cmd/velocity/       Entry point (Ebitengine game loop)
pkg/
  engine/           ECS framework, deterministic RNG, input, camera
  config/           Viper-based configuration loading
  rendering/        Sprite generation, animation, particles
  audio/            Adaptive music, SFX, positional audio
  ux/               HUD, menus, tutorial
  procgen/          Procedural wave/content generation
  procgen/genre/    Genre presets and post-processing
  combat/           Weapons, damage, status effects
  world/            Quests, weather, environment, loot, economy
  class/            Ship hull classes, upgrade trees
  balance/          Stat tuning tables, difficulty curves
  companion/        Wingman AI behaviour trees
  saveload/         Save/load serialization
  networking/       Client-server netcode
  security/         E2E encryption, authentication
  social/           Squadrons, leaderboards, federation
  hostplay/         Local-host authoritative server
  integration/      External service hooks
  errors/           Structured error types
  recovery/         Panic recovery middleware
  validation/       Input and config validation
  version/          Build version, save-file migration
  audit/            Frame-time telemetry, entity logging
  benchmark/        Per-system micro-benchmark harness
  stability/        Crash detection, watchdog
  visualtest/       Screenshot regression testing
mods/               Mod loader and scripted mod API
config.yaml         Default configuration file
```

## Build and Run

```sh
go build ./cmd/velocity/
./velocity
```

## Configuration

The game reads `config.yaml` from the working directory on startup. Defaults are applied for any missing values. See `config.yaml` for the full schema covering display, audio, gameplay, and control settings.

## Dependencies

- [Ebitengine](https://ebitengine.org/) v2 — 2-D game engine
- [Viper](https://github.com/spf13/viper) — configuration management

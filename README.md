# Velocity

Procedural arcade shooter built with Go and Ebitengine.

## Core Design Principle: 100% Procedural Generation

**All gameplay assets—including audio, visual, and narrative components—must be procedurally generated at runtime using deterministic algorithms.** This project strictly prohibits:

- Pre-rendered or embedded audio files (e.g., .mp3, .wav, .ogg)
- Static visual/image files (e.g., .png, .jpg, .svg, .gif)
- Pre-authored narrative content (e.g., hardcoded dialogue, pre-written cutscene scripts, fixed story arcs, embedded text assets)

Every asset—from sprites and particles to music and sound effects, from quest descriptions to world-building lore—must be generated procedurally and deterministically, producing identical output given identical seed inputs, with zero reliance on external, bundled, or pre-authored media or text files.

## Directory Structure

```
cmd/velocity/       Entry point (Ebitengine game loop)
pkg/
  engine/           ECS framework, deterministic RNG, input, camera
  config/           Viper-based configuration loading
  rendering/        Procedurally generated sprite generation, animation, particles (no image files)
  audio/            Procedurally generated adaptive music, SFX, positional audio (no audio files)
  ux/               HUD, menus, tutorial
  procgen/          Procedural wave/content generation
  procgen/genre/    Genre presets and post-processing
  combat/           Weapons, damage, status effects
  world/            Procedurally generated quests, weather, environment, loot, economy (no text files)
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

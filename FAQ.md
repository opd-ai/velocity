# Frequently Asked Questions

## General

### What is Velocity?

Velocity is a procedural arcade shooter in the Galaga × Asteroids style. It features Newtonian physics, procedural content generation, and five thematic genres. Everything is generated at runtime — no external asset files are needed.

### What platforms does Velocity support?

Velocity supports:
- Linux (amd64)
- macOS (amd64, arm64)
- Windows (amd64)
- WebAssembly (browser)

### How do I run the game?

```bash
# Build
go build ./cmd/velocity/

# Run
./velocity
```

The game reads configuration from `config.yaml` in the current directory.

## Gameplay

### How does the physics work?

Velocity uses Newtonian 2D physics:
- **Thrust** accelerates your ship in the direction it's facing
- **No friction** — you'll keep drifting until you counter-thrust
- **Drag** slowly reduces velocity over time (configurable)
- **Rotation** is instant — rotate to aim, then thrust

### What are the arena modes?

- **Wrap mode**: Asteroids-style — crossing a screen edge teleports you to the opposite side
- **Bounded mode**: Hitting a wall bounces your ship

Set the mode in `config.yaml`:
```yaml
gameplay:
  arena_mode: wrap  # or "bounded"
```

### How do waves work?

- Wave N spawns `N + 2` enemies
- Enemy health: `10 + N*5`
- Enemy speed: `1.0 + N*0.1`
- Clear all enemies to advance to the next wave

### What are the genres?

Velocity supports five thematic genres that change visuals and audio:

| Genre | Description |
|-------|-------------|
| SciFi | Lasers, nebulae, metallic ships |
| Fantasy | Magical effects, enchanted vessels |
| Horror | Organic enemies, void fog |
| Cyberpunk | Neon glow, grid backgrounds |
| Post-Apocalyptic | Dust, salvaged aesthetics |

Change genre in `config.yaml`:
```yaml
gameplay:
  genre: scifi
```

## Technical

### Why are there no asset files?

Velocity follows a "zero external assets" philosophy:
- **Single binary distribution** — nothing to unpack or organize
- **Deterministic generation** — same seed = same content
- **Smaller download size**
- **No asset loading bugs**

All sprites, sounds, and levels are generated procedurally at runtime.

### What is seed-based determinism?

Given the same seed value, the game generates identical:
- Wave patterns
- Enemy spawns
- Procedural content

This enables:
- Reproducible speedruns
- Per-seed leaderboards
- Bug reproduction

Set a seed in `config.yaml`:
```yaml
gameplay:
  seed: 12345
```

### How do I report a bug?

Open an issue on GitHub with:
1. Operating system and version
2. Steps to reproduce
3. Expected vs actual behavior
4. The seed value if gameplay-related
5. Any error messages from the console

### How do I contribute?

1. Fork the repository
2. Create a feature branch
3. Follow the coding conventions in the custom instructions
4. Add tests for new functionality
5. Submit a pull request

## Configuration

### Where is the configuration file?

The game looks for `config.yaml` in the current working directory.

### What can I configure?

```yaml
display:
  width: 800
  height: 600
  fullscreen: false

gameplay:
  genre: scifi
  arena_mode: wrap
  seed: 0  # 0 = random seed

audio:
  master_volume: 1.0
  music_volume: 0.7
  sfx_volume: 1.0

controls:
  thrust: W
  rotate_left: A
  rotate_right: D
  fire_primary: Space
```

### How do I reset to defaults?

Delete `config.yaml` and restart the game. A new config with defaults will be created.

## Performance

### What are the system requirements?

Minimum:
- Any system that can run Go and OpenGL 2.1
- 512MB RAM
- 100MB disk space

Recommended:
- Modern CPU (2015+)
- 1GB RAM
- OpenGL 3.0+ capable GPU

### The game is running slowly. What can I do?

1. Lower the resolution in `config.yaml`
2. Run in windowed mode instead of fullscreen
3. Close other applications
4. Ensure you're running the native build, not WASM

### What's the entity budget?

v1.0 targets ≤200 entities on screen:
- 1 player
- ~50 enemies max
- ~100 projectiles
- ~50 particles

## Multiplayer (Future)

### When will multiplayer be available?

Multiplayer is planned for v5.0. Features will include:
- Client-server netcode
- Co-op mode
- PvP arena
- Leaderboards
- Squadron system

### Will there be cross-platform multiplayer?

Yes, all platforms will be able to play together.

## Modding (Future)

### Can I mod the game?

Mod support is planned for v5.0+. The `mods/` package provides the foundation for:
- Custom ships
- New weapons
- Enemy variants
- Genre themes
- Wave scripts

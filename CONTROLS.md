# Velocity Controls

## Keyboard Controls

### Movement
| Key | Action |
|-----|--------|
| W | Thrust forward |
| A | Rotate left |
| D | Rotate right |
| S | Brake (reduce velocity) |

### Combat
| Key | Action |
|-----|--------|
| Space | Fire primary weapon |
| Shift | Fire secondary weapon |

### UI
| Key | Action |
|-----|--------|
| Escape | Pause / Open menu |
| Enter | Confirm selection |
| Arrow Keys | Navigate menus |

## Gameplay Tips

### Physics
- Your ship follows Newtonian physics — thrust adds velocity, but there's no friction in space
- Use short thrust bursts followed by coasting for efficient movement
- Counter-thrust (rotate 180° and thrust) to slow down

### Combat
- Lead your targets — projectiles travel at a fixed speed
- Enemies spawn in waves from screen edges
- Destroy all enemies to advance to the next wave
- Each wave increases enemy count, health, and speed

### Scoring
- Points are awarded per enemy destroyed
- Combo multiplier increases with rapid kills
- Multiplier resets if no kills within the timeout window

## Arena Modes

### Wrap Mode (default)
- Crossing a screen edge teleports you to the opposite side
- Projectiles also wrap around

### Bounded Mode
- Hitting a screen edge bounces your ship
- Use walls tactically to change direction quickly

## Configuration

Control bindings can be customised in `config.yaml`:

```yaml
controls:
  thrust: W
  rotate_left: A
  rotate_right: D
  brake: S
  fire_primary: Space
  fire_secondary: Shift
  pause: Escape
```

## Gamepad Support (Planned)

Gamepad support is planned for a future release. The following mappings are anticipated:

| Input | Action |
|-------|--------|
| Left Stick | Rotate ship |
| Right Trigger | Thrust |
| A Button | Fire primary |
| B Button | Fire secondary |
| Start | Pause |

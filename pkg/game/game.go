//go:build !noebiten

// Package game provides the main Game struct and game loop implementation.
package game

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/opd-ai/velocity/pkg/audio"
	"github.com/opd-ai/velocity/pkg/combat"
	"github.com/opd-ai/velocity/pkg/config"
	"github.com/opd-ai/velocity/pkg/engine"
	"github.com/opd-ai/velocity/pkg/procgen"
	"github.com/opd-ai/velocity/pkg/rendering"
	"github.com/opd-ai/velocity/pkg/saveload"
	"github.com/opd-ai/velocity/pkg/ux"
	"github.com/opd-ai/velocity/pkg/version"
)

// Physics tuning constants - adjust these to change ship handling feel.
const (
	// DefaultThrustForce is the acceleration applied when thrusting.
	DefaultThrustForce = 200.0
	// DefaultRotationSpeed is the angular velocity when turning (radians/sec).
	DefaultRotationSpeed = 4.0
	// DefaultDragCoeff is the velocity dampening factor per frame (0-1).
	DefaultDragCoeff = 0.98
	// DefaultMaxSpeed is the maximum velocity magnitude.
	DefaultMaxSpeed = 300.0
)

// Gameplay balance constants.
const (
	// DefaultPlayerHealth is the starting health for the player ship.
	DefaultPlayerHealth = 100.0
	// DefaultPlayerWeaponDamage is the damage per projectile for primary weapon.
	DefaultPlayerWeaponDamage = 10.0
	// DefaultPlayerWeaponCooldown is the cooldown between shots in seconds.
	DefaultPlayerWeaponCooldown = 0.15
	// PlayerSpriteSizePx is the pixel size for player and enemy sprites.
	PlayerSpriteSizePx = 16
	// PlayerBoundingBoxOffset is the offset from sprite center for collision.
	PlayerBoundingBoxOffset = -8
	// DefaultProjectileSize is the pixel size for projectile sprites.
	DefaultProjectileSize = 8
)

// Scoring constants.
const (
	// BaseKillScore is the base points awarded for killing an enemy.
	BaseKillScore = 100
	// WaveBonusMultiplier is multiplied by wave number for wave completion bonus.
	WaveBonusMultiplier = 100
	// ComboTimerDuration is how long in seconds before combo resets.
	ComboTimerDuration = 2.0
	// ComboTierDivisor is the combo count divisor for multiplier calculation.
	ComboTierDivisor = 5
)

// Audio intensity levels.
const (
	// AudioIntensityLow is used between waves and at game start.
	AudioIntensityLow = 0.3
	// AudioIntensityHigh is used during active wave combat.
	AudioIntensityHigh = 0.8
)

// Frame timing.
const (
	// FrameDeltaTime is the fixed timestep for physics (1/60 second).
	FrameDeltaTime = 1.0 / 60.0
)

// UI layout constants.
const (
	// TutorialBannerY is the Y position for tutorial text.
	TutorialBannerY = 60
	// TutorialKeyHintOffset is the Y offset from banner for key hint.
	TutorialKeyHintOffset = 20
	// TutorialProgressOffset is the Y offset from bottom for progress text.
	TutorialProgressOffset = 60
	// CharacterWidthApprox is the approximate width of debug font characters.
	CharacterWidthApprox = 6
	// TutorialStepCount is the total number of tutorial steps.
	TutorialStepCount = 4
	// MenuItemSpacing is the vertical spacing between menu items.
	MenuItemSpacing = 20
	// HUDBottomOffset is the Y offset from bottom for HUD text.
	HUDBottomOffset = 20
	// GameOverScoreOffset is the Y offset for score display on game over.
	GameOverScoreOffset = 80
	// ViewportCullMargin is the margin in pixels for partial visibility culling.
	ViewportCullMargin = 32
)

// savePath returns the path to the save file.
func savePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "velocity_save.json"
	}
	dir := filepath.Join(home, ".velocity")
	_ = os.MkdirAll(dir, 0o755)
	return filepath.Join(dir, "save.json")
}

// Game implements the ebiten.Game interface.
type Game struct {
	cfg      *config.Config
	world    *engine.World
	camera   *engine.Camera
	renderer *rendering.Renderer
	audio    *audio.Manager
	hud      *ux.HUD
	menu     *ux.Menu

	// Game state management
	stateManager   *ux.GameStateManager
	menuController *ux.MenuController

	// Core systems
	physicsSystem    *engine.PhysicsSystem
	inputSystem      *engine.InputSystem
	arenaSystem      *engine.ArenaSystem
	projectileSystem *combat.ProjectileSystem
	damageSystem     *combat.DamageSystem
	weaponSystem     *combat.WeaponSystem
	enemyAISystem    *procgen.EnemyAISystem

	// Particle effects
	particleSystem *rendering.ParticleSystem

	// Procedural generation
	generator   *procgen.Generator
	waveSpawner *procgen.WaveSpawner
	waveManager *procgen.WaveManager

	// Player tracking
	playerEntity engine.Entity
	score        int64
	combo        int
	comboTimer   float64

	// Save state
	hasSavedGame bool

	// Tutorial system
	tutorial *ux.Tutorial

	// Sprite rendering cache (converts *image.RGBA to *ebiten.Image)
	ebitenImageCache map[string]*ebiten.Image
}

// NewGame initializes a new game instance from configuration.
func NewGame(cfg *config.Config) *Game {
	g := &Game{
		cfg:              cfg,
		world:            engine.NewWorld(),
		camera:           engine.NewCamera(),
		renderer:         rendering.NewRenderer(),
		audio:            audio.NewManager(),
		hud:              ux.NewHUD(),
		menu:             ux.NewMenu(),
		particleSystem:   rendering.NewParticleSystem(),
		ebitenImageCache: make(map[string]*ebiten.Image),
	}

	genre := cfg.Gameplay.Genre
	g.renderer.SetGenre(genre)
	g.renderer.SetSeed(cfg.Gameplay.Seed)
	g.audio.SetGenre(genre)
	g.hud.SetGenre(genre)
	g.menu.SetGenre(genre)
	g.particleSystem.SetGenre(genre)
	g.particleSystem.SetSeed(cfg.Gameplay.Seed)
	g.audio.SetVolumes(cfg.Audio.MasterVolume, cfg.Audio.MusicVolume, cfg.Audio.SFXVolume)

	// Check for saved game
	g.hasSavedGame = g.checkSavedGame()

	// Initialize game state management
	g.stateManager = ux.NewGameStateManager()
	g.menuController = ux.NewMenuController(g.stateManager)

	// Add Continue option if save exists
	if g.hasSavedGame {
		g.menuController.AddContinueOption()
	}

	// Set up state change callbacks
	g.stateManager.SetStateChangeCallback(func(from, to ux.GameState) {
		if to == ux.StatePlaying && from == ux.StateMainMenu {
			g.startNewGame()
		}
		if to == ux.StatePlaying && from == ux.StateGameOver {
			g.startNewGame()
		}
		// Save on pause
		if to == ux.StatePaused && from == ux.StatePlaying {
			g.saveGame()
		}
	})

	g.menuController.SetActionCallback(func(action string) {
		switch action {
		case "quit":
			g.saveGame() // Save before quitting
			os.Exit(0)
		case "continue":
			g.loadAndResumeGame()
		case "quit_menu":
			g.saveGame() // Save before returning to menu
		}
	})

	// Initialize core systems
	g.initializeSystems()

	return g
}

// initializeSystems sets up all game systems.
func (g *Game) initializeSystems() {
	width := g.cfg.Display.Width
	height := g.cfg.Display.Height

	// Physics system
	physicsConfig := engine.PhysicsConfig{
		ThrustForce:   DefaultThrustForce,
		RotationSpeed: DefaultRotationSpeed,
		DragCoeff:     DefaultDragCoeff,
		MaxSpeed:      DefaultMaxSpeed,
	}
	g.physicsSystem = engine.NewPhysicsSystem(g.world, physicsConfig)

	// Input system
	bindings := engine.DefaultKeyBindings()
	inputReader := engine.NewEbitenInputReader()
	g.inputSystem = engine.NewInputSystem(g.world, g.physicsSystem, bindings, inputReader)

	// Arena system
	arenaMode := engine.ArenaModeWrap
	if g.cfg.Gameplay.ArenaMode == "bounded" {
		arenaMode = engine.ArenaModeBounded
	}
	g.arenaSystem = engine.NewArenaSystem(g.world, width, height, arenaMode)

	// Combat systems
	g.projectileSystem = combat.NewProjectileSystem(g.world)
	g.damageSystem = combat.NewDamageSystem(g.world)
	g.weaponSystem = combat.NewWeaponSystem(g.world, g.projectileSystem)
	g.weaponSystem.SetFireProvider(g.inputSystem)

	// Connect projectile hits to damage system
	g.projectileSystem.SetHitCallback(func(projectile, target engine.Entity, damage float64) {
		g.damageSystem.QueueDamage(target, projectile, damage, "projectile")
	})

	// Procedural generation
	g.generator = procgen.NewGenerator(g.cfg.Gameplay.Seed)
	g.generator.SetGenre(g.cfg.Gameplay.Genre)

	g.enemyAISystem = procgen.NewEnemyAISystem(g.world)
	g.waveSpawner = procgen.NewWaveSpawner(g.world, g.generator, width, height)
	g.waveManager = procgen.NewWaveManager(g.world, g.waveSpawner, g.enemyAISystem)

	// Wave callbacks
	g.waveManager.SetWaveStartCallback(func(wave int) {
		g.audio.PlaySFX("wave_start")
	})

	g.waveManager.SetWaveCompleteCallback(func(wave int) {
		g.score += int64(wave * WaveBonusMultiplier)
		g.audio.PlaySFX("wave_complete")
	})

	// Death callbacks for scoring
	g.damageSystem.SetDeathCallback(func(event combat.DeathEvent) {
		// Check if it's an enemy
		if tag, ok := g.world.GetComponent(event.Entity, "collisiontag"); ok {
			if ct := tag.(*combat.CollisionTag); ct.Tag == "enemy" {
				g.onEnemyKilled(event.Entity)
			}
			if ct := tag.(*combat.CollisionTag); ct.Tag == "player" {
				g.onPlayerDeath()
			}
		}
	})

	// Register systems with world (for Update calls)
	g.world.AddSystem(g.physicsSystem)
	g.world.AddSystem(g.arenaSystem)
}

// startNewGame resets the game state and spawns the player.
func (g *Game) startNewGame() {
	// Clear all entities
	g.clearAllEntities()

	// Reset game state
	g.score = 0
	g.combo = 0
	g.comboTimer = 0
	g.waveManager.Reset()

	// Enable tutorial for first-run (no save file exists)
	if !g.hasSavedGame {
		g.tutorial = ux.NewTutorial()
	} else {
		g.tutorial = nil
	}

	// Spawn player at center
	g.spawnPlayer()

	// Start background music
	g.audio.PlayMusic()
	g.audio.SetIntensity(AudioIntensityLow)

	// Start first wave after a brief delay
	g.waveManager.StartNextWave()
}

// spawnPlayer creates the player entity at screen center.
func (g *Game) spawnPlayer() {
	width := float64(g.cfg.Display.Width)
	height := float64(g.cfg.Display.Height)

	g.playerEntity = g.world.CreateEntity()

	// Position at center
	g.world.AddComponent(g.playerEntity, "position", &engine.Position{
		X: width / 2,
		Y: height / 2,
	})

	// Initial velocity (stationary)
	g.world.AddComponent(g.playerEntity, "velocity", &engine.Velocity{VX: 0, VY: 0})

	// Facing up (negative Y in screen space)
	g.world.AddComponent(g.playerEntity, "rotation", &engine.Rotation{Angle: -math.Pi / 2})

	// Player health
	g.world.AddComponent(g.playerEntity, "health", &combat.Health{
		Current: DefaultPlayerHealth,
		Max:     DefaultPlayerHealth,
	})

	// Collision tag and bounding box
	g.world.AddComponent(g.playerEntity, "collisiontag", &combat.CollisionTag{Tag: "player"})
	g.world.AddComponent(g.playerEntity, "boundingbox", &combat.BoundingBox{
		X: PlayerBoundingBoxOffset, Y: PlayerBoundingBoxOffset, Width: PlayerSpriteSizePx, Height: PlayerSpriteSizePx,
	})

	// Weapon
	primaryWeapon := combat.NewWeapon(combat.WeaponPrimary, DefaultPlayerWeaponDamage, DefaultPlayerWeaponCooldown)
	g.world.AddComponent(g.playerEntity, "weapon", combat.NewWeaponComponent(primaryWeapon))

	// Sprite
	g.world.AddComponent(g.playerEntity, "sprite", &rendering.SpriteComponent{
		Type:    rendering.SpriteTypeShip,
		Variant: 0,
		Size:    PlayerSpriteSizePx,
	})

	// Connect player to systems
	g.inputSystem.SetPlayerEntity(g.playerEntity)
	g.enemyAISystem.SetPlayerEntity(g.playerEntity)
}

// clearAllEntities removes all entities from the world.
func (g *Game) clearAllEntities() {
	toRemove := []engine.Entity{}
	g.world.ForEachEntity(func(e engine.Entity) {
		toRemove = append(toRemove, e)
	})
	for _, e := range toRemove {
		g.world.RemoveEntity(e)
	}
	g.playerEntity = 0
}

// onEnemyKilled handles scoring when an enemy dies.
func (g *Game) onEnemyKilled(entity engine.Entity) {
	g.waveManager.OnEnemyKilled()

	// Get entity position for particle effect
	if pos, ok := g.world.GetComponent(entity, "position"); ok {
		p := pos.(*engine.Position)
		g.particleSystem.Emit(p.X, p.Y, 20)
	}

	// Mark tutorial kill action
	if g.tutorial != nil && g.tutorial.Active {
		g.tutorial.MarkAction("kill")
	}

	// Base score
	baseScore := int64(BaseKillScore)

	// Apply combo multiplier
	g.combo++
	g.comboTimer = ComboTimerDuration
	multiplier := 1 + g.combo/ComboTierDivisor
	g.score += baseScore * int64(multiplier)

	g.audio.PlaySFX("explosion")
}

// onPlayerDeath handles game over when player dies.
func (g *Game) onPlayerDeath() {
	g.stateManager.GameOver(g.score, g.waveManager.CurrentWave())
	g.deleteSaveFile() // Clear save on game over
}

// checkSavedGame returns true if a save file exists.
func (g *Game) checkSavedGame() bool {
	_, err := os.Stat(savePath())
	return err == nil
}

// saveGame persists the current run state to disk.
func (g *Game) saveGame() {
	if !g.stateManager.IsPlaying() && !g.stateManager.IsPaused() {
		return
	}

	// Get player health
	var playerHealth float64
	if h, ok := g.world.GetComponent(g.playerEntity, "health"); ok {
		playerHealth = h.(*combat.Health).Current
	}

	state := &saveload.RunState{
		Version:    1,
		Seed:       g.cfg.Gameplay.Seed,
		Genre:      g.cfg.Gameplay.Genre,
		Wave:       g.waveManager.CurrentWave(),
		Score:      g.score,
		PlayerData: encodePlayerHealth(playerHealth),
	}

	if err := saveload.Save(savePath(), state); err != nil {
		log.Printf("Warning: failed to save game: %v", err)
	}
	g.hasSavedGame = true
}

// loadAndResumeGame loads a saved game and resumes play.
func (g *Game) loadAndResumeGame() {
	state, err := saveload.Load(savePath())
	if err != nil {
		log.Printf("Warning: failed to load save: %v", err)
		g.startNewGame()
		return
	}

	// Clear and set up new game state
	g.clearAllEntities()
	g.score = state.Score
	g.combo = 0
	g.comboTimer = 0

	// Set wave state
	g.waveManager.Reset()
	for i := 0; i < state.Wave-1; i++ {
		g.waveManager.StartNextWave()
		// Clear spawned enemies immediately for skipped waves
		g.clearEnemies()
	}

	// Spawn player with saved health
	g.spawnPlayer()
	playerHealth := decodePlayerHealth(state.PlayerData)
	if h, ok := g.world.GetComponent(g.playerEntity, "health"); ok {
		h.(*combat.Health).Current = playerHealth
	}

	// Start current wave
	g.waveManager.StartNextWave()
	g.stateManager.StartGame()
}

// clearEnemies removes all enemy entities without scoring.
func (g *Game) clearEnemies() {
	toRemove := []engine.Entity{}
	g.world.ForEachEntity(func(e engine.Entity) {
		if tag, ok := g.world.GetComponent(e, "collisiontag"); ok {
			if ct := tag.(*combat.CollisionTag); ct.Tag == "enemy" {
				toRemove = append(toRemove, e)
			}
		}
	})
	for _, e := range toRemove {
		g.world.RemoveEntity(e)
	}
}

// deleteSaveFile removes the save file.
func (g *Game) deleteSaveFile() {
	_ = os.Remove(savePath())
	g.hasSavedGame = false
}

// encodePlayerHealth encodes player health as bytes.
func encodePlayerHealth(health float64) []byte {
	return []byte(fmt.Sprintf("%.2f", health))
}

// decodePlayerHealth decodes player health from bytes.
func decodePlayerHealth(data []byte) float64 {
	var health float64
	if len(data) > 0 {
		_, _ = fmt.Sscanf(string(data), "%f", &health)
	}
	if health <= 0 {
		health = DefaultPlayerHealth
	}
	return health
}

// Update advances the game state by one tick.
func (g *Game) Update() error {
	const dt = FrameDeltaTime

	// Handle menu input
	if g.stateManager.IsMenuActive() {
		g.handleMenuInput()
	}

	// Only update game systems if playing
	if g.stateManager.IsPlaying() {
		g.updateGameplay(dt)
	}

	// Handle pause toggle
	if g.inputSystem.IsPausePressed() && g.stateManager.IsPlaying() {
		g.stateManager.PauseGame()
	}

	g.audio.Update()
	return nil
}

// Key state tracking for edge detection
var prevKeys = make(map[ebiten.Key]bool)

// handleMenuInput processes menu navigation.
func (g *Game) handleMenuInput() {
	// Simple menu controls (keyboard)
	if ebiten.IsKeyPressed(ebiten.KeyUp) && !wasKeyPressed(ebiten.KeyUp) {
		g.menuController.MoveUp()
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && !wasKeyPressed(ebiten.KeyDown) {
		g.menuController.MoveDown()
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) && !wasKeyPressed(ebiten.KeyEnter) {
		g.menuController.Select()
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !wasKeyPressed(ebiten.KeySpace) {
		g.menuController.Select()
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) && g.stateManager.IsPaused() {
		g.stateManager.ResumeGame()
	}

	updatePrevKeys()
}

// wasKeyPressed returns true if the key was pressed in the previous frame.
func wasKeyPressed(key ebiten.Key) bool {
	return prevKeys[key]
}

// updatePrevKeys stores the current key states for edge detection.
func updatePrevKeys() {
	prevKeys[ebiten.KeyUp] = ebiten.IsKeyPressed(ebiten.KeyUp)
	prevKeys[ebiten.KeyDown] = ebiten.IsKeyPressed(ebiten.KeyDown)
	prevKeys[ebiten.KeyEnter] = ebiten.IsKeyPressed(ebiten.KeyEnter)
	prevKeys[ebiten.KeySpace] = ebiten.IsKeyPressed(ebiten.KeySpace)
	prevKeys[ebiten.KeyEscape] = ebiten.IsKeyPressed(ebiten.KeyEscape)
}

// updateGameplay runs all game systems for one tick.
func (g *Game) updateGameplay(dt float64) {
	// Update combo timer
	if g.comboTimer > 0 {
		g.comboTimer -= dt
		if g.comboTimer <= 0 {
			g.combo = 0
		}
	}

	// Update core systems
	g.camera.Update(dt)
	g.inputSystem.Update(dt)
	g.world.Update(dt) // Updates physics and arena systems

	// Track tutorial actions
	g.updateTutorialActions()

	// Combat systems
	g.weaponSystem.Update(dt)
	g.projectileSystem.Update(dt)
	g.damageSystem.Update(dt)

	// Update particle system
	g.particleSystem.Update(dt)

	// AI and wave management
	g.enemyAISystem.Update(dt)
	g.waveManager.Update(dt)

	// Update music intensity based on wave state
	if g.waveManager.WaveInProgress() {
		g.audio.SetIntensity(AudioIntensityHigh)
	} else {
		g.audio.SetIntensity(AudioIntensityLow)
	}

	// Auto-advance to next wave if current wave is complete
	if !g.waveManager.WaveInProgress() && g.stateManager.IsPlaying() {
		g.waveManager.StartNextWave()
	}

	// Update HUD
	health := 0.0
	if h, ok := g.world.GetComponent(g.playerEntity, "health"); ok {
		health = h.(*combat.Health).Current
	}
	g.hud.Update(health, 0, g.score, g.waveManager.CurrentWave(), g.combo)
}

// updateTutorialActions checks player input and marks tutorial progress.
func (g *Game) updateTutorialActions() {
	if g.tutorial == nil || !g.tutorial.Active {
		return
	}

	state := g.inputSystem.GetState()

	// Check actions based on current step
	if state.Thrust {
		g.tutorial.MarkAction("thrust")
	}
	if state.RotateLeft || state.RotateRight {
		g.tutorial.MarkAction("rotate")
	}
	if state.Fire {
		g.tutorial.MarkAction("fire")
	}
}

// Draw renders the current frame.
func (g *Game) Draw(screen *ebiten.Image) {
	// Background color based on genre
	bgColor := g.getBackgroundColor()
	screen.Fill(bgColor)

	// Draw gameplay elements if playing or paused
	if g.stateManager.IsPlaying() || g.stateManager.IsPaused() {
		g.drawGameplay(screen)
	}

	// Draw HUD if playing
	if g.stateManager.IsPlaying() {
		g.drawHUD(screen)
	}

	// Draw tutorial overlay if active
	if g.stateManager.IsPlaying() && g.tutorial != nil && g.tutorial.Active {
		g.drawTutorial(screen)
	}

	// Draw menus
	if g.stateManager.IsMenuActive() {
		g.drawMenu(screen)
	}

	// Display FPS and version
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Velocity %s | FPS: %.0f", version.GetVersion(), ebiten.ActualFPS()))
}

// getBackgroundColor returns the background color for the current genre.
func (g *Game) getBackgroundColor() color.RGBA {
	switch g.cfg.Gameplay.Genre {
	case "fantasy":
		return color.RGBA{R: 10, G: 5, B: 30, A: 255}
	case "horror":
		return color.RGBA{R: 10, G: 5, B: 10, A: 255}
	case "cyberpunk":
		return color.RGBA{R: 5, G: 0, B: 20, A: 255}
	case "postapoc":
		return color.RGBA{R: 20, G: 15, B: 10, A: 255}
	default: // scifi
		return color.RGBA{R: 0, G: 0, B: 20, A: 255}
	}
}

// drawGameplay renders all game entities using procedural sprites.
func (g *Game) drawGameplay(screen *ebiten.Image) {
	// Set up viewport for culling
	viewport := rendering.NewViewport(g.cfg.Display.Width, g.cfg.Display.Height)
	cullContext := rendering.NewCullContext(viewport, ViewportCullMargin)

	// Create draw batches for efficient rendering
	batches := rendering.CreateDrawBatches(g.world)
	batches = rendering.SortBatchesByRenderOrder(batches)

	// Render each batch
	for _, batch := range batches {
		for _, e := range batch.Entities {
			g.drawEntity(screen, e, cullContext)
		}
	}

	// Also render entities without sprite components (e.g., projectiles with projectile tag only)
	g.world.ForEachEntity(func(e engine.Entity) {
		// Skip entities that were already rendered via batches
		if _, hasSprite := g.world.GetComponent(e, "sprite"); hasSprite {
			return
		}
		// Only render projectiles that don't have sprite components
		if _, hasProjectile := g.world.GetComponent(e, "projectile"); hasProjectile {
			g.drawEntity(screen, e, cullContext)
		}
	})

	// Render particles
	g.drawParticles(screen)
}

// drawParticles renders all active particles.
func (g *Game) drawParticles(screen *ebiten.Image) {
	particles := g.particleSystem.GetParticles()
	for _, p := range particles {
		col := computeParticleColor(p)
		size := clampParticleSize(int(p.Size))
		g.renderParticlePixels(screen, p.X, p.Y, size, col)
	}
}

// computeParticleColor applies alpha fade based on remaining life.
func computeParticleColor(p rendering.Particle) color.Color {
	alpha := float32(p.Life / p.MaxLife)
	if alpha < 0 {
		alpha = 0
	}
	col := p.Color
	col.A = uint8(float32(col.A) * alpha)
	return col
}

// clampParticleSize ensures minimum particle size of 1 pixel.
func clampParticleSize(size int) int {
	if size < 1 {
		return 1
	}
	return size
}

// renderParticlePixels draws a filled square particle at the given position.
func (g *Game) renderParticlePixels(screen *ebiten.Image, px, py float64, size int, col color.Color) {
	halfSize := size / 2
	for dy := 0; dy < size; dy++ {
		for dx := 0; dx < size; dx++ {
			x := int(px) + dx - halfSize
			y := int(py) + dy - halfSize
			if g.isWithinScreen(x, y) {
				screen.Set(x, y, col)
			}
		}
	}
}

// isWithinScreen checks if coordinates are within screen bounds.
func (g *Game) isWithinScreen(x, y int) bool {
	return x >= 0 && x < g.cfg.Display.Width && y >= 0 && y < g.cfg.Display.Height
}

// drawEntity renders a single entity with its sprite and rotation.
func (g *Game) drawEntity(screen *ebiten.Image, e engine.Entity, cullContext *rendering.CullContext) {
	posComp, hasPos := g.world.GetComponent(e, "position")
	if !hasPos {
		return
	}
	pos := posComp.(*engine.Position)

	// Determine sprite size (default to PlayerSpriteSizePx for entities without sprite component)
	spriteSize := PlayerSpriteSizePx
	if spriteComp, hasSprite := g.world.GetComponent(e, "sprite"); hasSprite {
		spriteSize = spriteComp.(*rendering.SpriteComponent).Size
	}

	// Apply viewport culling
	if !cullContext.ShouldRender(pos.X, pos.Y, float64(spriteSize), float64(spriteSize)) {
		return
	}

	// Get rotation if available
	var angle float64
	if rotComp, hasRot := g.world.GetComponent(e, "rotation"); hasRot {
		angle = rotComp.(*engine.Rotation).Angle
	}

	// Get or generate the sprite image
	img := g.getSpriteImage(e, spriteSize)
	if img == nil {
		return
	}

	// Set up draw options with rotation and position
	opts := &ebiten.DrawImageOptions{}

	// Center the sprite for rotation
	halfSize := float64(spriteSize) / 2
	opts.GeoM.Translate(-halfSize, -halfSize)
	opts.GeoM.Rotate(angle)
	opts.GeoM.Translate(pos.X, pos.Y)

	screen.DrawImage(img, opts)
}

// getSpriteImage returns the ebiten.Image for an entity, generating and caching it if needed.
func (g *Game) getSpriteImage(e engine.Entity, defaultSize int) *ebiten.Image {
	cacheKey, rgbaImg := g.resolveEntitySprite(e)
	if cacheKey == "" {
		return nil
	}

	if img, ok := g.ebitenImageCache[cacheKey]; ok {
		return img
	}

	return g.cacheConvertedSprite(cacheKey, rgbaImg)
}

// resolveEntitySprite determines the cache key and generates the sprite for an entity.
func (g *Game) resolveEntitySprite(e engine.Entity) (string, *image.RGBA) {
	if spriteComp, hasSprite := g.world.GetComponent(e, "sprite"); hasSprite {
		return g.resolveSpriteComponent(spriteComp.(*rendering.SpriteComponent))
	}

	if _, hasProjectile := g.world.GetComponent(e, "projectile"); hasProjectile {
		return g.resolveDefaultProjectile()
	}

	return "", nil
}

// resolveSpriteComponent generates sprite based on the component type.
func (g *Game) resolveSpriteComponent(sprite *rendering.SpriteComponent) (string, *image.RGBA) {
	cacheKey := fmt.Sprintf("%s:%d:%d", g.renderer.GetGenre(), sprite.Type, sprite.Variant)
	rgbaImg := g.generateSpriteByType(sprite)
	return cacheKey, rgbaImg
}

// generateSpriteByType dispatches to the appropriate sprite generator.
func (g *Game) generateSpriteByType(sprite *rendering.SpriteComponent) *image.RGBA {
	switch sprite.Type {
	case rendering.SpriteTypeShip:
		return g.renderer.GetOrCreateShipSprite(sprite.Variant, sprite.Size)
	case rendering.SpriteTypeEnemy:
		return g.renderer.GetOrCreateEnemySprite(sprite.Variant, sprite.Size)
	case rendering.SpriteTypeProjectile:
		return g.renderer.GetOrCreateProjectileSprite(sprite.Variant, sprite.Size)
	default:
		return nil
	}
}

// resolveDefaultProjectile returns the default projectile sprite key and image.
func (g *Game) resolveDefaultProjectile() (string, *image.RGBA) {
	cacheKey := fmt.Sprintf("%s:projectile:0", g.renderer.GetGenre())
	rgbaImg := g.renderer.GetOrCreateProjectileSprite(0, DefaultProjectileSize)
	return cacheKey, rgbaImg
}

// cacheConvertedSprite converts RGBA to ebiten.Image and caches the result.
func (g *Game) cacheConvertedSprite(cacheKey string, rgbaImg *image.RGBA) *ebiten.Image {
	if rgbaImg == nil {
		return nil
	}
	ebitenImg := ebiten.NewImageFromImage(rgbaImg)
	g.ebitenImageCache[cacheKey] = ebitenImg
	return ebitenImg
}

// drawHUD renders the heads-up display.
func (g *Game) drawHUD(screen *ebiten.Image) {
	hudText := fmt.Sprintf("Score: %d | Wave: %d | Combo: x%d | Health: %.0f",
		g.hud.Score, g.hud.Wave, g.hud.Combo+1, g.hud.Health)
	ebitenutil.DebugPrintAt(screen, hudText, 10, g.cfg.Display.Height-HUDBottomOffset)
}

// drawTutorial renders the tutorial overlay.
func (g *Game) drawTutorial(screen *ebiten.Image) {
	prompt := g.tutorial.CurrentPrompt()
	if prompt == nil {
		return
	}

	width := g.cfg.Display.Width
	height := g.cfg.Display.Height

	// Draw semi-transparent banner at top of screen
	y := TutorialBannerY

	// Draw prompt text centered
	text := prompt.Text
	textWidth := len(text) * CharacterWidthApprox
	x := (width - textWidth) / 2
	ebitenutil.DebugPrintAt(screen, text, x, y)

	// Draw key hint below
	if prompt.KeyHint != "" {
		hintText := fmt.Sprintf("Press: %s", prompt.KeyHint)
		hintWidth := len(hintText) * CharacterWidthApprox
		ebitenutil.DebugPrintAt(screen, hintText, (width-hintWidth)/2, y+TutorialKeyHintOffset)
	}

	// Draw progress indicator
	progressText := fmt.Sprintf("Tutorial %d/%d", g.tutorial.Step+1, TutorialStepCount)
	ebitenutil.DebugPrintAt(screen, progressText, (width-len(progressText)*CharacterWidthApprox)/2, height-TutorialProgressOffset)
}

// drawMenu renders the current menu.
func (g *Game) drawMenu(screen *ebiten.Image) {
	items := g.menuController.GetCurrentItems()
	if items == nil {
		return
	}

	// Draw menu title
	title := g.getMenuTitle()
	width := g.cfg.Display.Width
	height := g.cfg.Display.Height

	ebitenutil.DebugPrintAt(screen, title, width/2-len(title)*(CharacterWidthApprox/2), height/3)

	// Draw menu items
	for i, item := range items {
		y := height/2 + i*MenuItemSpacing
		prefix := "  "
		if i == g.menuController.SelectionIndex() {
			prefix = "> "
		}
		ebitenutil.DebugPrintAt(screen, prefix+item.Label, width/2-40, y)
	}

	// Draw score on game over
	if g.stateManager.State() == ux.StateGameOver {
		scoreText := fmt.Sprintf("Final Score: %d | Wave Reached: %d",
			g.stateManager.FinalScore(), g.stateManager.FinalWave())
		ebitenutil.DebugPrintAt(screen, scoreText, width/2-len(scoreText)*(CharacterWidthApprox/2), height/2+GameOverScoreOffset)
	}
}

// getMenuTitle returns the title for the current menu state.
func (g *Game) getMenuTitle() string {
	switch g.stateManager.State() {
	case ux.StateMainMenu:
		return "VELOCITY"
	case ux.StatePaused:
		return "PAUSED"
	case ux.StateGameOver:
		return "GAME OVER"
	default:
		return ""
	}
}

// Layout returns the logical screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.cfg.Display.Width, g.cfg.Display.Height
}

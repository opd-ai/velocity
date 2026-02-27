// Velocity — a procedural arcade shooter built with Ebitengine.
package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/opd-ai/velocity/pkg/audio"
	"github.com/opd-ai/velocity/pkg/config"
	"github.com/opd-ai/velocity/pkg/engine"
	"github.com/opd-ai/velocity/pkg/recovery"
	"github.com/opd-ai/velocity/pkg/rendering"
	"github.com/opd-ai/velocity/pkg/ux"
	"github.com/opd-ai/velocity/pkg/version"
)

// Game implements the ebiten.Game interface.
type Game struct {
	cfg      *config.Config
	world    *engine.World
	camera   *engine.Camera
	renderer *rendering.Renderer
	audio    *audio.Manager
	hud      *ux.HUD
	menu     *ux.Menu
}

// NewGame initializes a new game instance from configuration.
func NewGame(cfg *config.Config) *Game {
	g := &Game{
		cfg:      cfg,
		world:    engine.NewWorld(),
		camera:   engine.NewCamera(),
		renderer: rendering.NewRenderer(),
		audio:    audio.NewManager(),
		hud:      ux.NewHUD(),
		menu:     ux.NewMenu(),
	}
	genre := cfg.Gameplay.Genre
	g.renderer.SetGenre(genre)
	g.audio.SetGenre(genre)
	g.hud.SetGenre(genre)
	g.menu.SetGenre(genre)
	g.audio.SetVolumes(cfg.Audio.MasterVolume, cfg.Audio.MusicVolume, cfg.Audio.SFXVolume)
	return g
}

// Update advances the game state by one tick.
func (g *Game) Update() error {
	const dt = 1.0 / 60.0
	g.camera.Update(dt)
	g.world.Update(dt)
	g.audio.Update()
	return nil
}

// Draw renders the current frame.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0, G: 0, B: 20, A: 255})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Velocity %s | Genre: %s", version.GetVersion(), g.cfg.Gameplay.Genre))
}

// Layout returns the logical screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.cfg.Display.Width, g.cfg.Display.Height
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	game := NewGame(cfg)

	ebiten.SetWindowSize(cfg.Display.Width, cfg.Display.Height)
	ebiten.SetWindowTitle("Velocity")
	ebiten.SetVsyncEnabled(cfg.Display.VSync)

	if cfg.Display.Fullscreen {
		ebiten.SetFullscreen(true)
	}

	recovery.WithRecovery(func() {
		if err := ebiten.RunGame(game); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}, recovery.DefaultHandler)
}

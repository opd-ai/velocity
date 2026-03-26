//go:build !noebiten

// Velocity — a procedural arcade shooter built with Ebitengine.
// This file is the entry point; see pkg/game for the main Game implementation.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/opd-ai/velocity/pkg/config"
	"github.com/opd-ai/velocity/pkg/game"
	"github.com/opd-ai/velocity/pkg/recovery"
)

// main initializes the game configuration and starts the Ebitengine game loop.
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	g := game.NewGame(cfg)

	ebiten.SetWindowSize(cfg.Display.Width, cfg.Display.Height)
	ebiten.SetWindowTitle("Velocity")
	ebiten.SetVsyncEnabled(cfg.Display.VSync)

	if cfg.Display.Fullscreen {
		ebiten.SetFullscreen(true)
	}

	recovery.WithRecovery(func() {
		if err := ebiten.RunGame(g); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}, recovery.DefaultHandler)
}

// Package config provides Viper-based configuration loading for the game.
package config

import (
	"github.com/spf13/viper"
)

// Config holds the game configuration.
type Config struct {
	Display  DisplayConfig  `mapstructure:"display"`
	Audio    AudioConfig    `mapstructure:"audio"`
	Gameplay GameplayConfig `mapstructure:"gameplay"`
	Controls ControlsConfig `mapstructure:"controls"`
}

// DisplayConfig holds display/window settings.
type DisplayConfig struct {
	Width      int  `mapstructure:"width"`
	Height     int  `mapstructure:"height"`
	Fullscreen bool `mapstructure:"fullscreen"`
	VSync      bool `mapstructure:"vsync"`
}

// AudioConfig holds audio volume settings.
type AudioConfig struct {
	MasterVolume float64 `mapstructure:"master_volume"`
	MusicVolume  float64 `mapstructure:"music_volume"`
	SFXVolume    float64 `mapstructure:"sfx_volume"`
}

// GameplayConfig holds gameplay settings.
type GameplayConfig struct {
	Genre     string `mapstructure:"genre"`
	ArenaMode string `mapstructure:"arena_mode"`
	Seed      int64  `mapstructure:"seed"`
}

// ControlsConfig holds key binding settings.
type ControlsConfig struct {
	Thrust      string `mapstructure:"thrust"`
	RotateLeft  string `mapstructure:"rotate_left"`
	RotateRight string `mapstructure:"rotate_right"`
	Fire        string `mapstructure:"fire"`
	Secondary   string `mapstructure:"secondary"`
	Pause       string `mapstructure:"pause"`
}

// Load reads the configuration file and returns a Config struct.
// It searches for config.yaml in the current directory.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("display.width", 800)
	viper.SetDefault("display.height", 600)
	viper.SetDefault("display.fullscreen", false)
	viper.SetDefault("display.vsync", true)

	viper.SetDefault("audio.master_volume", 0.8)
	viper.SetDefault("audio.music_volume", 0.6)
	viper.SetDefault("audio.sfx_volume", 0.8)

	viper.SetDefault("gameplay.genre", "scifi")
	viper.SetDefault("gameplay.arena_mode", "wrap")
	viper.SetDefault("gameplay.seed", 0)

	viper.SetDefault("controls.thrust", "W")
	viper.SetDefault("controls.rotate_left", "A")
	viper.SetDefault("controls.rotate_right", "D")
	viper.SetDefault("controls.fire", "Space")
	viper.SetDefault("controls.secondary", "Shift")
	viper.SetDefault("controls.pause", "Escape")
}

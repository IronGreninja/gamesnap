package internal

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Options struct {
	BackupMax   int8    `mapstructure:"backup_max"`
	Destination TPath_t `mapstructure:"destination"`
}

type Config struct {
	Games     map[string][]TPath_t `mapstructure:"games"`
	Opts      Options              `mapstructure:"options"`
	Variables map[string]string    `mapstructure:"variables,omitempty"`
}

type Path struct {
	Name string `mapstructure:"name" json:"name"`
	Type string `mapstructure:"type" json:"type"`
}

type BRConfig struct { // Backup Restore Config
	Game  string          `mapstructure:"game" json:"game"`
	Paths map[string]Path `mapstructure:"paths" json:"paths"`
}

var config *Config
var once sync.Once

func GetConfig() *Config {
	if config == nil {
		once.Do(func() {
			config = &Config{}
			if err := viper.UnmarshalExact(config); err != nil {
				fmt.Fprintln(os.Stderr, "Error: could not parse config")
				os.Exit(1)
			}
		})
	}
	return config
}

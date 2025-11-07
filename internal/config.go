package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

type Options struct {
	BackupMax   int8    `koanf:"backup_max"`
	Destination TPath_t `koanf:"destination"`
}

type Config struct {
	Games     map[string][]TPath_t `koanf:"games"`
	Opts      Options              `koanf:"options"`
	Variables map[string]string    `koanf:"variables,omitempty"`
}

type Path struct {
	Name string `koanf:"name" json:"name"`
	Type string `koanf:"type" json:"type"`
}

type BRConfig struct { // Backup Restore Config
	Game  string          `koanf:"game" json:"game"`
	Paths map[string]Path `koanf:"paths" json:"paths"`
}

var config *Config
var once sync.Once

var CfgFile string

func CheckErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func parseTOML(F string, o interface{}) {
	k := koanf.New(".")
	if err := k.Load(file.Provider(F), toml.Parser()); err != nil {
		CheckErr(fmt.Errorf("could not read config file: %v", err))
	}
	CheckErr(k.Unmarshal("", o))
}

func GetConfig() *Config {
	if config == nil {
		once.Do(func() {
			if CfgFile == "" {
				F, searchP := getExistingConfigPath()
				if F == "" {
					CheckErr(fmt.Errorf(
						"no config file found.\nsearched in: %v\ncreate one there, or specify using '-c'",
						strings.Join(searchP, ", ")))
				}
				CfgFile = F
			}
			config = &Config{}
			parseTOML(CfgFile, config)
		})
	}
	return config
}

func GetBRConfig(path string) *BRConfig {
	brCfg := &BRConfig{}
	parseTOML(path, brCfg)
	return brCfg
}

func WriteBRConfig(path string, bfCfg *BRConfig) {
	var k = koanf.New(".")
	k.Load(structs.Provider(bfCfg, "koanf"), nil)

	b, _ := k.Marshal(toml.Parser())

	out, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	CheckErr(err)
	defer out.Close()

	fmt.Fprint(out, string(b))
}

// Returns config file path, search paths
// if existing config file is not found in any
// search paths, empty string is returned
func getExistingConfigPath() (string, []string) {
	defCfgFName := "gamesnap.toml"

	cwd, err := filepath.Abs(".")
	CheckErr(err)

	paths := []string{cwd}
	home, _ := os.UserHomeDir()
	paths = append(paths, home)

	var foundPath string
	for _, P := range paths {
		P = filepath.Join(P, defCfgFName)
		if _, err := os.Stat(P); err == nil {
			foundPath = P
			break
		}
	}
	return foundPath, paths
}

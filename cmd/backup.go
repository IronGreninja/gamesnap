package cmd

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sync"
	"time"

	"github.com/IronGreninja/gamesnap/internal"
	"github.com/spf13/cobra"
)

var snapCmd = &cobra.Command{
	Use:   "snap [game]...",
	Short: "snapshot game saves",
	Run:   backupGames,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			games := internal.GetConfig().Games
			for _, name := range args {
				if _, ok := games[name]; !ok {
					CheckErr(fmt.Errorf("game %s not found", name))
				}
			}
		}
		return nil
	},
}

func saveBRConfig(brCfg *internal.BRConfig, path string) {
	path = filepath.Join(path, "info.toml")
	internal.WriteBRConfig(path, brCfg)
}

func backupSingleGame(name string, paths []internal.TPath_t, bkupPath string, wg *sync.WaitGroup) {
	defer wg.Done()

	vars := internal.GetConfig().Variables
	brCfg := internal.BRConfig{
		Game:  name,
		Paths: make(map[string]internal.Path),
	}
	for i, path := range paths {
		path, err := path.Resolve(vars)
		CheckErr(err)

		srcPath, err := internal.AbsPath(path)
		CheckErr(err)

		pname := fmt.Sprintf("path%d", i+1)
		snkDir := filepath.Join(bkupPath, pname)
		fInfo, err := os.Stat(srcPath)
		CheckErr(err)

		Path := internal.Path{
			Name: srcPath,
		}
		if fInfo.IsDir() {
			Path.Type = "dir"
		} else {
			Path.Type = "file"
		}
		CheckErr(internal.Copy(srcPath, snkDir))
		brCfg.Paths[pname] = Path
	}
	saveBRConfig(&brCfg, bkupPath)

	// create symlink 'latest' to this snapshot
	// needs elevated privileges on windows so skip for now :(
	// FIXME
	if runtime.GOOS != "windows" {
		target := bkupPath
		link := filepath.Join(filepath.Dir(target), "latest")
		CheckErr(internal.CreateSymlinkForce(target, link))
	}
}

func backupGames(cmd *cobra.Command, args []string) {
	cfg := internal.GetConfig()
	fmt.Println("using config:", internal.CfgFile)
	CheckErr(os.Chdir(filepath.Dir(internal.CfgFile)))
	vars := cfg.Variables
	bkupTime := time.Now().Format("2006-01-02_15.04.05") // windows does not allow ':' in pathnames

	bkupDestRoot, err := cfg.Opts.Destination.Resolve(vars)
	CheckErr(err)

	var games []string
	if len(args) == 0 {
		games = slices.Collect(maps.Keys(cfg.Games))
	} else {
		games = args
	}

	var wg sync.WaitGroup
	wg.Add(len(games))

	for _, game := range games {
		paths := cfg.Games[game]
		bkupPath := filepath.Join(bkupDestRoot, game, bkupTime)
		CheckErr(os.MkdirAll(bkupPath, 0700))
		go backupSingleGame(game, paths, bkupPath, &wg)
	}

	wg.Wait()
}

func init() {
	rootCmd.AddCommand(snapCmd)
}

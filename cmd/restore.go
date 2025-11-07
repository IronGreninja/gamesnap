package cmd

import (
	"os"
	"path/filepath"

	"github.com/IronGreninja/gamesnap/internal"
	"github.com/spf13/cobra"
)

var resnapCmd = &cobra.Command{
	Use:   "resnap snapshotDir",
	Short: "restore snapshot (Warning: Overwrites latest saves)",
	Args:  cobra.ExactArgs(1),
	Run:   restoreSnapshot,
}

func restoreSnapshot(cmd *cobra.Command, args []string) {
	bkupPath, err := internal.RealPath(args[0])
	CheckErr(err)

	brCfg := internal.GetBRConfig(filepath.Join(bkupPath, "info.toml"))

	for pathSrc, pType := range brCfg.Paths {
		if pType.Type == "file" {
			filename := filepath.Base(pType.Name)
			srcFile := filepath.Join(bkupPath, pathSrc, filename)
			destDir := filepath.Dir(pType.Name)

			internal.Copy(srcFile, destDir)
		} else {
			os.RemoveAll(pType.Name)
			err := internal.Copy(filepath.Join(bkupPath, pathSrc), pType.Name)
			CheckErr(err)
		}
	}
}

func init() {
	rootCmd.AddCommand(resnapCmd)
}

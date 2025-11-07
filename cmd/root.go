package cmd

import (
	"fmt"
	"os"

	"github.com/IronGreninja/gamesnap/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     internal.APPNAME,
	Short:   "A dead simple cli tool for game save snapshots",
	Version: internal.VERSION,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func CheckErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&internal.CfgFile, "config", "c", "", "config file")
}

func Execute() {
	CheckErr(rootCmd.Execute())
}

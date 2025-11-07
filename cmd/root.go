package cmd

import (
	"fmt"
	"os"

	"github.com/IronGreninja/gamesnap/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     internal.APPNAME,
	Short:   "A dead simple cli tool for game save snapshots",
	Version: internal.VERSION,
}

var cfgFile string

func CheckErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		cfgPaths := getConfigPaths()
		for _, path := range cfgPaths {
			viper.AddConfigPath(path)
		}
		viper.SetConfigName("." + internal.APPNAME)
		// viper.SetConfigType("yaml")
	}

	// viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config:", viper.ConfigFileUsed())
	} else {
		CheckErr(err)
	}
}

func getConfigPaths() []string {
	paths := []string{"."}
	home, _ := os.UserHomeDir()
	paths = append(paths, home)
	return paths
}

func Execute() {
	CheckErr(rootCmd.Execute())
}

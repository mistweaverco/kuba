package kuba

import (
	"fmt"
	"os"

	"github.com/mistweaverco/kuba/internal/config"
	"github.com/mistweaverco/kuba/internal/lib/log"
	"github.com/mistweaverco/kuba/internal/lib/version"
	"github.com/spf13/cobra"
)

var cfg = config.NewConfig(config.Config{
	Flags: config.ConfigFlags{},
})

var rootCmd = &cobra.Command{
	Use:   "kuba",
	Short: "Kuba CLI",
	Long:  "Kuba is a CLI tool for accessing secrets and environment variables in a secure and efficient way.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize logging with debug mode
		log.SetDebugMode(cfg.Flags.Debug)
	},
	Run: func(cmd *cobra.Command, files []string) {
		if cfg.Flags.Version {
			fmt.Println(version.VERSION)
			return
		} else {
			// TODO: Add a help message
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		osExit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&cfg.Flags.Debug, "debug", "d", false, "Enable debug mode for verbose logging")
}

// osExit is a variable to allow overriding in tests
var osExit = os.Exit

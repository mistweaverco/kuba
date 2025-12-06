package kuba

import (
	"github.com/mistweaverco/kuba/internal/lib/fileutils"
	"github.com/mistweaverco/kuba/internal/lib/log"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a default configuration file",
	Long:  "This command initializes a default configuration file for kuba, if it does not already exist.",
	Run: func(cmd *cobra.Command, files []string) {
		logger := log.NewLogger()
		logger.Debug("Initializing default kuba configuration")

		created := fileutils.GenerateDefaultKubaConfig()
		if created {
			logger.Debug("Default configuration file created successfully")
		} else {
			logger.Debug("Configuration file already exists, no action taken")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

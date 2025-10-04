package kuba

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage kuba configuration",
	Long: `Manage global kuba configuration settings.

This command allows you to configure various aspects of kuba's behavior,
including caching, logging, and other global settings.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// Add config command to root
	rootCmd.AddCommand(configCmd)

	// Add cache subcommand to config
	configCmd.AddCommand(configCacheCmd)
}

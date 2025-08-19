package kuba

import (
	"github.com/mistweaverco/kuba/internal/lib/fileutils"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a default configuration file",
	Long:  "This command initializes a default configuration file for kuba, if it does not already exist.",
	Run: func(cmd *cobra.Command, files []string) {
		fileutils.GenerateDefaultKubaConfig()
	},
}

package kuba

import (
	"fmt"

	"github.com/mistweaverco/kuba/internal/lib/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Kuba Version",
	Long:  "Displays the current version of Kuba CLI",
	Run: func(cmd *cobra.Command, files []string) {
		fmt.Println(version.VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().BoolVar(&cfg.Flags.Version, "version", false, "Kuba version")
}

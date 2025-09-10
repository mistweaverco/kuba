package kuba

import (
	"context"
	"fmt"

	"github.com/mistweaverco/kuba/internal/config"
	"github.com/mistweaverco/kuba/internal/lib/log"
	"github.com/mistweaverco/kuba/internal/lib/secrets"
	"github.com/spf13/cobra"
)

var (
	testEnvironment string
	testConfigFile  string
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test secret retrieval for an environment",
	Long: `Attempt to retrieve all mapped values for the selected environment.

This command will:
1. Locate and load the kuba.yaml configuration
2. Resolve the selected environment
3. Attempt to fetch all mapped values (secrets, paths, and literals)

It prints any errors encountered during retrieval.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTest()
	},
}

func init() {
	testCmd.Flags().StringVarP(&testEnvironment, "env", "e", "default", "Environment to use (default: default)")
	testCmd.Flags().StringVarP(&testConfigFile, "config", "c", "", "Path to kuba.yaml configuration file")
	rootCmd.AddCommand(testCmd)
}

func runTest() error {
	logger := log.NewLogger()

	// Find configuration file if not specified
	cfgPath := testConfigFile
	if cfgPath == "" {
		logger.Debug("No config file specified, searching for kuba.yaml")
		path, err := config.FindConfigFile()
		if err != nil {
			return fmt.Errorf("failed to find configuration file: %w", err)
		}
		cfgPath = path
		logger.Debug("Found configuration file", "path", cfgPath)
	} else {
		logger.Debug("Using specified configuration file", "path", cfgPath)
	}

	// Load configuration
	logger.Debug("Loading configuration from file")
	kubaConfig, err := config.LoadKubaConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	logger.Debug("Configuration loaded successfully")

	// Get environment configuration
	logger.Debug("Getting environment configuration", "environment", testEnvironment)
	env, err := kubaConfig.GetEnvironment(testEnvironment)
	if err != nil {
		return fmt.Errorf("failed to get environment '%s': %w", testEnvironment, err)
	}

	// Create secrets manager factory and attempt retrieval
	logger.Debug("Creating secrets manager factory")
	factory := secrets.NewSecretManagerFactory()
	ctx := context.Background()

	logger.Debug("Fetching secrets and values for environment")
	values, err := factory.GetSecretsForEnvironment(ctx, env)
	if err != nil {
		return fmt.Errorf("failed to retrieve values: %w", err)
	}

	// Success summary
	fmt.Printf("Successfully retrieved %d values for environment '%s'\n", len(values), testEnvironment)
	return nil
}

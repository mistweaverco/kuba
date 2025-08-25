package kuba

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/mistweaverco/kuba/internal/config"
	"github.com/mistweaverco/kuba/internal/lib/log"
	"github.com/mistweaverco/kuba/internal/lib/secrets"
	"github.com/spf13/cobra"
)

var (
	environment string
	configFile  string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command with secrets from cloud providers",
	Long: `Run a command with environment variables populated from secrets stored in cloud providers.
	
This command will:
1. Load the kuba.yaml configuration file
2. Fetch secrets from the specified cloud providers
3. Set the secrets as environment variables
4. Execute the provided command with those environment variables

Example:
  kuba run -- node server.js
  kuba run --env production -- python app.py
  kuba run --config ./config/kuba.yaml -- docker-compose up`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCommand(args)
	},
}

func init() {
	runCmd.Flags().StringVarP(&environment, "env", "e", "default", "Environment to use (default: default)")
	runCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to kuba.yaml configuration file")
	rootCmd.AddCommand(runCmd)
}

func runCommand(args []string) error {
	logger := log.NewLogger()

	// Find configuration file if not specified
	if configFile == "" {
		var err error
		logger.Debug("No config file specified, searching for kuba.yaml")
		configFile, err = config.FindConfigFile()
		if err != nil {
			return fmt.Errorf("failed to find configuration file: %w", err)
		}
		logger.Debug("Found configuration file", "path", configFile)
	} else {
		logger.Debug("Using specified configuration file", "path", configFile)
	}

	// Load configuration
	logger.Debug("Loading configuration from file")
	kubaConfig, err := config.LoadKubaConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	logger.Debug("Configuration loaded successfully")

	// Get environment configuration
	logger.Debug("Getting environment configuration", "environment", environment)
	env, err := kubaConfig.GetEnvironment(environment)
	if err != nil {
		return fmt.Errorf("failed to get environment '%s': %w", environment, err)
	}
	logger.Debug("Environment configuration retrieved", "environment", environment, "provider", env.Provider, "mappings_count", len(env.Mappings))

	// Create secrets manager factory
	logger.Debug("Creating secrets manager factory")
	factory := secrets.NewSecretManagerFactory()

	// Get secrets for the environment
	ctx := context.Background()
	logger.Debug("Fetching secrets from cloud providers")
	secrets, err := factory.GetSecretsForEnvironment(ctx, env)
	if err != nil {
		return fmt.Errorf("failed to get secrets: %w", err)
	}
	logger.Debug("Secrets retrieved successfully", "count", len(secrets))

	// Prepare command
	command := args[0]
	commandArgs := args[1:]
	logger.Debug("Preparing command execution", "command", command, "args", commandArgs)

	// Create command
	cmd := exec.Command(command, commandArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set environment variables
	cmd.Env = os.Environ()
	for key, value := range secrets {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	logger.Debug("Environment variables set", "secrets_count", len(secrets), "total_env_vars", len(cmd.Env))

	// Execute command
	logger.Debug("Executing command")
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			logger.Debug("Command exited with non-zero status", "exit_code", exitErr.ExitCode())
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("command failed: %w", err)
	}

	logger.Debug("Command executed successfully")
	return nil
}

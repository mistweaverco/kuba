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
	contain     bool
	commandFlag string
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

By default, secrets are merged with the current OS environment. Use --contain to only use
environment variables from kuba.yaml.

Example:
  kuba run -- node server.js
  kuba run --env production -- python app.py
  kuba run --config ./config/kuba.yaml -- docker-compose up
  kuba run --contain -- node server.js
  kuba run --command 'echo "$SOME_SECRET"'`,
	Args: func(cmd *cobra.Command, args []string) error {
		// If --command is provided, args are optional
		if cmd.Flags().Changed("command") {
			return nil
		}
		// Otherwise, require at least one argument
		if len(args) < 1 {
			return fmt.Errorf("requires at least 1 arg(s), only received %d", len(args))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCommand(args)
	},
}

func init() {
	runCmd.Flags().StringVarP(&environment, "env", "e", "default", "Environment to use (default: default)")
	runCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to kuba.yaml configuration file")
	runCmd.Flags().BoolVar(&contain, "contain", false, "Only use environment variables from kuba.yaml, do not merge with OS environment")
	runCmd.Flags().StringVar(&commandFlag, "command", "", "Run an arbitrary command string in a shell with access to injected environment variables")
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
	logger.Debug("Environment configuration retrieved", "environment", environment, "provider", env.Provider, "env_count", len(env.Env))

	// Create secrets manager factory
	logger.Debug("Creating secrets manager factory")
	factory := secrets.NewSecretManagerFactory()

	// Get secrets for the environment
	ctx := context.Background()
	logger.Debug("Fetching secrets from cloud providers")
	secrets, err := factory.GetSecretsForEnvironmentWithCache(ctx, env, configFile, environment)
	if err != nil {
		return fmt.Errorf("failed to get secrets: %w", err)
	}
	logger.Debug("Secrets retrieved successfully", "count", len(secrets))

	// Prepare command execution
	var cmd *exec.Cmd
	if commandFlag != "" {
		// Execute command string in a shell
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}
		logger.Debug("Preparing shell command execution", "shell", shell, "command", commandFlag)
		cmd = exec.Command(shell, "-c", commandFlag)
	} else {
		// Execute command directly (existing behavior)
		command := args[0]
		commandArgs := args[1:]
		logger.Debug("Preparing command execution", "command", command, "args", commandArgs)
		cmd = exec.Command(command, commandArgs...)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set environment variables
	if contain {
		// Only use secrets from kuba.yaml, do not merge with OS environment
		cmd.Env = make([]string, 0, len(secrets))
	} else {
		// Default behavior: merge OS environment with secrets
		cmd.Env = os.Environ()
	}

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

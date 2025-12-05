package kuba

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mistweaverco/kuba/internal/config"
	"github.com/mistweaverco/kuba/internal/lib/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	convertFrom    string
	convertEnv     string
	convertInfile  string
	convertOutfile string
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert configuration from other formats to kuba.yaml",
	Long: `Convert configuration files from other formats (e.g., dotenv) to kuba.yaml format.

This command helps migrate existing configurations to kuba.yaml format.
For dotenv files, it will create environment variable entries using the 'value' field.

Note: When updating an existing kuba.yaml file, comments within the modified
environment section will be lost as the section is regenerated. This is a limitation
of YAML manipulation - to preserve structure and data, comments in modified sections
cannot be retained. Consider backing up your kuba.yaml file before conversion if
comments are important.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConvert()
	},
}

func init() {
	convertCmd.Flags().StringVar(&convertFrom, "from", "", "Source format (e.g., 'dotenv')")
	convertCmd.Flags().StringVarP(&convertEnv, "env", "e", "default", "Environment name to use in kuba.yaml (default: default)")
	convertCmd.Flags().StringVar(&convertInfile, "infile", "", "Input file path (e.g., .env.example)")
	convertCmd.Flags().StringVar(&convertOutfile, "outfile", "", "Output kuba.yaml file path (default: kuba.yaml in current directory)")

	convertCmd.MarkFlagRequired("from")
	convertCmd.MarkFlagRequired("infile")

	rootCmd.AddCommand(convertCmd)
}

func runConvert() error {
	logger := log.NewLogger()

	if convertFrom != "dotenv" {
		return fmt.Errorf("unsupported source format: %s (only 'dotenv' is currently supported)", convertFrom)
	}

	// Determine output file path
	outPath := convertOutfile
	if outPath == "" {
		outPath = "kuba.yaml"
	}

	logger.Debug("Converting dotenv to kuba.yaml", "infile", convertInfile, "outfile", outPath, "env", convertEnv)

	// Read and parse dotenv file
	logger.Debug("Reading dotenv file", "path", convertInfile)
	envVars, err := parseDotenvFile(convertInfile)
	if err != nil {
		return fmt.Errorf("failed to parse dotenv file: %w", err)
	}
	logger.Debug("Parsed dotenv file", "variables_count", len(envVars))

	// Load existing kuba.yaml if it exists, or create new config
	var kubaConfig *config.KubaConfig
	var existingRawContent []byte
	var existingFileExists bool

	if _, err := os.Stat(outPath); err == nil {
		existingFileExists = true
		logger.Debug("Loading existing kuba.yaml", "path", outPath)

		// Read raw content to potentially preserve comments
		existingRawContent, err = os.ReadFile(outPath)
		if err != nil {
			return fmt.Errorf("failed to read existing kuba.yaml: %w", err)
		}

		// Also load as config struct for manipulation
		kubaConfig, err = config.LoadKubaConfig(outPath)
		if err != nil {
			return fmt.Errorf("failed to load existing kuba.yaml: %w", err)
		}
		logger.Debug("Loaded existing kuba.yaml", "environments_count", len(kubaConfig.Environments))
	} else {
		existingFileExists = false
		logger.Debug("No existing kuba.yaml found, creating new config")
		kubaConfig = &config.KubaConfig{
			Environments: make(map[string]config.Environment),
		}
	}

	// Create or update the environment
	env, exists := kubaConfig.Environments[convertEnv]
	if !exists {
		logger.Debug("Creating new environment", "env", convertEnv)
		// Create a new environment with local provider (since we're using values)
		env = config.Environment{
			Provider: "local",
			Project:  "",
			Env:      make(map[string]config.EnvItem),
		}
	} else {
		logger.Debug("Updating existing environment", "env", convertEnv)
		if env.Env == nil {
			env.Env = make(map[string]config.EnvItem)
		}
	}

	// Add dotenv entries to the environment
	// Since dotenv files contain actual values, we'll use the 'value' field
	// If the environment uses a different provider, we'll keep that but still add values
	// The user can later convert values to secrets if needed
	// Skip empty values - they should not be included in the config
	for key, value := range envVars {
		// Skip empty values
		if strings.TrimSpace(value) == "" {
			logger.Debug("Skipping empty environment variable", "key", key)
			continue
		}
		env.Env[key] = config.EnvItem{
			Value: value,
		}
		logger.Debug("Added environment variable", "key", key)
	}

	// Clean up empty values from the environment before writing
	cleanupEmptyValues(&env)

	// Update the environment in config
	kubaConfig.Environments[convertEnv] = env

	// Write the updated kuba.yaml
	logger.Debug("Writing kuba.yaml", "path", outPath)
	if err := writeKubaConfigWithCommentPreservation(outPath, kubaConfig, existingRawContent, existingFileExists); err != nil {
		return fmt.Errorf("failed to write kuba.yaml: %w", err)
	}

	fmt.Printf("Successfully converted %d variables from %s to kuba.yaml (environment: %s)\n", len(envVars), convertInfile, convertEnv)
	logger.Debug("Conversion completed successfully")
	return nil
}

// parseDotenvFile reads and parses a dotenv file
// It handles:
// - Comments (lines starting with #)
// - Blank lines
// - KEY=VALUE pairs
// - Quoted values (single and double quotes)
// - Multiline values (basic support)
func parseDotenvFile(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	envVars := make(map[string]string)
	scanner := bufio.NewScanner(file)

	var currentKey string
	var currentValue strings.Builder

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Check if this line continues a previous multiline value
		if currentKey != "" && (strings.HasPrefix(line, "\"") || strings.HasPrefix(line, "'")) {
			// This might be a continuation, but for simplicity, we'll treat each line independently
			currentKey = ""
			currentValue.Reset()
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			// Try to continue a multiline value
			if currentKey != "" {
				currentValue.WriteString("\n")
				currentValue.WriteString(line)
				continue
			}
			// Skip malformed lines
			continue
		}

		// If we had a previous key being built, save it now
		if currentKey != "" {
			envVars[currentKey] = strings.TrimSpace(currentValue.String())
			currentKey = ""
			currentValue.Reset()
		}

		key := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])

		// Skip empty keys
		if key == "" {
			continue
		}

		// Handle quoted values
		valueStr = unquoteValue(valueStr)

		// Check for multiline value (values ending with \)
		if strings.HasSuffix(valueStr, "\\") && !strings.HasSuffix(valueStr, "\\\\") {
			// Start building a multiline value
			currentKey = key
			currentValue.WriteString(strings.TrimSuffix(valueStr, "\\"))
			continue
		}

		envVars[key] = valueStr
	}

	// Handle any remaining multiline value
	if currentKey != "" {
		envVars[currentKey] = strings.TrimSpace(currentValue.String())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return envVars, nil
}

// unquoteValue removes surrounding quotes from a value if present
func unquoteValue(value string) string {
	value = strings.TrimSpace(value)

	// Handle double quotes
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		// Remove quotes and unescape
		unquoted := value[1 : len(value)-1]
		// Basic unescaping for common cases
		unquoted = strings.ReplaceAll(unquoted, "\\n", "\n")
		unquoted = strings.ReplaceAll(unquoted, "\\t", "\t")
		unquoted = strings.ReplaceAll(unquoted, "\\\"", "\"")
		unquoted = strings.ReplaceAll(unquoted, "\\\\", "\\")
		return unquoted
	}

	// Handle single quotes
	if len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'' {
		return value[1 : len(value)-1]
	}

	return value
}

// cleanupEmptyValues removes environment variables with empty values from the environment
// This ensures that empty values are not written to the YAML file
func cleanupEmptyValues(env *config.Environment) {
	cleanedEnv := make(map[string]config.EnvItem)
	for key, item := range env.Env {
		// Check if the item has any meaningful content
		hasContent := false

		// Check value
		if item.Value != nil {
			valueStr := fmt.Sprintf("%v", item.Value)
			if strings.TrimSpace(valueStr) != "" {
				hasContent = true
			}
		}

		// Check secret-key
		if item.SecretKey != "" {
			hasContent = true
		}

		// Check secret-path
		if item.SecretPath != "" {
			hasContent = true
		}

		// Only include items that have some content
		if hasContent {
			cleanedEnv[key] = item
		}
	}
	env.Env = cleanedEnv
}

// writeKubaConfigWithCommentPreservation writes a KubaConfig to a YAML file
// It attempts to preserve comments when updating existing files by using yaml.Node
func writeKubaConfigWithCommentPreservation(filePath string, cfg *config.KubaConfig, existingRawContent []byte, existingFileExists bool) error {
	schemaComment := "# yaml-language-server: $schema=https://kuba.mwco.app/kuba.schema.json\n---\n"
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// If we have existing content, try to preserve comments using yaml.Node
	if existingFileExists && len(existingRawContent) > 0 {
		// Parse existing YAML into a node tree (preserves comments)
		var existingNode yaml.Node
		if err := yaml.Unmarshal(existingRawContent, &existingNode); err == nil {
			// Try to update only the specific environment section
			// This is a best-effort attempt - some comments may still be lost
			if err := updateEnvironmentInNode(&existingNode, convertEnv, cfg.Environments[convertEnv]); err == nil {
				// Successfully updated the node tree, write it back
				var buf strings.Builder
				encoder := yaml.NewEncoder(&buf)
				encoder.SetIndent(2)
				if err := encoder.Encode(&existingNode); err == nil {
					encoder.Close()
					content := buf.String()

					// Ensure schema comment is present
					if !strings.Contains(content, "yaml-language-server") {
						content = schemaComment + content
					}

					return os.WriteFile(filePath, []byte(content), 0644)
				}
			}
			// If node-based update failed, fall through to struct-based marshaling
		}
	}

	// Fallback: marshal from struct (comments will be lost, but structure is correct)
	var buf strings.Builder
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}
	encoder.Close()

	content := buf.String()
	// Add schema comment at the top if file is new or doesn't have it
	if !strings.Contains(content, schemaComment) {
		content = schemaComment + content
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// updateEnvironmentInNode updates a specific environment section in a yaml.Node tree
// NOTE: This function replaces the entire environment node, which means comments within
// that environment section will be lost. Comments in other environments are preserved.
func updateEnvironmentInNode(rootNode *yaml.Node, envName string, env config.Environment) error {
	// The root node should be a document node
	if rootNode.Kind != yaml.DocumentNode && rootNode.Kind != yaml.MappingNode {
		return fmt.Errorf("unexpected root node kind: %v", rootNode.Kind)
	}

	// Find the mapping node (the actual content)
	var mappingNode *yaml.Node
	if rootNode.Kind == yaml.DocumentNode && len(rootNode.Content) > 0 {
		mappingNode = rootNode.Content[0]
	} else {
		mappingNode = rootNode
	}

	if mappingNode.Kind != yaml.MappingNode {
		return fmt.Errorf("expected mapping node, got %v", mappingNode.Kind)
	}

	// Find the environment key-value pair
	envNodeIndex := -1
	for i := 0; i < len(mappingNode.Content); i += 2 {
		if i+1 < len(mappingNode.Content) {
			keyNode := mappingNode.Content[i]
			if keyNode.Value == envName {
				envNodeIndex = i + 1
				break
			}
		}
	}

	// Create new environment node from the config
	// WARNING: This replaces the entire node, losing all comments within this environment section
	newEnvNode := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	// Add provider
	newEnvNode.Content = append(newEnvNode.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "provider"},
		&yaml.Node{Kind: yaml.ScalarNode, Value: env.Provider},
	)

	// Add project if present and not empty
	if strings.TrimSpace(env.Project) != "" {
		newEnvNode.Content = append(newEnvNode.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "project"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: env.Project},
		)
	}

	// Add env map
	envMapNode := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}
	for key, item := range env.Env {
		itemNode := &yaml.Node{
			Kind: yaml.MappingNode,
			Tag:  "!!map",
		}
		hasContent := false

		// Only add value if it's non-empty
		if item.Value != nil {
			valueStr := fmt.Sprintf("%v", item.Value)
			if strings.TrimSpace(valueStr) != "" {
				itemNode.Content = append(itemNode.Content,
					&yaml.Node{Kind: yaml.ScalarNode, Value: "value"},
					&yaml.Node{Kind: yaml.ScalarNode, Value: valueStr},
				)
				hasContent = true
			}
		}

		// Add secret-key if present
		if item.SecretKey != "" {
			itemNode.Content = append(itemNode.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: "secret-key"},
				&yaml.Node{Kind: yaml.ScalarNode, Value: item.SecretKey},
			)
			hasContent = true
		}

		// Add secret-path if present
		if item.SecretPath != "" {
			itemNode.Content = append(itemNode.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: "secret-path"},
				&yaml.Node{Kind: yaml.ScalarNode, Value: item.SecretPath},
			)
			hasContent = true
		}

		// Add provider if present and different from env-level provider
		if item.Provider != "" && item.Provider != env.Provider {
			itemNode.Content = append(itemNode.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: "provider"},
				&yaml.Node{Kind: yaml.ScalarNode, Value: item.Provider},
			)
			hasContent = true
		}

		// Add project if present and different from env-level project
		if item.Project != "" && item.Project != env.Project {
			itemNode.Content = append(itemNode.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: "project"},
				&yaml.Node{Kind: yaml.ScalarNode, Value: item.Project},
			)
			hasContent = true
		}

		// Only add the env item if it has some content
		if hasContent {
			envMapNode.Content = append(envMapNode.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: key},
				itemNode,
			)
		}
	}
	newEnvNode.Content = append(newEnvNode.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "env"},
		envMapNode,
	)

	// Update or add the environment
	if envNodeIndex >= 0 {
		// Update existing - this replaces the node, losing comments
		mappingNode.Content[envNodeIndex] = newEnvNode
	} else {
		// Add new environment
		mappingNode.Content = append(mappingNode.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: envName},
			newEnvNode,
		)
	}

	return nil
}

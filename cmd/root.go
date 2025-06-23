package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ironlabsdev/iron/cmd/generate"
	"github.com/ironlabsdev/iron/internal/config"
	"github.com/ironlabsdev/iron/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	jsonOutput  bool
	shortOutput bool
)

const (
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
	ColorBlue  = "\033[34m"
	ColorReset = "\033[0m"
)

var rootCmd = &cobra.Command{
	Use:   "iron",
	Short: "CLI tool for scaffolding educational programming projects",
	Long: `Iron CLI helps students scaffold programming projects and focus on core learning 
while simplifying the bootstrap process.

Iron provides templates and tools to generate common project structures,
allowing students to quickly get started with their learning without getting
bogged down in setup and configuration.`,
	Version: version.GetVersion(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Print error in red color
		_, _ = fmt.Fprintf(os.Stderr, "%s✗ %s%s\n", ColorRed, err.Error(), ColorReset)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Add subcommands
	rootCmd.AddCommand(generate.GenerateCmd)
	rootCmd.AddCommand(versionCmd)

	// Version command flags
	versionCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output version information in JSON format")
	versionCmd.Flags().BoolVarP(&shortOutput, "short", "s", false, "Output only the version number")

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.iron/config.yaml)")

	// Silence errors and usage on error (we handle them ourselves)
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	// Improve command suggestions
	rootCmd.DisableSuggestions = false
	rootCmd.SuggestionsMinimumDistance = 2

	// Disable default completion command (we can add our own later if needed)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Hide the help command (help is available via --help flag)
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	// Set version template
	rootCmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Create .iron directory if it doesn't exist
		ironDir := filepath.Join(home, config.Directory)
		if err := os.MkdirAll(ironDir, 0750); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error creating %s directory: %v\n", config.Directory, err)
			os.Exit(1)
		}

		// Search config in .iron directory
		viper.AddConfigPath(ironDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ironlabsdev/iron/cmd/generate"
	"github.com/ironlabsdev/iron/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var Verbose bool

const (
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
	ColorReset = "\033[0m"
)

var rootCmd = &cobra.Command{
	Use:     "iron",
	Short:   "",
	Long:    ``,
	Version: "0.1.0",
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

	rootCmd.AddCommand(generate.GenerateCmd)
	// Silence errors
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	rootCmd.DisableSuggestions = false
	rootCmd.SuggestionsMinimumDistance = 4
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})
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

		// Create .beam directory if it doesn't exist
		beamDir := filepath.Join(home, config.Directory)
		if err := os.MkdirAll(beamDir, 0750); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error creating %s directory: %v\n", config.Directory, err)
			os.Exit(1)
		}

		// Search config in .beam directory
		viper.AddConfigPath(beamDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// _, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

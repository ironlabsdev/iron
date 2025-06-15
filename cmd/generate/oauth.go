package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// OAuthCmd generates OAuth authentication boilerplate
var OAuthCmd = &cobra.Command{
	Use:   "oauth <project-name>",
	Short: "Generate OAuth authentication boilerplate",
	Long: `Generate a complete OAuth authentication implementation including:
- OAuth providers configuration
- Authentication handlers
- Token management
- User sessions
- Example routes and middleware

Use "." as project-name to use the current working directory.`,
	Args: cobra.ExactArgs(1), // Require exactly one argument
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}

		var fullPath string

		if projectName == "." {
			fullPath = cwd
		} else {
			fullPath = filepath.Join(cwd, projectName)
		}

		fullPath = filepath.Clean(fullPath)

		return FromTemplate("oauth", fullPath)
	},
}

func init() {
	GenerateCmd.AddCommand(OAuthCmd)
}

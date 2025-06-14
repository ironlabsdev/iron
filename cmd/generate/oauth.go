package generate

import (
	"github.com/spf13/cobra"
)

// OAuthCmd generates OAuth authentication boilerplate
var OAuthCmd = &cobra.Command{
	Use:   "oauth [project-name]",
	Short: "Generate OAuth authentication boilerplate",
	Long: `Generate a complete OAuth authentication implementation including:
- OAuth providers configuration
- Authentication handlers
- Token management
- User sessions
- Example routes and middleware`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := "oauth-project"
		if len(args) > 0 {
			projectName = args[0]
		}

		return FromTemplate("oauth", projectName)
	},
}

func init() {
	GenerateCmd.AddCommand(OAuthCmd)
}

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ironlabsdev/iron/internal/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print detailed version information about Iron CLI including build date, git commit, and platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		if jsonOutput {
			info := version.GetBuildInfo()
			output, err := json.MarshalIndent(info, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%sError marshaling version info: %s%s\n", ColorRed, err.Error(), ColorReset)
				os.Exit(1)
			}
			fmt.Println(string(output))
		} else if shortOutput {
			fmt.Println(version.GetVersion())
		} else {
			fmt.Println(version.GetFullVersion())
		}
	},
}

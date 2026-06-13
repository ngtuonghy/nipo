package version

import (
	"fmt"
	"nipo/pkg/ui"
	"github.com/spf13/cobra"
)

// Command creates and returns the version cobra subcommand to print version info.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Nipo",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Nipo v%s\n", ui.AppVersion)
		},
	}

	return cmd
}

package status

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Command creates and returns the status cobra subcommand to check active tunnels.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check tunnel status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Tunnel status: Active")
			// TODO: Implement actual status check
		},
	}

	return cmd
}

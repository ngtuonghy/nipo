package http

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"nipo/pkg/config"
	"nipo/pkg/ui"
)

// Command creates and returns the http cobra subcommand to start a single tunnel.
func Command() *cobra.Command {
	var subdomain string

	cmd := &cobra.Command{
		Use:   "http [port]",
		Short: "Start an HTTP tunnel",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config.EnsureLanguageSelected()
			
			port, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Invalid port number")
				return
			}

			// Load subdomain from flag, fallback to config default
			finalSubdomain := config.AppConfig.DefaultSubdomain
			if subdomain != "" {
				finalSubdomain = subdomain
			}

			model := ui.InitialModel(port, finalSubdomain, config.DefaultBackendURL)
			p := tea.NewProgram(model)

			// Handle termination signals gracefully
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
			go func() {
				<-sigs
				p.Quit()
			}()

			finalModel, err := p.Run()
			if err != nil {
				fmt.Printf("Error starting UI: %v\n", err)
				os.Exit(1)
			}
			// Trigger cleanup routines on exit
			if m, ok := finalModel.(ui.TunnelModel); ok {
				m.ShutdownWithFeedback()
			}
		},
	}

	cmd.Flags().StringVarP(&subdomain, "subdomain", "s", "", "Custom subdomain")
	return cmd
}

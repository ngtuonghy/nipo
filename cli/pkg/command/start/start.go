package start

import (
	"fmt"
	"nipo/pkg/config"
	"nipo/pkg/ui"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// Command creates and returns the start cobra subcommand to run configured tunnels.
func Command() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "start [tunnel-name...]",
		Short: "Start tunnels defined in config",
		Run: func(cmd *cobra.Command, args []string) {
			config.EnsureLanguageSelected()

			// If configFile is set, load config from that path (can be URL or local file)
			if configFile != "" {
				if err := config.LoadConfigFromFile(configFile); err != nil {
					fmt.Printf("Error loading config file: %v\n", err)
					os.Exit(1)
				}
			}

			tunnels := config.AppConfig.Tunnels
			if len(tunnels) == 0 {
				fmt.Println("No tunnels defined in configuration file.")
				os.Exit(1)
			}

			// If specific tunnels are requested, filter them
			if len(args) > 0 {
				var filtered []config.TunnelConfig
				for _, arg := range args {
					found := false
					for _, t := range tunnels {
						if t.Name == arg {
							filtered = append(filtered, t)
							found = true
							break
						}
					}
					if !found {
						fmt.Printf("Error: tunnel '%s' not found in configuration file.\n", arg)
						os.Exit(1)
					}
				}
				tunnels = filtered
			}

			model := ui.InitialMultiModel(tunnels, config.DefaultBackendURL)
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

	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Path or URL to the configuration file (e.g. -c nipo.yml or -c https://example.com/nipo.yml)")
	return cmd
}

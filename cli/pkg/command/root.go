package command

import (
	"strings"

	config_cmd "nipo/pkg/command/config"
	"nipo/pkg/command/http"
	"nipo/pkg/command/start"
	"nipo/pkg/command/status"
	"nipo/pkg/command/version"
	"nipo/pkg/config"
	"nipo/pkg/ui"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nipo",
	Short: "Nipo is a production-ready tunneling platform",
	Long:  `Nipo allows you to create secure HTTP tunnels from localhost to the internet.`,
	Example: `  # Start an HTTP tunnel on port 3000
  nipo http 3000

  # Start a tunnel with custom subdomain
  nipo http 3000 --subdomain myapp

  # Start all tunnels defined in local configuration file (nipo.yml)
  nipo start

  # Start specific tunnels from the config file
  nipo start web api

  # View and configure interface language
  nipo config
  nipo config --language vi`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// It also handles configuration initialization and TUI-level localization mapping.
func Execute() error {
	config.InitConfig()

	// Explicitly initialize default commands and flags so we can translate them
	rootCmd.InitDefaultHelpCmd()
	rootCmd.InitDefaultCompletionCmd()
	rootCmd.InitDefaultHelpFlag()

	tr := ui.GetT()
	rootCmd.Short = tr.RootShort
	rootCmd.Long = tr.RootLong
	rootCmd.Example = tr.RootExample

	if config.AppConfig.Lang == "vi" {
		template := rootCmd.UsageTemplate()
		template = strings.ReplaceAll(template, "Usage:", "Cách dùng:")
		template = strings.ReplaceAll(template, "Examples:", "Ví dụ:")
		template = strings.ReplaceAll(template, "Available Commands:", "Các lệnh khả dụng:")
		template = strings.ReplaceAll(template, "Flags:", "Tùy chọn:")
		template = strings.ReplaceAll(template, "Global Flags:", "Tùy chọn chung:")
		template = strings.ReplaceAll(template, "Additional help topics:", "Chủ đề trợ giúp bổ sung:")
		template = strings.ReplaceAll(template, "Use \"{{.CommandPath}} [command] --help\" for more information about a command.", "Sử dụng \"{{.CommandPath}} [command] --help\" để xem thông tin chi tiết về lệnh đó.")
		rootCmd.SetUsageTemplate(template)
	}

	if helpFlag := rootCmd.Flags().Lookup("help"); helpFlag != nil {
		helpFlag.Usage = tr.FlagHelpUsage
	}

	for _, cmd := range rootCmd.Commands() {
		switch cmd.Name() {
		case "config":
			cmd.Short = tr.CmdConfigShort
		case "http":
			cmd.Short = tr.CmdHttpShort
		case "start":
			cmd.Short = tr.CmdStartShort
		case "status":
			cmd.Short = tr.CmdStatusShort
		case "version":
			cmd.Short = tr.CmdVersionShort
		case "help":
			cmd.Short = tr.CmdHelpShort
		case "completion":
			cmd.Short = tr.CmdCompletionShort
		}
	}

	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(config.InitConfig)

	rootCmd.AddCommand(http.Command())
	rootCmd.AddCommand(status.Command())
	rootCmd.AddCommand(start.Command())
	rootCmd.AddCommand(config_cmd.Command())
	rootCmd.AddCommand(version.Command())
}

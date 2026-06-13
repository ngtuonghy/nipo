package config_cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"nipo/pkg/config"
	"nipo/pkg/ui"
)

// Command creates and returns the config cobra subcommand to manage language options.
func Command() *cobra.Command {
	var language string

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure Nipo CLI settings",
		Run: func(cmd *cobra.Command, args []string) {
			tr := ui.GetT()

			if language != "" {
				if language != "en" && language != "vi" {
					fmt.Println(tr.ConfigErrUnsupported)
					os.Exit(1)
				}
				
				config.AppConfig.Lang = language

				home, err := os.UserHomeDir()
				if err != nil {
					fmt.Printf("Error getting home directory: %v\n", err)
					os.Exit(1)
				}

				dir := filepath.Join(home, ".nipo")
				os.MkdirAll(dir, 0755)
				globalPath := filepath.Join(dir, "config.yml")

				vGlobal := viper.New()
				vGlobal.SetConfigFile(globalPath)
				vGlobal.Set("lang", language)
				if err := vGlobal.WriteConfigAs(globalPath); err != nil {
					fmt.Printf("Error writing config file: %v\n", err)
					os.Exit(1)
				}
				
				// Re-fetch translation block since language just updated
				tr = ui.GetT()
				fmt.Printf(tr.ConfigLangUpdated+"\n", language)
				return
			}

			// If no flags are provided, show current configuration
			fmt.Println(tr.ConfigHeader)
			fmt.Printf(tr.ConfigLangLabel+"\n", config.AppConfig.Lang)
			fmt.Println(tr.ConfigHint)
		},
	}

	cmd.Flags().StringVarP(&language, "language", "l", "", "Set interface language ('en' or 'vi')")
	return cmd
}

package config

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	"github.com/charmbracelet/lipgloss"
)

var Version = "v0.2.0" // x-release-please-version
// TunnelConfig holds specifications for a single tunnel.
type TunnelConfig struct {
	Name      string `mapstructure:"name"`
	Protocol  string `mapstructure:"protocol"`
	Port      int    `mapstructure:"port"`
	Subdomain string `mapstructure:"subdomain"`
}

// DefaultBackendURL is the hardcoded API gateway backend for Nipo tunnels.
const DefaultBackendURL = "https://api.ngtuonghy.online"

// Config represents the application configuration schema.
type Config struct {
	DefaultSubdomain string         `mapstructure:"default_subdomain"`
	Lang             string         `mapstructure:"lang"`
	Tunnels          []TunnelConfig `mapstructure:"-"`
}

// AppConfig is the globally accessible active application configuration.
var AppConfig Config

// parseTunnels converts the raw tunnels input (map or array) into a structured TunnelConfig slice.
func parseTunnels(rawTunnels interface{}) ([]TunnelConfig, error) {
	var list []TunnelConfig

	switch m := rawTunnels.(type) {
	case map[string]interface{}:
		for name, val := range m {
			tc := TunnelConfig{Name: name}
			switch v := val.(type) {
			case int:
				tc.Port = v
			case int32:
				tc.Port = int(v)
			case int64:
				tc.Port = int(v)
			case float64:
				tc.Port = int(v)
			case string:
				port, err := strconv.Atoi(v)
				if err == nil {
					tc.Port = port
				}
			case map[string]interface{}:
				if p, ok := v["port"]; ok {
					switch pv := p.(type) {
					case int:
						tc.Port = pv
					case int32:
						tc.Port = int(pv)
					case int64:
						tc.Port = int(pv)
					case float64:
						tc.Port = int(pv)
					}
				}
				if sub, ok := v["subdomain"].(string); ok {
					tc.Subdomain = sub
				}
				if proto, ok := v["protocol"].(string); ok {
					tc.Protocol = proto
				}
			}
			if tc.Port > 0 {
				list = append(list, tc)
			}
		}
	case []interface{}:
		for _, item := range m {
			if imap, ok := item.(map[string]interface{}); ok {
				tc := TunnelConfig{}
				if name, ok := imap["name"].(string); ok {
					tc.Name = name
				}
				if portVal, ok := imap["port"]; ok {
					switch pv := portVal.(type) {
					case int:
						tc.Port = pv
					case int32:
						tc.Port = int(pv)
					case int64:
						tc.Port = int(pv)
					case float64:
						tc.Port = int(pv)
					}
				}
				if sub, ok := imap["subdomain"].(string); ok {
					tc.Subdomain = sub
				}
				if proto, ok := imap["protocol"].(string); ok {
					tc.Protocol = proto
				}
				if tc.Name != "" && tc.Port > 0 {
					list = append(list, tc)
				}
			}
		}
	}

	// Sort by name for deterministic loading order
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list, nil
}

// InitConfig initializes the config system, loading global config and local configs sequentially.
func InitConfig() {
	// Read global settings from user home directory for language preference
	home, err := os.UserHomeDir()
	if err == nil {
		globalPath := filepath.Join(home, ".nipo", "config.yml")
		if _, err := os.Stat(globalPath); err == nil {
			vGlobal := viper.New()
			vGlobal.SetConfigFile(globalPath)
			if err := vGlobal.ReadInConfig(); err == nil {
				AppConfig.Lang = vGlobal.GetString("lang")
			}
		}
	}

	if AppConfig.Lang == "" {
		AppConfig.Lang = "en"
	}

	// Read local config file (nipo.yml) if present
	viper.SetConfigName("nipo")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err == nil {
		var localConfig struct {
			DefaultSubdomain string      `mapstructure:"default_subdomain"`
			Lang             string      `mapstructure:"lang"`
			Tunnels          interface{} `mapstructure:"tunnels"`
		}
		if err := viper.Unmarshal(&localConfig); err == nil {
			if localConfig.DefaultSubdomain != "" {
				AppConfig.DefaultSubdomain = localConfig.DefaultSubdomain
			}
			if localConfig.Lang != "" {
				AppConfig.Lang = localConfig.Lang
			}
			if localConfig.Tunnels != nil {
				tunnels, _ := parseTunnels(localConfig.Tunnels)
				AppConfig.Tunnels = tunnels
			}
		}
	}

	// Apply default values if empty
	if AppConfig.DefaultSubdomain == "" {
		AppConfig.DefaultSubdomain = fmt.Sprintf("nipo-%d", rand.Intn(90000)+10000)
	}

	// Ensure all tunnels have assigned subdomains
	for i := range AppConfig.Tunnels {
		if AppConfig.Tunnels[i].Subdomain == "" {
			AppConfig.Tunnels[i].Subdomain = fmt.Sprintf("%s-%s", AppConfig.DefaultSubdomain, AppConfig.Tunnels[i].Name)
		}
	}
}

type langSelectModel struct {
	cursor   int
	quitting bool
}

func (m langSelectModel) Init() tea.Cmd {
	return nil
}

func (m langSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			os.Exit(0)
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 1 {
				m.cursor++
			}
		case "enter", " ":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m langSelectModel) View() string {
	if m.quitting {
		return ""
	}
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Bold(true)
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	s := "\n" + titleStyle.Render("Please choose your language / Vui lòng chọn ngôn ngữ:") + "\n\n"

	options := []string{"English (en)", "Tiếng Việt (vi)"}
	for i, opt := range options {
		if m.cursor == i {
			s += cursorStyle.Render("> ") + selectedStyle.Render(opt) + "\n"
		} else {
			s += "  " + opt + "\n"
		}
	}

	s += "\n" + dimStyle.Render("[↑↓ navigate  Enter select]") + "\n"
	return s
}

// EnsureLanguageSelected prompts the user to select their interface language (TUI selection)
// if it is not already configured in the global config.yml file.
func EnsureLanguageSelected() {
	// Check if global config contains a valid language option
	home, err := os.UserHomeDir()
	if err == nil {
		globalPath := filepath.Join(home, ".nipo", "config.yml")
		if _, err := os.Stat(globalPath); err == nil {
			vGlobal := viper.New()
			vGlobal.SetConfigFile(globalPath)
			if err := vGlobal.ReadInConfig(); err == nil {
				lang := vGlobal.GetString("lang")
				if lang == "en" || lang == "vi" {
					AppConfig.Lang = lang
					return
				}
			}
		}
	}

	p := tea.NewProgram(langSelectModel{})
	m, errRun := p.Run()
	if errRun != nil {
		AppConfig.Lang = "en"
		return
	}

	selectedLang := "en"
	if model, ok := m.(langSelectModel); ok && model.cursor == 1 {
		selectedLang = "vi"
	}

	AppConfig.Lang = selectedLang

	// Persist the selected language globally to config.yml
	if err == nil && home != "" {
		dir := filepath.Join(home, ".nipo")
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Warning: failed to create config directory: %v\n", err)
			return
		}
		globalPath := filepath.Join(dir, "config.yml")

		vGlobal := viper.New()
		vGlobal.SetConfigFile(globalPath)
		vGlobal.Set("lang", selectedLang)
		if err := vGlobal.WriteConfigAs(globalPath); err != nil {
			fmt.Printf("Warning: failed to save language setting: %v\n", err)
		}

		step0, step1, step2, step3 := "Nipo Tunnel", "Hello", "Xin chào", "Nipo Tunnel - From localhost to the world!"
		if selectedLang == "vi" {
			step0, step1, step2, step3 = "Nipo Tunnel", "Xin chào", "Hello", "Nipo Tunnel - Đưa code nhà làm ra biển lớn!"
		}

		badgeStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1).
			Bold(true)

		plainStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Bold(true)

		// Print initial spacing
		fmt.Print("\n\n")

		// printStep clears the current line using a carriage return and spaces, then prints the new text.
		printStep := func(text string) {
			fmt.Print("\r                                                                ")
			fmt.Printf("\r  %s", text)
		}

		// Step 0
		printStep(badgeStyle.Render(step0))
		time.Sleep(1500 * time.Millisecond)

		// Step 1
		printStep(badgeStyle.Render(step1))
		time.Sleep(1500 * time.Millisecond)

		// Step 2
		printStep(badgeStyle.Render(step2))
		time.Sleep(1500 * time.Millisecond)

		// Step 3
		printStep(plainStyle.Render(step3))
		time.Sleep(2500 * time.Millisecond)
		fmt.Println() // final newline before starting app
	}
}

// LoadConfigFromFile loads the configuration from a local file or remote URL path.
func LoadConfigFromFile(path string) error {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(path)
		if err != nil {
			return fmt.Errorf("failed to fetch config from URL: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to fetch config from URL: status %d", resp.StatusCode)
		}
		viper.SetConfigType("yml")
		if err := viper.ReadConfig(resp.Body); err != nil {
			return fmt.Errorf("failed to parse remote config: %w", err)
		}
	} else {
		viper.SetConfigFile(path)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("read config file: %w", err)
		}
	}

	var rawConfig struct {
		DefaultSubdomain string      `mapstructure:"default_subdomain"`
		Lang             string      `mapstructure:"lang"`
		Tunnels          interface{} `mapstructure:"tunnels"`
	}
	if err := viper.Unmarshal(&rawConfig); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	if rawConfig.DefaultSubdomain != "" {
		AppConfig.DefaultSubdomain = rawConfig.DefaultSubdomain
	}
	if rawConfig.Lang != "" {
		AppConfig.Lang = rawConfig.Lang
	}
	if rawConfig.Tunnels != nil {
		tunnels, err := parseTunnels(rawConfig.Tunnels)
		if err != nil {
			return err
		}
		AppConfig.Tunnels = tunnels
	}

	if AppConfig.DefaultSubdomain == "" {
		AppConfig.DefaultSubdomain = fmt.Sprintf("nipo-%d", rand.Intn(90000)+10000)
	}

	if AppConfig.Lang == "" {
		AppConfig.Lang = "en"
	}

	for i := range AppConfig.Tunnels {
		if AppConfig.Tunnels[i].Subdomain == "" {
			AppConfig.Tunnels[i].Subdomain = fmt.Sprintf("%s-%s", AppConfig.DefaultSubdomain, AppConfig.Tunnels[i].Name)
		}
	}
	return nil
}

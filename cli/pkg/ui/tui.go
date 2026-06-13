package ui

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"sort"
	"strings"
	"time"

	"nipo-tunnel/internal/tunnel"
	"nipo-tunnel/pkg/config"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

var AppVersion = "v0.2.0" // x-release-please-version

const (
	stateStartingProxy state = iota
	stateDownloading
	stateStartingTunnel
	stateRegistering
	stateVerifyingDNS
	stateDone
	stateError
	stateConflict // stateConflict indicates the subdomain is already in use, waiting for user choice
)

// TunnelState tracks the runtime state of a single tunnel instance.
type TunnelState struct {
	Config         config.TunnelConfig
	State          state
	Err            error
	Proxy          *tunnel.Proxy
	NodeURL        string
	PublicURL      string
	CloudflaredCmd *exec.Cmd // CloudflaredCmd tracks the subprocess execution
}

// TunnelModel acts as the main state machine for Bubble Tea TUI rendering.
type TunnelModel struct {
	BackendURL      string
	Spinner         spinner.Model
	Progress        progress.Model
	DownloadPercent float64
	DownloadChan    chan float64
	Tunnels         []TunnelState
	BinPath         string
	DownloadErr     error
	Width           int
	SessionEndsAt   time.Time
	ConflictCursor  int // ConflictCursor tracks menu selection index: 0 = random, 1 = exit
	Ctx             context.Context
	Cancel          context.CancelFunc
}

type tickMsg time.Time

// InitialMultiModel sets up the state model for multi-tunnel mode.
func InitialMultiModel(configs []config.TunnelConfig, backend string) TunnelModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	pg := progress.New(progress.WithDefaultGradient())
	pg.Width = 40

	states := make([]TunnelState, len(configs))
	for i, c := range configs {
		if c.Name == "" {
			c.Name = fmt.Sprintf("tunnel-%d", i+1)
		}
		states[i] = TunnelState{
			Config: c,
			State:  stateStartingProxy,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	return TunnelModel{
		BackendURL:   backend,
		Spinner:      s,
		Progress:     pg,
		DownloadChan: make(chan float64, 50),
		Tunnels:      states,
		Ctx:          ctx,
		Cancel:       cancel,
	}
}

// InitialModel configures and initializes the TUI model for single tunnel mode.
func InitialModel(port int, subdomain, backend string) TunnelModel {
	return InitialMultiModel([]config.TunnelConfig{
		{Name: "web", Port: port, Subdomain: subdomain},
	}, backend)
}

// ShutdownWithFeedback prints a releasing message, runs Cleanup, then erases
// the message lines so the terminal is left clean.
func (m TunnelModel) ShutdownWithFeedback() {
	var subs []string
	for _, t := range m.Tunnels {
		if t.Config.Subdomain != "" && m.BackendURL != "" {
			subs = append(subs, t.Config.Subdomain)
		}
	}

	linesWritten := 0
	if len(subs) > 0 {
		fmt.Print("\n")
		linesWritten++
		for _, s := range subs {
			fmt.Printf("  \033[33m⣿\033[0m  Đang giải phóng subdomain '%s'...\n", s)
			linesWritten++
		}
		fmt.Print("\n")
		linesWritten++
	}

	m.Cleanup()

	// Erase the lines we printed above
	if linesWritten > 0 {
		fmt.Printf("\033[%dA\033[J", linesWritten)
	}
}

// Cleanup unregisters all active subdomains from KV and kills our cloudflared processes.
func (m TunnelModel) Cleanup() {
	// Cancel the tunnel context first to stop all running goroutines
	if m.Cancel != nil {
		m.Cancel()
	}
	// Kill cloudflared processes
	for _, t := range m.Tunnels {
		if t.CloudflaredCmd != nil && t.CloudflaredCmd.Process != nil {
			t.CloudflaredCmd.Process.Kill()
		}
	}
	// Use a fresh context for unregister — m.Ctx is already cancelled at this point
	unregCtx, unregCancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer unregCancel()
	for _, t := range m.Tunnels {
		if t.Config.Subdomain != "" && m.BackendURL != "" {
			UnregisterSubdomain(unregCtx, m.BackendURL, t.Config.Subdomain)
		}
	}
}

func uiTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m TunnelModel) Init() tea.Cmd {
	cmds := []tea.Cmd{m.Spinner.Tick, uiTick(), downloadTask(m.DownloadChan), listenToProgress(m.DownloadChan)}
	for i, t := range m.Tunnels {
		cmds = append(cmds, startProxyTask(m.Ctx, i, t.Config.Name, t.Config.Port))
	}
	return tea.Batch(cmds...)
}

func (m TunnelModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle interactive menu navigation when a tunnel is in conflict state
		for i, t := range m.Tunnels {
			if t.State == stateConflict {
				switch msg.String() {
				case "up", "k":
					if m.ConflictCursor > 0 {
						m.ConflictCursor--
					}
				case "down", "j":
					if m.ConflictCursor < 1 {
						m.ConflictCursor++
					}
				case "enter", " ":
					if m.ConflictCursor == 0 {
						// Use config subdomain as base, strip trailing numeric suffix if present
						base := t.Config.Subdomain
						parts := strings.Split(base, "-")
						if len(parts) > 1 {
							lastPart := parts[len(parts)-1]
							allDigits := len(lastPart) > 0
							for _, c := range lastPart {
								if c < '0' || c > '9' {
									allDigits = false
									break
								}
							}
							if allDigits {
								base = strings.Join(parts[:len(parts)-1], "-")
							}
						}
						newSub := fmt.Sprintf("%s-%05d", base, rand.Intn(100000))
						m.Tunnels[i].Config.Subdomain = newSub
						m.Tunnels[i].State = stateRegistering
						return m, registerTask(m.Ctx, i, m.BackendURL, newSub, m.Tunnels[i].NodeURL)
					}
					return m, tea.Quit
				case "ctrl+c", "q":
					return m, tea.Quit
				}
				return m, nil
			}
		}
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		return m, nil

	case tickMsg:
		if !m.SessionEndsAt.IsZero() && time.Now().After(m.SessionEndsAt) {
			return m, tea.Quit
		}
		return m, uiTick()

	case proxyStartedMsg:
		if msg.err != nil {
			m.Tunnels[msg.id].State = stateError
			m.Tunnels[msg.id].Err = msg.err
			return m, nil
		}
		m.Tunnels[msg.id].Proxy = msg.proxy

		// If download already finished, start tunnel right away
		if m.BinPath != "" {
			m.Tunnels[msg.id].State = stateStartingTunnel
			delay := time.Duration(msg.id) * 2 * time.Second
			return m, startTunnelTask(m.Ctx, msg.id, m.BinPath, m.Tunnels[msg.id].Proxy.ListenPort, delay)
		} else if m.DownloadErr != nil {
			m.Tunnels[msg.id].State = stateError
			m.Tunnels[msg.id].Err = m.DownloadErr
		} else {
			m.Tunnels[msg.id].State = stateDownloading
		}
		return m, nil

	case downloadProgressMsg:
		m.DownloadPercent = float64(msg)
		cmd := m.Progress.SetPercent(float64(msg))
		return m, tea.Batch(cmd, listenToProgress(m.DownloadChan))

	case progress.FrameMsg:
		newModel, cmd := m.Progress.Update(msg)
		if newProgressModel, ok := newModel.(progress.Model); ok {
			m.Progress = newProgressModel
		}
		return m, cmd

	case downloadCompleteMsg:
		if msg.err != nil {
			m.DownloadErr = msg.err
			for i := range m.Tunnels {
				m.Tunnels[i].State = stateError
				m.Tunnels[i].Err = msg.err
			}
			return m, tea.Quit
		}
		m.BinPath = msg.binPath
		var cmds []tea.Cmd
		for i, t := range m.Tunnels {
			if t.State == stateDownloading && t.Proxy != nil {
				m.Tunnels[i].State = stateStartingTunnel
				delay := time.Duration(i) * 2 * time.Second
				cmds = append(cmds, startTunnelTask(m.Ctx, i, m.BinPath, t.Proxy.ListenPort, delay))
			}
		}
		return m, tea.Batch(cmds...)

	case tunnelStartedMsg:
		if msg.err != nil {
			m.Tunnels[msg.id].State = stateError
			m.Tunnels[msg.id].Err = msg.err
			return m, nil
		}
		m.Tunnels[msg.id].NodeURL = msg.url
		m.Tunnels[msg.id].CloudflaredCmd = msg.cmd // Save exact process reference

		if m.BackendURL != "" && m.Tunnels[msg.id].Config.Subdomain != "" {
			m.Tunnels[msg.id].State = stateRegistering
			return m, registerTask(m.Ctx, msg.id, m.BackendURL, m.Tunnels[msg.id].Config.Subdomain, msg.url)
		}
		m.Tunnels[msg.id].State = stateDone
		return m, nil

	case registerCompleteMsg:
		if msg.err != nil {
			if strings.Contains(msg.err.Error(), "already in use") {
				// Handle subdomain conflict by prompting the user
				m.Tunnels[msg.id].State = stateConflict
				m.ConflictCursor = 0
			} else {
				m.Tunnels[msg.id].State = stateError
				m.Tunnels[msg.id].Err = msg.err
			}
		} else {
			m.Tunnels[msg.id].PublicURL = fmt.Sprintf("https://%s.ngtuonghy.online", m.Tunnels[msg.id].Config.Subdomain)
			m.Tunnels[msg.id].State = stateVerifyingDNS
			return m, verifyDNSTask(m.Ctx, msg.id, m.Tunnels[msg.id].PublicURL)
		}
		return m, nil

	case verifyDNSCompleteMsg:
		if msg.err != nil {
			m.Tunnels[msg.id].State = stateError
			m.Tunnels[msg.id].Err = msg.err
			return m, tea.Quit
		} else {
			m.Tunnels[msg.id].State = stateDone

			// Start the session timer only when the tunnel successfully connects
			if m.SessionEndsAt.IsZero() {
				m.SessionEndsAt = time.Now().Add(9 * time.Hour)
			}
		}

		// Check if ALL tunnels are now done, then start a single heartbeat for all
		allOnline := true
		var activeSubdomains []string
		for _, t := range m.Tunnels {
			if t.State != stateDone {
				allOnline = false
			}
			if t.Config.Subdomain != "" {
				activeSubdomains = append(activeSubdomains, t.Config.Subdomain)
			}
		}
		if allOnline && m.BackendURL != "" && len(activeSubdomains) > 0 {
			return m, heartbeatTask(m.Ctx, m.BackendURL, activeSubdomains)
		}
		return m, nil

	case heartbeatMsg:
		// Re-schedule the next heartbeat with the current active subdomains
		var activeSubdomains []string
		for _, t := range m.Tunnels {
			if t.State == stateDone && t.Config.Subdomain != "" {
				activeSubdomains = append(activeSubdomains, t.Config.Subdomain)
			}
		}
		if m.BackendURL != "" && len(activeSubdomains) > 0 {
			return m, heartbeatTask(m.Ctx, m.BackendURL, activeSubdomains)
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m TunnelModel) View() string {
	tr := GetT()

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Margin(1, 0, 1, 2)

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Width(15).MarginLeft(2)
	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	urlStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("87")).Bold(true).Underline(true)

	doc := titleStyle.Render(tr.Title) + "\n"
	doc += lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Margin(0, 0, 1, 2).Render(tr.DropStar) + "\n"

	var (
		isError       bool
		errMsg        string
		totalTraffic  uint64
		unifiedLogs   []tunnel.RequestLog
		allDone       = true
		anyStarting   = false
		minState      = stateDone
		isDownloading bool
	)
	for _, t := range m.Tunnels {
		if t.State == stateError {
			isError = true
			errMsg = t.Err.Error()
		}
		if t.State != stateDone {
			allDone = false
		}
		if t.State != stateDone && t.State != stateError && t.State != stateConflict {
			if t.State < minState {
				minState = t.State
			}
		}
		if t.State == stateStartingProxy || t.State == stateDownloading || t.State == stateStartingTunnel || t.State == stateRegistering || t.State == stateVerifyingDNS {
			anyStarting = true
		}
		if t.State == stateDownloading {
			isDownloading = true
		}
		if t.Proxy != nil {
			totalTraffic += t.Proxy.Stats.Bytes.Load()
			unifiedLogs = append(unifiedLogs, t.Proxy.Stats.GetRecent()...)
		}
	}

	if isDownloading && !isError {
		doc := titleStyle.Render(tr.Title) + "\n"
		msgStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Margin(1, 0, 1, 2)
		doc += msgStyle.Render(m.Spinner.View()+tr.DownloadBinary) + "\n\n"
		doc += lipgloss.NewStyle().MarginLeft(2).Render(m.Progress.View()) + "\n\n"
		doc += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginLeft(2).Render(tr.PressQuit) + "\n"
		return doc
	}

	// 1. Render conflict menu
	// Check if any tunnel is in conflict state and render the interactive menu
	for i, t := range m.Tunnels {
		if t.State == stateConflict {
			doc += labelStyle.Render(tr.Version) + valStyle.Render(AppVersion) + "\n"
			doc += labelStyle.Render(tr.Status) + lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true).Render(tr.Conflict) + "\n\n"
			doc += lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("252")).Render(
				fmt.Sprintf(tr.SubdomainInUse, t.Config.Subdomain)) + "\n"
			doc += lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("241")).Render(
				fmt.Sprintf(tr.TunnelPortInfo, t.Config.Name, t.Config.Port)) + "\n\n"

			options := []string{tr.UseRandomSubdomain, tr.ExitOption}
			for j, opt := range options {
				cursor := "  "
				var optStyle lipgloss.Style
				if m.ConflictCursor == j {
					cursor = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true).Render("> ")
					optStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Bold(true)
				} else {
					optStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
				}
				doc += lipgloss.NewStyle().MarginLeft(2).Render(cursor+optStyle.Render(opt)) + "\n"
			}
			_ = i
			doc += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(0, 0, 0, 2).Render(tr.NavigateSelect) + "\n"
			return doc
		}
	}

	// 2. Render Status Line
	var statusVal string
	if isError {
		if strings.Contains(errMsg, "Cloudflare rate limited") {
			errMsg = tr.ErrCFRateLimit
		}
		statusVal = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render(tr.ErrorPrefix + errMsg)
	} else if anyStarting {
		var stateText string
		switch minState {
		case stateStartingProxy:
			stateText = tr.InitProxies
		case stateDownloading:
			stateText = tr.DownloadBinary
		case stateStartingTunnel:
			stateText = tr.EstablishTunnel
		case stateRegistering:
			stateText = tr.MapSubdomain
		case stateVerifyingDNS:
			stateText = tr.VerifyingDNS
		default:
			stateText = tr.Starting
		}
		statusVal = m.Spinner.View() + lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render(stateText)
		if minState == stateDownloading {
			statusVal += "\n\n" + lipgloss.NewStyle().MarginLeft(2).Render(m.Progress.View())
		}
	} else if allDone {
		statusVal = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true).Render(tr.ONLINE)
	} else {
		statusVal = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(tr.Waiting)
	}
	doc += labelStyle.Render(tr.Version) + valStyle.Render(AppVersion) + "\n"
	doc += labelStyle.Render(tr.Status) + statusVal + "\n"

	// 3. Render Traffic & Session (only when ONLINE)
	if allDone && !isError {
		var trafficStr string
		if totalTraffic > 1024*1024 {
			trafficStr = fmt.Sprintf("%.2f MB", float64(totalTraffic)/(1024*1024))
		} else if totalTraffic > 1024 {
			trafficStr = fmt.Sprintf("%.2f KB", float64(totalTraffic)/1024)
		} else {
			trafficStr = fmt.Sprintf("%d B", totalTraffic)
		}
		doc += labelStyle.Render(tr.Traffic) + valStyle.Render(trafficStr) + "\n"

		if !m.SessionEndsAt.IsZero() {
			remaining := time.Until(m.SessionEndsAt)
			if remaining < 0 {
				remaining = 0
			}
			hours := int(remaining.Hours())
			minutes := int(remaining.Minutes()) % 60
			seconds := int(remaining.Seconds()) % 60
			sessionStr := fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
			doc += labelStyle.Render(tr.Session) + valStyle.Render(sessionStr) + "\n\n"
		} else {
			doc += labelStyle.Render(tr.Session) + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(tr.SessionStarting) + "\n\n"
		}
	} else {
		doc += "\n"
	}

	// 4. Render Forwarding Info
	if allDone && !isError {
		isMulti := len(m.Tunnels) > 1
		for i, t := range m.Tunnels {
			label := tr.Forwarding
			if isMulti {
				if i == 0 {
					label = tr.Forwarding
				} else {
					label = ""
				}
			}

			// If PublicURL is available, show it. Otherwise, fallback to NodeURL.
			if t.PublicURL != "" {
				pubStr := urlStyle.Render(t.PublicURL) + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(fmt.Sprintf(" -> http://localhost:%d", t.Config.Port))
				if isMulti {
					nameFormat := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220")).Width(12).Render(t.Config.Name)
					doc += labelStyle.Render(label) + nameFormat + pubStr + "\n"
				} else {
					doc += labelStyle.Render(label) + pubStr + "\n"
				}
				label = ""
			} else if t.NodeURL != "" {
				nodeStr := fmt.Sprintf("%s -> http://localhost:%d", t.NodeURL, t.Config.Port)
				if isMulti {
					nameFormat := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220")).Width(12).Render(t.Config.Name)
					doc += labelStyle.Render(label) + nameFormat + valStyle.Render(nodeStr) + "\n"
				} else {
					doc += labelStyle.Render(label) + valStyle.Render(nodeStr) + "\n"
				}
				label = ""
			}
		}
		doc += "\n"
	}

	// 5. Render HTTP Logs
	sort.Slice(unifiedLogs, func(i, j int) bool {
		return unifiedLogs[i].Time.Before(unifiedLogs[j].Time)
	})
	if len(unifiedLogs) > 10 {
		unifiedLogs = unifiedLogs[len(unifiedLogs)-10:]
	}

	if len(unifiedLogs) > 0 {
		doc += lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63")).Margin(0, 0, 0, 2).Render(tr.HttpRequests) + "\n"
		for _, log := range unifiedLogs {
			var sColor lipgloss.Color
			if log.Status >= 200 && log.Status < 300 {
				sColor = lipgloss.Color("42")
			} else if log.Status >= 300 && log.Status < 400 {
				sColor = lipgloss.Color("220")
			} else {
				sColor = lipgloss.Color("196")
			}

			// Status is always 3 digits, so it aligns perfectly without padding
			statusTxt := lipgloss.NewStyle().Foreground(sColor).Render(fmt.Sprintf("%d", log.Status))
			// Method is padded to 6 chars, aligned left
			methodTxt := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("252")).Width(6).Align(lipgloss.Left).Render(log.Method)
			timeTxt := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginLeft(2).Render(log.Time.Format("15:04:05"))

			var nameTxt string
			if len(m.Tunnels) > 1 {
				nameTxt = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Width(8).Render("["+log.TunnelName+"]") + " "
			}

			path := log.Path

			// Dynamic width calculation for path capping
			fixedWidth := 30 // Time(10) + Status(4) + Method(6) + Spacing(10)
			if len(m.Tunnels) > 1 {
				fixedWidth += 9 // Name(9)
			}

			pathWidth := m.Width - fixedWidth
			if pathWidth < 20 {
				pathWidth = 40 // Default fallback if terminal is too narrow
			}

			// Cap the maximum width
			if pathWidth > 80 {
				pathWidth = 80
			}

			if len(path) > pathWidth {
				path = path[:pathWidth-3] + "..."
			}

			pathTxt := lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render(path)

			doc += fmt.Sprintf("%s %s %s %s%s\n", timeTxt, statusTxt, methodTxt, nameTxt, pathTxt)
		}
		doc += "\n"
	}

	doc += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(0, 0, 0, 2).Render(tr.PressQuit) + "\n"

	return doc
}

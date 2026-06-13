package ui

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"nipo/internal/tunnel"
)

type downloadCompleteMsg struct {
	binPath string
	err     error
}

type tunnelStartedMsg struct {
	id  int
	url string
	err error
	cmd *exec.Cmd // the exact cloudflared process we started
}

type proxyStartedMsg struct {
	id    int
	proxy *tunnel.Proxy
	err   error
}

type registerCompleteMsg struct {
	id  int
	err error
}

type heartbeatMsg struct {
	err error
}

type downloadProgressMsg float64

// downloadTask returns a Bubble Tea command that installs the tunnel binary if it is missing,
// reporting progress to the given channel.
func downloadTask(progressChan chan<- float64) tea.Cmd {
	return func() tea.Msg {
		binPath, err := tunnel.InstallIfMissingWithProgress(progressChan)
		if progressChan != nil {
			close(progressChan)
		}
		return downloadCompleteMsg{binPath: binPath, err: err}
	}
}

// listenToProgress reads from the progress channel and returns a downloadProgressMsg.
func listenToProgress(ch <-chan float64) tea.Cmd {
	return func() tea.Msg {
		percent, ok := <-ch
		if !ok {
			return nil
		}
		return downloadProgressMsg(percent)
	}
}

// startProxyTask returns a Bubble Tea command that starts the local proxy for the given port.
func startProxyTask(ctx context.Context, id int, tunnelName string, targetPort int) tea.Cmd {
	return func() tea.Msg {
		p, err := tunnel.StartProxy(ctx, tunnelName, targetPort)
		return proxyStartedMsg{id: id, proxy: p, err: err}
	}
}

// startTunnelTask returns a Bubble Tea command that launches the cloudflared daemon.
func startTunnelTask(ctx context.Context, id int, binPath string, port int, delay time.Duration) tea.Cmd {
	return func() tea.Msg {
		if delay > 0 {
			select {
			case <-ctx.Done():
				return tunnelStartedMsg{id: id, err: ctx.Err()}
			case <-time.After(delay):
			}
		}
		t, err := tunnel.StartTunnel(ctx, binPath, port)
		if err != nil {
			return tunnelStartedMsg{id: id, err: err}
		}
		return tunnelStartedMsg{id: id, url: t.URL, cmd: t.Cmd}
	}
}

// registerTask returns a Bubble Tea command that maps the subdomain on the API backend.
func registerTask(ctx context.Context, id int, backendURL, subdomain, targetURL string) tea.Cmd {
	return func() tea.Msg {
		payload := map[string]string{
			"subdomain": subdomain,
			"target":    targetURL,
		}
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return registerCompleteMsg{id: id, err: fmt.Errorf("marshal register payload: %w", err)}
		}

		client := &http.Client{Timeout: 10 * time.Second}
		var resp *http.Response
		var reqErr error

		// Retry register request up to 3 times
		for i := 0; i < 3; i++ {
			req, err := http.NewRequestWithContext(ctx, "POST", backendURL+"/register", bytes.NewBuffer(jsonData))
			if err != nil {
				return registerCompleteMsg{id: id, err: fmt.Errorf("create register request: %w", err)}
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer nipo-secret")

			resp, reqErr = client.Do(req)
			if reqErr == nil {
				break
			}
			// Wait before retry unless context is done
			select {
			case <-ctx.Done():
				return registerCompleteMsg{id: id, err: ctx.Err()}
			case <-time.After(1 * time.Second):
			}
		}

		if reqErr != nil {
			return registerCompleteMsg{id: id, err: fmt.Errorf("register subdomain %s: %w", subdomain, reqErr)}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var errResp struct {
				Error string `json:"error"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Error != "" {
				return registerCompleteMsg{id: id, err: errors.New(errResp.Error)}
			}
			return registerCompleteMsg{id: id, err: fmt.Errorf("backend returned status %d", resp.StatusCode)}
		}

		return registerCompleteMsg{id: id, err: nil}
	}
}

// UnregisterSubdomain deletes a subdomain mapping from the backend KV store.
func UnregisterSubdomain(ctx context.Context, backendURL, subdomain string) {
	if backendURL == "" || subdomain == "" {
		return
	}
	payload := map[string]string{"subdomain": subdomain}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", backendURL+"/unregister", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer nipo-secret")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
}

// heartbeatTask periodically renews TTL (10 minutes interval) for active subdomains.
func heartbeatTask(ctx context.Context, backendURL string, subdomains []string) tea.Cmd {
	return func() tea.Msg {
		select {
		case <-ctx.Done():
			return heartbeatMsg{err: ctx.Err()}
		case <-time.After(10 * time.Minute):
		}

		client := &http.Client{Timeout: 5 * time.Second}
		for _, subdomain := range subdomains {
			payload := map[string]string{"subdomain": subdomain}
			jsonData, err := json.Marshal(payload)
			if err != nil {
				continue
			}
			req, err := http.NewRequestWithContext(ctx, "POST", backendURL+"/heartbeat", bytes.NewBuffer(jsonData))
			if err != nil {
				continue
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer nipo-secret")
			resp, err := client.Do(req)
			if err == nil {
				resp.Body.Close()
			}
		}
		return heartbeatMsg{}
	}
}

package tunnel

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"time"
)

// Tunnel represents the running cloudflared process and its active public URL.
type Tunnel struct {
	URL string
	Cmd *exec.Cmd
}

// KillOrphanedTunnels kills any lingering cloudflared processes from previous runs.
func KillOrphanedTunnels() {
	// Windows: taskkill, Unix: pkill
	_ = exec.Command("taskkill", "/F", "/IM", "cloudflared.exe", "/T").Run()
	_ = exec.Command("pkill", "-f", "cloudflared").Run()
	time.Sleep(500 * time.Millisecond) // Give the OS time to release ports
}

// StartTunnel runs the cloudflared process, parses its stdout/stderr, and returns the tunnel details.
func StartTunnel(ctx context.Context, binPath string, localPort int) (*Tunnel, error) {
	cmd := exec.CommandContext(ctx, binPath, "tunnel", "--url", fmt.Sprintf("http://localhost:%d", localPort))

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("get stderr pipe for cloudflared: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start cloudflared process: %w", err)
	}

	tunnelURLChan := make(chan string, 1)

	go func() {
		scanner := bufio.NewScanner(stderr)
		re := regexp.MustCompile(`https://[a-zA-Z0-9-]+\.trycloudflare\.com`)
		urlSent := false

		for scanner.Scan() {
			line := scanner.Text()
			match := re.FindString(line)
			if !urlSent && match != "" && match != "https://api.trycloudflare.com" {
				tunnelURLChan <- match
				urlSent = true
			}
		}

		if err := scanner.Err(); err != nil {
			// We can log or just trace if needed, but scanner ending is handled below
		}

		if !urlSent {
			tunnelURLChan <- "" // Signal that no URL was found
		}
	}()

	// Wait up to 30 seconds for the tunnel URL or context cancellation
	select {
	case url := <-tunnelURLChan:
		if url == "" {
			_ = cmd.Process.Kill()
			return nil, fmt.Errorf("failed to get tunnel URL (Cloudflare rate limited - try again in a moment)")
		}
		return &Tunnel{URL: url, Cmd: cmd}, nil
	case <-ctx.Done():
		_ = cmd.Process.Kill()
		return nil, fmt.Errorf("tunnel start cancelled by context: %w", ctx.Err())
	case <-time.After(30 * time.Second):
		_ = cmd.Process.Kill()
		return nil, fmt.Errorf("timeout waiting for tunnel URL (cloudflared took too long)")
	}
}

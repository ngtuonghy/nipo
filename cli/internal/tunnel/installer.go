package tunnel

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// InstallIfMissing checks for the cloudflared binary in the user's home directory.
// If missing, it downloads it. Backward compatibility helper.
func InstallIfMissing() (string, error) {
	return InstallIfMissingWithProgress(nil)
}

// progressWriter implements io.Writer to report write progress to a channel.
type progressWriter struct {
	total      int64
	written    int64
	progressCh chan<- float64
	lastSent   float64
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.written += int64(n)
	if pw.progressCh != nil && pw.total > 0 {
		percent := float64(pw.written) / float64(pw.total)
		// Send progress if changed by at least 1% to prevent channel clogging
		if percent-pw.lastSent >= 0.01 || percent == 1.0 {
			pw.progressCh <- percent
			pw.lastSent = percent
		}
	}
	return n, nil
}

// InstallIfMissingWithProgress downloads the appropriate cloudflared binary with download progress tracking.
func InstallIfMissingWithProgress(progressChan chan<- float64) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get user home dir: %w", err)
	}

	binDir := filepath.Join(home, ".nipo", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return "", fmt.Errorf("create bin directory %s: %w", binDir, err)
	}

	var exeName string
	var downloadURL string

	switch runtime.GOOS {
	case "windows":
		exeName = "cloudflared.exe"
		if runtime.GOARCH == "amd64" {
			downloadURL = "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-windows-amd64.exe"
		}
	case "darwin":
		exeName = "cloudflared"
		downloadURL = "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-darwin-amd64.tgz"
	case "linux":
		exeName = "cloudflared"
		if runtime.GOARCH == "amd64" {
			downloadURL = "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64"
		}
	}

	if downloadURL == "" {
		return "", fmt.Errorf("unsupported OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	binPath := filepath.Join(binDir, exeName)
	if _, err := os.Stat(binPath); err == nil {
		return binPath, nil // Already installed
	}

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	resp, err := client.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("request download from %s: %w", downloadURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download from %s: HTTP %d", downloadURL, resp.StatusCode)
	}

	out, err := os.OpenFile(binPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return "", fmt.Errorf("open target binary path %s: %w", binPath, err)
	}

	pw := &progressWriter{
		total:      resp.ContentLength,
		progressCh: progressChan,
	}

	if _, err = io.Copy(io.MultiWriter(out, pw), resp.Body); err != nil {
		out.Close()
		os.Remove(binPath)
		return "", fmt.Errorf("write binary content to %s: %w", binPath, err)
	}

	if err := out.Close(); err != nil {
		os.Remove(binPath)
		return "", fmt.Errorf("finalize binary file %s: %w", binPath, err)
	}

	return binPath, nil
}

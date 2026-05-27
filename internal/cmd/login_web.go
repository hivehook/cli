package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type cliStartResp struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

type cliPollResp struct {
	Status string `json:"status"`
	APIKey string `json:"api_key"`
}

func loginWeb(ctx context.Context) error {
	host, _ := os.Hostname()
	start, err := cliPost[cliStartResp](ctx, flagEndpoint+"/api/v1/saas/auth/cli/start",
		map[string]string{"client": host})
	if err != nil {
		return fmt.Errorf("starting login: %w", err)
	}

	fmt.Fprintf(os.Stderr, "\n  Confirm this code in your browser:\n\n      %s\n\n", start.UserCode)
	if openBrowser(start.VerificationURIComplete) == nil {
		fmt.Fprintf(os.Stderr, "  Opened %s\n\n", start.VerificationURIComplete)
	} else {
		fmt.Fprintf(os.Stderr, "  Open this URL to continue:\n  %s\n\n", start.VerificationURIComplete)
	}
	fmt.Fprint(os.Stderr, "  Waiting for approval")

	interval := time.Duration(max(start.Interval, 1)) * time.Second
	deadline := time.Now().Add(time.Duration(max(start.ExpiresIn, 60)) * time.Second)

	for {
		select {
		case <-ctx.Done():
			fmt.Fprintln(os.Stderr)
			return ctx.Err()
		case <-time.After(interval):
		}
		if time.Now().After(deadline) {
			fmt.Fprintln(os.Stderr)
			return fmt.Errorf("login request expired; run 'hivehook login' again")
		}

		poll, err := cliPost[cliPollResp](ctx, flagEndpoint+"/api/v1/saas/auth/cli/poll",
			map[string]string{"device_code": start.DeviceCode})
		if err != nil {
			fmt.Fprintln(os.Stderr)
			return fmt.Errorf("polling: %w", err)
		}
		switch poll.Status {
		case "approved":
			fmt.Fprintln(os.Stderr, " done.")
			return persistLogin(ctx, poll.APIKey)
		case "denied":
			fmt.Fprintln(os.Stderr)
			return fmt.Errorf("authorization was cancelled in the browser")
		case "expired":
			fmt.Fprintln(os.Stderr)
			return fmt.Errorf("login request expired; run 'hivehook login' again")
		default:
			fmt.Fprint(os.Stderr, ".")
		}
	}
}

func cliPost[T any](ctx context.Context, url string, body any) (*T, error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		var e struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&e)
		if e.Error != "" {
			return nil, fmt.Errorf("%s", e.Error)
		}
		return nil, fmt.Errorf("server returned %d", resp.StatusCode)
	}
	var out T
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}

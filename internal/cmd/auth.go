package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newLoginCmd() *cobra.Command {
	var apiKey string
	var paste bool
	c := &cobra.Command{
		Use:   "login",
		Short: "Authenticate the CLI",
		Long:  "Opens your browser to approve this device and saves an API key to ~/.config/hivehook.\nUse --api-key to pass a key directly, or --paste to enter one at a hidden prompt.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if key := firstNonEmpty(apiKey, os.Getenv("HIVEHOOK_API_KEY")); key != "" {
				return persistLogin(cmd.Context(), key)
			}
			if paste {
				key, err := promptAPIKey()
				if err != nil {
					return err
				}
				return persistLogin(cmd.Context(), key)
			}
			return loginWeb(cmd.Context())
		},
	}
	c.Flags().StringVar(&apiKey, "api-key", "", "authenticate with this API key instead of the browser")
	c.Flags().BoolVar(&paste, "paste", false, "paste an API key at a hidden prompt instead of the browser")
	return c
}

func persistLogin(ctx context.Context, key string) error {
	flagAPIKey = key
	cl, err := client()
	if err != nil {
		return err
	}
	st, err := cl.Status.Get(ctx)
	if err != nil {
		return fmt.Errorf("could not validate the API key against %s: %w", flagEndpoint, err)
	}
	if err := saveCredentials(credentials{Endpoint: flagEndpoint, APIKey: key}); err != nil {
		return fmt.Errorf("saving credentials: %w", err)
	}
	fmt.Printf("Logged in to %s (server %s). Credentials saved.\n", flagEndpoint, st.Version)
	return nil
}

func promptAPIKey() (string, error) {
	fmt.Fprint(os.Stderr, "Create an API key in the dashboard (Settings → API keys), then paste it.\nAPI key: ")
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", fmt.Errorf("reading API key: %w", err)
	}
	key := strings.TrimSpace(string(b))
	if key == "" {
		return "", fmt.Errorf("no API key provided")
	}
	return key, nil
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove stored credentials",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := clearCredentials(); err != nil {
				return err
			}
			fmt.Println("Logged out.")
			return nil
		},
	}
}

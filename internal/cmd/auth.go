package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newLoginCmd() *cobra.Command {
	var apiKey string
	c := &cobra.Command{
		Use:   "login",
		Short: "Authenticate the CLI with an API key",
		Long:  "Stores an API key in ~/.config/hivehook so every command is authenticated. Create a key in the dashboard under Settings → API keys.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			key := apiKey
			if key == "" {
				key = os.Getenv("HIVEHOOK_API_KEY")
			}
			if key == "" {
				fmt.Fprint(os.Stderr, "Create an API key in the dashboard (Settings → API keys), then paste it.\nAPI key: ")
				b, err := term.ReadPassword(int(os.Stdin.Fd()))
				fmt.Fprintln(os.Stderr)
				if err != nil {
					return fmt.Errorf("reading API key: %w", err)
				}
				key = strings.TrimSpace(string(b))
			}
			if key == "" {
				return fmt.Errorf("no API key provided")
			}

			flagAPIKey = key
			cl, err := client()
			if err != nil {
				return err
			}
			st, err := cl.Status.Get(cmd.Context())
			if err != nil {
				return fmt.Errorf("could not validate the API key against %s: %w", flagEndpoint, err)
			}
			if err := saveCredentials(credentials{Endpoint: flagEndpoint, APIKey: key}); err != nil {
				return fmt.Errorf("saving credentials: %w", err)
			}
			fmt.Printf("Logged in to %s (server %s). Credentials saved.\n", flagEndpoint, st.Version)
			return nil
		},
	}
	c.Flags().StringVar(&apiKey, "api-key", "", "API key (prompted if omitted)")
	return c
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

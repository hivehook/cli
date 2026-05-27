package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	hivehook "github.com/hivehook/sdk-go"
	"github.com/spf13/cobra"
)

const defaultEndpoint = "https://app.hivehook.com"

var (
	flagEndpoint string
	flagAPIKey   string
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "hivehook",
		Short:         "HiveHook CLI, webhook infrastructure for modern teams",
		Long:          "Manage HiveHook from the command line: sources, destinations, subscriptions, applications, endpoints, messages, and more, over the HiveHook API.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().StringVar(&flagEndpoint, "endpoint", envOr("HIVEHOOK_ENDPOINT", "https://app.hivehook.com"), "HiveHook API endpoint (env HIVEHOOK_ENDPOINT)")
	root.PersistentFlags().StringVar(&flagAPIKey, "api-key", os.Getenv("HIVEHOOK_API_KEY"), "API key (env HIVEHOOK_API_KEY)")

	root.AddCommand(
		newLoginCmd(),
		newLogoutCmd(),
		newStatusCmd(),
		newVersionCmd(),
	)
	addResourceCommands(root)
	return root
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func client() (*hivehook.Client, error) {
	key, ep := flagAPIKey, flagEndpoint
	if key == "" {
		if c, err := loadCredentials(); err == nil && c.APIKey != "" {
			key = c.APIKey
			if c.Endpoint != "" && os.Getenv("HIVEHOOK_ENDPOINT") == "" && ep == defaultEndpoint {
				ep = c.Endpoint
			}
		}
	}
	if key == "" {
		return nil, fmt.Errorf("not logged in: run 'hivehook login' (or pass --api-key / set HIVEHOOK_API_KEY)")
	}
	return hivehook.New(hivehook.WithBaseURL(ep), hivehook.WithAPIKey(key)), nil
}

func emit(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func readBody(data, file string) ([]byte, error) {
	switch {
	case data != "":
		return []byte(data), nil
	case file == "-" || (data == "" && file == ""):
		return io.ReadAll(os.Stdin)
	default:
		return os.ReadFile(file)
	}
}

type ListParams struct {
	Limit  *int
	Search *string
}

// addCRUD registers list/get/create/update/delete; pass nil to skip an operation.
func addCRUD(
	parent *cobra.Command,
	singular string,
	list func(context.Context, ListParams) (any, error),
	get func(context.Context, string) (any, error),
	create func(context.Context, []byte) (any, error),
	update func(context.Context, string, []byte) (any, error),
	del func(context.Context, string) error,
) {
	if list != nil {
		var limit int
		var search string
		c := &cobra.Command{
			Use:   "list",
			Short: "List " + singular + "s",
			RunE: func(cmd *cobra.Command, _ []string) error {
				p := ListParams{}
				if cmd.Flags().Changed("limit") {
					p.Limit = &limit
				}
				if search != "" {
					p.Search = &search
				}
				out, err := list(cmd.Context(), p)
				if err != nil {
					return err
				}
				return emit(out)
			},
		}
		c.Flags().IntVar(&limit, "limit", 0, "max results")
		c.Flags().StringVar(&search, "search", "", "search filter")
		parent.AddCommand(c)
	}
	if get != nil {
		parent.AddCommand(&cobra.Command{
			Use:   "get <id>",
			Short: "Get a " + singular + " by ID",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				out, err := get(cmd.Context(), args[0])
				if err != nil {
					return err
				}
				return emit(out)
			},
		})
	}
	if create != nil {
		var data, file string
		c := &cobra.Command{
			Use:   "create",
			Short: "Create a " + singular + " from a JSON body (--data, --file, or stdin)",
			RunE: func(cmd *cobra.Command, _ []string) error {
				body, err := readBody(data, file)
				if err != nil {
					return err
				}
				out, err := create(cmd.Context(), body)
				if err != nil {
					return err
				}
				return emit(out)
			},
		}
		c.Flags().StringVar(&data, "data", "", "JSON body")
		c.Flags().StringVarP(&file, "file", "f", "", "JSON body file ('-' for stdin)")
		parent.AddCommand(c)
	}
	if update != nil {
		var data, file string
		c := &cobra.Command{
			Use:   "update <id>",
			Short: "Update a " + singular + " from a JSON body (--data, --file, or stdin)",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				body, err := readBody(data, file)
				if err != nil {
					return err
				}
				out, err := update(cmd.Context(), args[0], body)
				if err != nil {
					return err
				}
				return emit(out)
			},
		}
		c.Flags().StringVar(&data, "data", "", "JSON body")
		c.Flags().StringVarP(&file, "file", "f", "", "JSON body file ('-' for stdin)")
		parent.AddCommand(c)
	}
	if del != nil {
		parent.AddCommand(&cobra.Command{
			Use:   "delete <id>",
			Short: "Delete a " + singular + " by ID",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := del(cmd.Context(), args[0]); err != nil {
					return err
				}
				fmt.Fprintf(os.Stdout, "deleted %s\n", args[0])
				return nil
			},
		})
	}
}

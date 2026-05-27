package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	hivehook "github.com/hivehook/sdk-go"
	"github.com/spf13/cobra"
)

func addResourceCommands(root *cobra.Command) {
	root.AddCommand(
		sourcesCmd(),
		destinationsCmd(),
		subscriptionsCmd(),
		applicationsCmd(),
		endpointsCmd(),
		messagesCmd(),
		eventsCmd(),
		deliveriesCmd(),
		outboundDeliveriesCmd(),
		dlqCmd(),
		outboundDLQCmd(),
		apiKeysCmd(),
		alertRulesCmd(),
		bookmarksCmd(),
		eventTypeSchemasCmd(),
		transformationsCmd(),
		auditLogsCmd(),
	)
}

func sourcesCmd() *cobra.Command {
	c := &cobra.Command{Use: "sources", Short: "Manage inbound webhook sources"}
	addCRUD(c, "source",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			n, _, err := cl.Sources.List(ctx, &hivehook.ListSourcesOptions{ListOptions: listOpts(p)})
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.Sources.Get(ctx, id) }),
		bodyCreate(func(cl *hivehook.Client, ctx context.Context, in *hivehook.CreateSourceInput) (any, error) { return cl.Sources.Create(ctx, in) }),
		bodyUpdate(func(cl *hivehook.Client, ctx context.Context, id string, in *hivehook.UpdateSourceInput) (any, error) { return cl.Sources.Update(ctx, id, in) }),
		idDelete(func(cl *hivehook.Client, ctx context.Context, id string) error { return cl.Sources.Delete(ctx, id) }),
	)
	c.AddCommand(
		idActionCmd("rotate-secret", "Rotate the signing secret", idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.Sources.RotateSecret(ctx, id) })),
		idActionCmd("clear-secondary-secret", "Clear the secondary signing secret", idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.Sources.ClearSecondarySecret(ctx, id) })),
	)
	return c
}

func destinationsCmd() *cobra.Command {
	c := &cobra.Command{Use: "destinations", Short: "Manage outbound destinations"}
	addCRUD(c, "destination",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			n, _, err := cl.Destinations.List(ctx, &hivehook.ListDestinationsOptions{ListOptions: listOpts(p)})
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.Destinations.Get(ctx, id) }),
		bodyCreate(func(cl *hivehook.Client, ctx context.Context, in *hivehook.CreateDestinationInput) (any, error) { return cl.Destinations.Create(ctx, in) }),
		bodyUpdate(func(cl *hivehook.Client, ctx context.Context, id string, in *hivehook.UpdateDestinationInput) (any, error) { return cl.Destinations.Update(ctx, id, in) }),
		idDelete(func(cl *hivehook.Client, ctx context.Context, id string) error { return cl.Destinations.Delete(ctx, id) }),
	)
	c.AddCommand(idActionCmd("rotate-secret", "Rotate the signing secret", idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.Destinations.RotateSecret(ctx, id) })))
	return c
}

func subscriptionsCmd() *cobra.Command {
	c := &cobra.Command{Use: "subscriptions", Short: "Manage source-to-destination routing rules"}
	addCRUD(c, "subscription",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			n, _, err := cl.Subscriptions.List(ctx, &hivehook.ListSubscriptionsOptions{ListOptions: listOpts(p)})
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.Subscriptions.Get(ctx, id) }),
		bodyCreate(func(cl *hivehook.Client, ctx context.Context, in *hivehook.CreateSubscriptionInput) (any, error) { return cl.Subscriptions.Create(ctx, in) }),
		bodyUpdate(func(cl *hivehook.Client, ctx context.Context, id string, in *hivehook.UpdateSubscriptionInput) (any, error) { return cl.Subscriptions.Update(ctx, id, in) }),
		idDelete(func(cl *hivehook.Client, ctx context.Context, id string) error { return cl.Subscriptions.Delete(ctx, id) }),
	)
	return c
}

func applicationsCmd() *cobra.Command {
	c := &cobra.Command{Use: "applications", Short: "Manage outbound applications"}
	addCRUD(c, "application",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			opts := listOpts(p)
			n, _, err := cl.Applications.List(ctx, &opts)
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.Applications.Get(ctx, id) }),
		bodyCreate(func(cl *hivehook.Client, ctx context.Context, in *hivehook.CreateApplicationInput) (any, error) { return cl.Applications.Create(ctx, in) }),
		bodyUpdate(func(cl *hivehook.Client, ctx context.Context, id string, in *hivehook.UpdateApplicationInput) (any, error) { return cl.Applications.Update(ctx, id, in) }),
		idDelete(func(cl *hivehook.Client, ctx context.Context, id string) error { return cl.Applications.Delete(ctx, id) }),
	)
	return c
}

func endpointsCmd() *cobra.Command {
	c := &cobra.Command{Use: "endpoints", Short: "Manage application endpoints"}
	addCRUD(c, "endpoint",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			n, _, err := cl.Endpoints.List(ctx, &hivehook.ListEndpointsOptions{ListOptions: listOpts(p)})
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.Endpoints.Get(ctx, id) }),
		bodyCreate(func(cl *hivehook.Client, ctx context.Context, in *hivehook.CreateEndpointInput) (any, error) { return cl.Endpoints.Create(ctx, in) }),
		bodyUpdate(func(cl *hivehook.Client, ctx context.Context, id string, in *hivehook.UpdateEndpointInput) (any, error) { return cl.Endpoints.Update(ctx, id, in) }),
		idDelete(func(cl *hivehook.Client, ctx context.Context, id string) error { return cl.Endpoints.Delete(ctx, id) }),
	)
	c.AddCommand(idActionCmd("rotate-secret", "Rotate the signing secret", idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.Endpoints.RotateSecret(ctx, id) })))
	return c
}

func messagesCmd() *cobra.Command {
	c := &cobra.Command{Use: "messages", Short: "Send and inspect outbound messages"}
	addCRUD(c, "message",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			n, _, err := cl.Messages.List(ctx, &hivehook.ListMessagesOptions{ListOptions: listOpts(p)})
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.Messages.Get(ctx, id) }),
		nil, nil, nil,
	)
	c.AddCommand(
		bodyActionCmd("send", "Send a message (JSON body)", bodyAction(func(cl *hivehook.Client, ctx context.Context, in *hivehook.SendMessageInput) (any, error) { return cl.Messages.Send(ctx, in) })),
		bodyActionCmd("broadcast", "Broadcast to every endpoint in an application (JSON body)", bodyAction(func(cl *hivehook.Client, ctx context.Context, in *hivehook.BroadcastMessageInput) (any, error) { return cl.Messages.Broadcast(ctx, in) })),
		bodyActionCmd("send-dynamic", "Send to an ad-hoc URL (JSON body)", bodyAction(func(cl *hivehook.Client, ctx context.Context, in *hivehook.SendDynamicMessageInput) (any, error) { return cl.Messages.SendDynamic(ctx, in) })),
	)
	return c
}

func eventsCmd() *cobra.Command {
	c := &cobra.Command{Use: "events", Short: "Inspect ingested events"}
	addCRUD(c, "event",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			n, _, err := cl.Events.List(ctx, &hivehook.ListEventsOptions{ListOptions: listOpts(p)})
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.Events.Get(ctx, id) }),
		nil, nil, nil,
	)
	return c
}

func deliveriesCmd() *cobra.Command {
	c := &cobra.Command{Use: "deliveries", Short: "Inspect inbound delivery attempts"}
	addCRUD(c, "delivery",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			n, _, err := cl.Deliveries.List(ctx, &hivehook.ListDeliveriesOptions{ListOptions: listOpts(p)})
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.Deliveries.Get(ctx, id) }),
		nil, nil, nil,
	)
	return c
}

func outboundDeliveriesCmd() *cobra.Command {
	c := &cobra.Command{Use: "outbound-deliveries", Short: "Inspect outbound delivery attempts"}
	addCRUD(c, "outbound delivery",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			n, _, err := cl.OutboundDeliveries.List(ctx, &hivehook.ListOutboundDeliveriesOptions{ListOptions: listOpts(p)})
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.OutboundDeliveries.Get(ctx, id) }),
		nil, nil, nil,
	)
	return c
}

func dlqCmd() *cobra.Command {
	c := &cobra.Command{Use: "dlq", Short: "Inbound dead-letter queue"}
	addCRUD(c, "DLQ entry",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			n, _, err := cl.DLQ.List(ctx, &hivehook.ListDLQOptions{ListOptions: listOpts(p)})
			return n, err
		}),
		nil, nil, nil, nil,
	)
	c.AddCommand(idActionVoidCmd("replay", "Replay a DLQ entry", func(cl *hivehook.Client, ctx context.Context, id string) error {
		return cl.DLQ.Replay(ctx, id)
	}))
	return c
}

func outboundDLQCmd() *cobra.Command {
	c := &cobra.Command{Use: "outbound-dlq", Short: "Outbound dead-letter queue"}
	addCRUD(c, "outbound DLQ entry",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			n, _, err := cl.OutboundDLQ.List(ctx, &hivehook.ListOutboundDLQOptions{ListOptions: listOpts(p)})
			return n, err
		}),
		nil, nil, nil, nil,
	)
	c.AddCommand(idActionVoidCmd("replay", "Replay an outbound DLQ entry", func(cl *hivehook.Client, ctx context.Context, id string) error {
		return cl.OutboundDLQ.Replay(ctx, id)
	}))
	return c
}

func apiKeysCmd() *cobra.Command {
	c := &cobra.Command{Use: "api-keys", Short: "Manage API keys"}
	addCRUD(c, "API key",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			opts := listOpts(p)
			n, _, err := cl.APIKeys.List(ctx, &opts)
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.APIKeys.Get(ctx, id) }),
		bodyCreate(func(cl *hivehook.Client, ctx context.Context, in *hivehook.CreateAPIKeyInput) (any, error) { return cl.APIKeys.Create(ctx, in) }),
		nil, nil,
	)
	c.AddCommand(idActionVoidCmd("revoke", "Revoke an API key", func(cl *hivehook.Client, ctx context.Context, id string) error { return cl.APIKeys.Revoke(ctx, id) }))
	return c
}

func alertRulesCmd() *cobra.Command {
	c := &cobra.Command{Use: "alert-rules", Short: "Manage DLQ and delivery alert rules"}
	addCRUD(c, "alert rule",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			n, _, err := cl.AlertRules.List(ctx, &hivehook.ListAlertRulesOptions{ListOptions: listOpts(p)})
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.AlertRules.Get(ctx, id) }),
		bodyCreate(func(cl *hivehook.Client, ctx context.Context, in *hivehook.CreateAlertRuleInput) (any, error) { return cl.AlertRules.Create(ctx, in) }),
		bodyUpdate(func(cl *hivehook.Client, ctx context.Context, id string, in *hivehook.UpdateAlertRuleInput) (any, error) { return cl.AlertRules.Update(ctx, id, in) }),
		idDelete(func(cl *hivehook.Client, ctx context.Context, id string) error { return cl.AlertRules.Delete(ctx, id) }),
	)
	c.AddCommand(idActionVoidCmd("test", "Fire a test alert", func(cl *hivehook.Client, ctx context.Context, id string) error { return cl.AlertRules.Test(ctx, id) }))
	return c
}

func bookmarksCmd() *cobra.Command {
	c := &cobra.Command{Use: "bookmarks", Short: "Manage event bookmarks"}
	addCRUD(c, "bookmark",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			n, _, err := cl.Bookmarks.List(ctx, &hivehook.ListBookmarksOptions{ListOptions: listOpts(p)})
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.Bookmarks.Get(ctx, id) }),
		nil, nil,
		idDelete(func(cl *hivehook.Client, ctx context.Context, id string) error { return cl.Bookmarks.Delete(ctx, id) }),
	)

	var eventID, name, notes string
	create := &cobra.Command{
		Use:   "create",
		Short: "Bookmark an event",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cl, err := client()
			if err != nil {
				return err
			}
			b, err := cl.Bookmarks.Create(cmd.Context(), eventID, optStr(name), optStr(notes))
			if err != nil {
				return err
			}
			return emit(b)
		},
	}
	create.Flags().StringVar(&eventID, "event-id", "", "event ID to bookmark (required)")
	create.Flags().StringVar(&name, "name", "", "bookmark name")
	create.Flags().StringVar(&notes, "notes", "", "notes")
	_ = create.MarkFlagRequired("event-id")
	c.AddCommand(create)
	return c
}

func eventTypeSchemasCmd() *cobra.Command {
	c := &cobra.Command{Use: "event-type-schemas", Short: "Manage event-type JSON schemas"}
	addCRUD(c, "event type schema",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			opts := listOpts(p)
			n, _, err := cl.EventTypeSchemas.List(ctx, &opts)
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.EventTypeSchemas.Get(ctx, id) }),
		bodyCreate(func(cl *hivehook.Client, ctx context.Context, in *hivehook.CreateEventTypeSchemaInput) (any, error) { return cl.EventTypeSchemas.Create(ctx, in) }),
		bodyUpdate(func(cl *hivehook.Client, ctx context.Context, id string, in *hivehook.UpdateEventTypeSchemaInput) (any, error) { return cl.EventTypeSchemas.Update(ctx, id, in) }),
		idDelete(func(cl *hivehook.Client, ctx context.Context, id string) error { return cl.EventTypeSchemas.Delete(ctx, id) }),
	)
	return c
}

func transformationsCmd() *cobra.Command {
	c := &cobra.Command{Use: "transformations", Short: "Manage payload transformations"}
	addCRUD(c, "transformation",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			n, _, err := cl.Transformations.List(ctx, &hivehook.ListTransformationsOptions{ListOptions: listOpts(p)})
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.Transformations.Get(ctx, id) }),
		bodyCreate(func(cl *hivehook.Client, ctx context.Context, in *hivehook.CreateTransformationInput) (any, error) { return cl.Transformations.Create(ctx, in) }),
		bodyUpdate(func(cl *hivehook.Client, ctx context.Context, id string, in *hivehook.UpdateTransformationInput) (any, error) { return cl.Transformations.Update(ctx, id, in) }),
		idDelete(func(cl *hivehook.Client, ctx context.Context, id string) error { return cl.Transformations.Delete(ctx, id) }),
	)
	c.AddCommand(bodyActionCmd("test", "Test a transformation against a sample event (JSON body)", bodyAction(func(cl *hivehook.Client, ctx context.Context, in *hivehook.TestTransformationInput) (any, error) { return cl.Transformations.Test(ctx, in) })))
	return c
}

func auditLogsCmd() *cobra.Command {
	c := &cobra.Command{Use: "audit-logs", Short: "Inspect audit log entries"}
	addCRUD(c, "audit log",
		listOf(func(cl *hivehook.Client, ctx context.Context, p ListParams) (any, error) {
			n, _, err := cl.AuditLogs.List(ctx, &hivehook.ListAuditLogsOptions{ListOptions: listOpts(p)})
			return n, err
		}),
		idGet(func(cl *hivehook.Client, ctx context.Context, id string) (any, error) { return cl.AuditLogs.Get(ctx, id) }),
		nil, nil, nil,
	)
	return c
}

func listOpts(p ListParams) hivehook.ListOptions {
	return hivehook.ListOptions{Limit: p.Limit, Search: p.Search}
}

func optStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// The adapters below inject an authenticated client and decode JSON bodies so
// each resource closure is a single SDK call.

func listOf(fn func(*hivehook.Client, context.Context, ListParams) (any, error)) func(context.Context, ListParams) (any, error) {
	return func(ctx context.Context, p ListParams) (any, error) {
		cl, err := client()
		if err != nil {
			return nil, err
		}
		return fn(cl, ctx, p)
	}
}

func idGet(fn func(*hivehook.Client, context.Context, string) (any, error)) func(context.Context, string) (any, error) {
	return func(ctx context.Context, id string) (any, error) {
		cl, err := client()
		if err != nil {
			return nil, err
		}
		return fn(cl, ctx, id)
	}
}

func bodyCreate[T any](fn func(*hivehook.Client, context.Context, *T) (any, error)) func(context.Context, []byte) (any, error) {
	return func(ctx context.Context, body []byte) (any, error) {
		cl, err := client()
		if err != nil {
			return nil, err
		}
		in, err := decode[T](body)
		if err != nil {
			return nil, err
		}
		return fn(cl, ctx, in)
	}
}

func bodyUpdate[T any](fn func(*hivehook.Client, context.Context, string, *T) (any, error)) func(context.Context, string, []byte) (any, error) {
	return func(ctx context.Context, id string, body []byte) (any, error) {
		cl, err := client()
		if err != nil {
			return nil, err
		}
		in, err := decode[T](body)
		if err != nil {
			return nil, err
		}
		return fn(cl, ctx, id, in)
	}
}

func bodyAction[T any](fn func(*hivehook.Client, context.Context, *T) (any, error)) func(context.Context, []byte) (any, error) {
	return bodyCreate(fn)
}

func idDelete(fn func(*hivehook.Client, context.Context, string) error) func(context.Context, string) error {
	return func(ctx context.Context, id string) error {
		cl, err := client()
		if err != nil {
			return err
		}
		return fn(cl, ctx, id)
	}
}

func decode[T any](body []byte) (*T, error) {
	var in T
	if err := json.Unmarshal(body, &in); err != nil {
		return nil, fmt.Errorf("invalid JSON body: %w", err)
	}
	return &in, nil
}

func idActionCmd(use, short string, fn func(context.Context, string) (any, error)) *cobra.Command {
	return &cobra.Command{
		Use:   use + " <id>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := fn(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return emit(out)
		},
	}
}

func idActionVoidCmd(use, short string, fn func(*hivehook.Client, context.Context, string) error) *cobra.Command {
	return &cobra.Command{
		Use:   use + " <id>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl, err := client()
			if err != nil {
				return err
			}
			if err := fn(cl, cmd.Context(), args[0]); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "ok: %s\n", args[0])
			return nil
		},
	}
}

func bodyActionCmd(use, short string, fn func(context.Context, []byte) (any, error)) *cobra.Command {
	var data, file string
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := readBody(data, file)
			if err != nil {
				return err
			}
			out, err := fn(cmd.Context(), body)
			if err != nil {
				return err
			}
			return emit(out)
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON body")
	c.Flags().StringVarP(&file, "file", "f", "", "JSON body file ('-' for stdin)")
	return c
}

# ⬢ HiveHook CLI

The command-line interface for [HiveHook](https://hivehook.com), webhook infrastructure for modern teams. Manage sources, destinations, subscriptions, applications, endpoints, messages, and more from your terminal, over the HiveHook API.

## Install

**macOS / Linux**

```sh
curl -fsSL https://hivehook.com/install.sh | sh
```

**Windows** (PowerShell)

```powershell
irm https://hivehook.com/install.ps1 | iex
```

**Go**

```sh
go install github.com/hivehook/cli/cmd/hivehook@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/hivehook/cli/releases). The CLI is a single static binary with no runtime dependencies.

## Quickstart

Authenticate with an API key (create one in the dashboard under **Settings → API keys**):

```sh
hivehook login
```

Then use it:

```sh
hivehook status
hivehook sources list
hivehook sources create --data '{"slug":"stripe","providerType":"stripe","secret":"whsec_..."}'
hivehook events list --limit 20
```

Every command prints JSON, so it pipes into `jq` and scripts.

## Authentication

Credentials resolve in this order:

1. `--api-key` flag
2. `HIVEHOOK_API_KEY` environment variable
3. Stored credentials from `hivehook login` (`~/.config/hivehook/credentials.json`, mode `0600`)

Point the CLI at a different endpoint with `--endpoint` or `HIVEHOOK_ENDPOINT` (default `https://app.hivehook.com`).

## Commands

Resource groups follow a consistent `list` / `get` / `create` / `update` / `delete` shape, plus resource-specific actions:

| Group | Notable actions |
|---|---|
| `sources` | `rotate-secret`, `clear-secondary-secret` |
| `destinations` | `rotate-secret` |
| `subscriptions` | |
| `applications` | |
| `endpoints` | `rotate-secret` |
| `messages` | `send`, `broadcast`, `send-dynamic` |
| `events`, `deliveries`, `outbound-deliveries` | read-only |
| `dlq`, `outbound-dlq` | `replay` |
| `api-keys` | `revoke` |
| `alert-rules` | `test` |
| `bookmarks` | |
| `event-type-schemas` | |
| `transformations` | `test` |
| `audit-logs` | read-only |

`create` and `update` read a JSON body from `--data`, a file (`--file`/`-f`), or stdin:

```sh
echo '{"slug":"github","providerType":"github"}' | hivehook sources create
hivehook destinations update dst_9a2f -f dest.json
```

Run `hivehook <group> --help` for the full list of subcommands and flags.

## License

See [LICENSE](LICENSE).

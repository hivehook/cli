// Command hivehook is the official CLI for HiveHook, talking to the hosted API
// (or any HiveHook endpoint) over the GraphQL admin API with an API key.
package main

import (
	"fmt"
	"os"

	"github.com/hivehook/cli/internal/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

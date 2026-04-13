// TODO: REPLACE ALL io.ReadAll WITH STREAMING METHODS

package main

import (
	"context"
	"os"

	"github.com/mclucy/lucy/cmd"
	"github.com/mclucy/lucy/logger"
)

func main() {
	defer logger.DumpHistory() // Whether DumpHistory actually does anything depend on the flag.
	// Keep the literal "--" case alive for completion before urfave/cli sees it.
	args := cmd.NormalizeCompletionArgs(os.Args)
	os.Args = args
	if err := cmd.Cli.Run(context.Background(), args); err != nil {
		// Error is already reported by decoratorLogAndExitOnError in commands.
		// Reporting here again would cause duplicate error messages.
	}
}

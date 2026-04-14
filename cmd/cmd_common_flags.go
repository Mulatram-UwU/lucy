package cmd

import (
	"errors"

	"github.com/mclucy/lucy/types"
	"github.com/spf13/cobra"

	"github.com/urfave/cli/v3"
)

const (
	flagJsonName      = "json"
	flagLongName      = "long"
	flagNoStyleName   = "no-style"
	flagSourceName    = "source"
	flagLogFileName   = "log-file"
	flagPrintLogsName = "print-logs"
	flagDebugName     = "debug"
	flagDumpLogsName  = "dump-logs"
)

// addJsonFlag adds the --json flag to a command.
func addJsonFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(flagJsonName, false, "Print raw JSON response")
}

// addLongFlag adds the --long/-l flag to a command.
func addLongFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP(flagLongName, "l", false, "Show hidden or collapsed output")
}

// addNoStyleFlag adds the --no-style flag to a command (local, not persistent).
func addNoStyleFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(flagNoStyleName, false, "Disable colored and styled output")
}

// addSourceFlag adds the --source/-s flag to a command.
// Validation of the source value is done in PreRunE of each command.
func addSourceFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(flagSourceName, "s", "", "To fetch info from SOURCE")
}

// validateSourceFlag validates the --source flag value.
// Returns an error if the source is not recognized.
func validateSourceFlag(cmd *cobra.Command) error {
	source, _ := cmd.Flags().GetString(flagSourceName)
	if source != "" && types.ParseSource(source) == types.SourceUnknown {
		return errors.New("unknown source " + source)
	}
	return nil
}

// Legacy urfave/cli flag variables - kept for backward compatibility
// with files that haven't been migrated to Cobra yet.
// These will be removed once all cmd_*.go files are migrated.
var (
	flagJsonOutput = &cli.BoolFlag{
		Name:  flagJsonName,
		Usage: "Print raw JSON response",
		Value: false,
	}

	flagLongOutput = &cli.BoolFlag{
		Name:    flagLongName,
		Usage:   "Show hidden or collapsed output",
		Value:   false,
		Aliases: []string{"l"},
	}

	flagSource = &cli.StringFlag{
		Name:    flagSourceName,
		Aliases: []string{"s"},
		Usage:   "To fetch info from `SOURCE`",
		Value:   "",
		Validator: func(s string) error {
			if types.ParseSource(s) == types.SourceUnknown {
				return errors.New("unknown source " + s)
			}
			return nil
		},
	}

	flagNoStyle = &cli.BoolFlag{
		Name:  flagNoStyleName,
		Usage: "Disable colored and styled output",
		Value: false,
	}
)

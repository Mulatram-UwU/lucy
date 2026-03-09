package cmd

import (
	"context"

	"github.com/mclucy/lucy/tools"

	"github.com/mclucy/lucy/logger"

	"github.com/urfave/cli/v3"
)

// decoratorBaseCommandFlags provides the base command `lucy` some necessary
// flag actions.
func decoratorBaseCommandFlags(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Bool(flagNoStyleName) {
			tools.TurnOffStyles()
		}
		if cmd.Bool(flagLogFileName) {
			println("Log file at", logger.GetLogFile().Name())
		}
		return f(ctx, cmd)
	}
}

// decoratorGlobalFlags is a higher-order function that appends global flag actions
// to the action function.
func decoratorGlobalFlags(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Bool(flagPrintLogsName) {
			logger.EnablePrintLogs()
		}
		if cmd.Bool(flagDebugName) {
			logger.EnableDebug()
		}
		if cmd.Bool(flagDumpLogsName) {
			logger.EnableDumpHistory()
		}
		if cmd.Bool(flagNoStyleName) {
			tools.TurnOffStyles()
		}
		return f(ctx, cmd)
	}
}

// decoratorHelpAndExitOnNoArg is a higher-order function that takes a cli.ActionFunc and
// returns a cli.ActionFunc that prints help and exit when there's no args specified.
//
// This function is not necessarily applicable to every action function, as some
// sub-commands are expected to have no args, e.g., `lucy status`.
func decoratorHelpAndExitOnNoArg(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() == 0 {
			cli.ShowSubcommandHelpAndExit(cmd, 0)
		}
		return f(ctx, cmd)
	}
}

// hasAnyFlagSet returns true if any flag for this command has a non-zero value.
// It attempts to infer whether a flag was set by checking common flag types
// (bool, string, int) via their zero values.
func hasAnyFlagSet(cmd *cli.Command) bool {
	for _, name := range cmd.FlagNames() {
		// bool flags: default is typically false; true implies set
		if cmd.Bool(name) {
			return true
		}
		// string flags: non-empty implies set
		if cmd.String(name) != "" {
			return true
		}
		// int flags: non-zero implies set
		if cmd.Int(name) != 0 {
			return true
		}
	}
	return false
}

// decoratorHelpAndExitOnNoFlag is similar to decoratorHelpAndExitOnNoArg, but
// it checks for flags instead of args. This is useful for commands that
// require at least one flag.
func decoratorHelpAndExitOnNoFlag(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if !hasAnyFlagSet(cmd) {
			cli.ShowSubcommandHelpAndExit(cmd, 0)
		}
		return f(ctx, cmd)
	}
}

func decoratorLogAndExitOnError(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		err := f(ctx, cmd)
		if err != nil {
			logger.ReportError(err)
			return err
		}
		return nil
	}
}

// decoratorHelpAndExitOnError exits with an error code and prints the help
//
// This means, with this decorator, you MUST NOT throw unexpected errors
// in your action function, as it will be caught and printed to the
// user.
//
// ONLY errors readable by the user should be thrown.
//
// Comparingly, decoratorLogAndExitOnError is more suitable for
// most of the action functions.
func decoratorHelpAndExitOnError(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		err := f(ctx, cmd)
		if err != nil {
			cli.ShowSubcommandHelpAndExit(cmd, 1)
		}
		return err
	}
}

package cmd

import (
	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/site"

	"github.com/spf13/cobra"
)

func xdebug(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "xdebug [on/off]",
		Short: "Turns Xdebug on or off without having to stop and start the site.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			status := kanaSite.IsXdebugRunning(consoleOutput)

			consoleOutput.Println(outputXdebugStatus(status))
		},
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	onCommand := &cobra.Command{
		Use:   "on",
		Short: "Starts Xdebug to aid in PHP debugging",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			err = kanaSite.StartXdebug(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			status := kanaSite.IsXdebugRunning(consoleOutput)

			consoleOutput.Println(outputXdebugStatus(status))
		},
	}

	commandsRequiringSite = append(commandsRequiringSite, onCommand.Use)

	offCommand := &cobra.Command{
		Use:   "off",
		Short: "Stops xdebug and removes its configuration",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			status := kanaSite.IsXdebugRunning(consoleOutput)

			if status {
				err = kanaSite.StopXdebug(consoleOutput)
				if err != nil {
					consoleOutput.Error(err)
				}

				status = kanaSite.IsXdebugRunning(consoleOutput)
			}

			consoleOutput.Println(outputXdebugStatus(status))
		},
	}

	commandsRequiringSite = append(commandsRequiringSite, offCommand.Use)

	cmd.AddCommand(
		onCommand,
		offCommand,
	)

	return cmd
}

func outputXdebugStatus(status bool) string {
	displayStatus := "off"

	if status {
		displayStatus = "on"
	}

	return displayStatus
}

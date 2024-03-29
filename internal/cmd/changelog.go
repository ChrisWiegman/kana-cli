package cmd

import (
	"github.com/ChrisWiegman/kana/internal/console"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func changelog(consoleOutput *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changelog",
		Short: "Open Kana's changelog in your browser",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := browser.OpenURL("https://github.com/ChrisWiegman/kana/releases")
			if err != nil {
				consoleOutput.Error(err)
			}

			consoleOutput.Success("The Kana changelog has been opened in your default browser.")
		},
	}

	return cmd
}

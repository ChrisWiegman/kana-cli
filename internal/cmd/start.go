package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/internal/setup"
	"github.com/ChrisWiegman/kana/internal/site"
	"github.com/ChrisWiegman/kana/internal/traefik"

	"github.com/spf13/cobra"
)

var FlagXdebug bool

func newStartCommand(appConfig config.AppConfig) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a new environment in the local folder.",
		Run: func(cmd *cobra.Command, args []string) {
			runStart(cmd, args, appConfig)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := setup.SetupApp(appConfig)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVar(&FlagXdebug, "xdebug", false, "Enable Xdebug when starting the container.")

	return cmd

}

func runStart(cmd *cobra.Command, args []string, appConfig config.AppConfig) {

	site, err := site.NewSite(appConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Starting development site: %s\n", site.GetURL(false))

	traefikClient, err := traefik.NewTraefik(appConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = traefikClient.StartTraefik()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = site.StartWordPress()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = site.VerifySite()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = site.InstallWordPress()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if FlagXdebug {
		_, err = site.InstallXdebug()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	err = site.OpenSite()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

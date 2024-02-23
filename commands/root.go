// Copyright (c) 2024 Owen Waller. All rights reserved.
package commands

import (
	//"fmt"

	"github.com/owenwaller/emailformgateway/server"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "emailformgateway <[flags]>",
	Short: "A simple gateway that sends emails",
	Long: `Longer description..
            feel free to use a few lines here.
            `,
	RunE: rootCmd, // function t oRun but returns an error
}

var configFilename string
var host string
var port string
var route string
var domain string

func init() {
	// set up the flags
	RootCmd.PersistentFlags().StringVarP(&configFilename, "config", "c", "", "The name of the configuration file to use, minus the TOML extension")
	RootCmd.PersistentFlags().StringVarP(&host, "host", "h", "localhost", "The hostname of the gateway server")
	RootCmd.PersistentFlags().StringVarP(&port, "port", "p", "9301", "The port the gateway server listens on")
	RootCmd.PersistentFlags().StringVarP(&route, "route", "r", "/", "The URL path that the form data is POSTed to")
	RootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "example.com", "The domain the email form gateway receives form data from")
	// Remove the default -h help shorthand as we want this for hostname as above.
	RootCmd.PersistentFlags().BoolP("help", "", false, "help for this command")
}

func rootCmd(cmd *cobra.Command, args []string) error {
	s := server.NewServer(host, port, domain)
	s.SetRouteHandler(route)
	if err := s.ReadConfig(configFilename); err != nil {
		return err
	}
	return s.Start()
}

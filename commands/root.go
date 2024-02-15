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
	Run: rootCmd,
}

var configFilename string
var host string
var port string
var route string

func init() {
	// set up the flags
	RootCmd.PersistentFlags().StringVarP(&configFilename, "config", "c", "", "The name of the configuration file to use, minus the TOML extension")
	RootCmd.PersistentFlags().StringVarP(&host, "host", "h", "localhost", "The hostname of the gateway server")
	RootCmd.PersistentFlags().StringVarP(&port, "port", "p", "9301", "The port the gateway server listens on")
	RootCmd.PersistentFlags().StringVarP(&route, "route", "r", "/", "The URL path that the form data is POSTed to")

}

func rootCmd(cmd *cobra.Command, args []string) {
	s := server.NewServer(host, port, route)
	s.Start(configFilename)
}

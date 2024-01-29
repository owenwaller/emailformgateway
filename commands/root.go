package commands

import (
	//"fmt"
	"os"

	"github.com/owenwaller/emailformgateway/config"
	conf "github.com/owenwaller/emailformgateway/config"
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

func init() {
	// set up the flags
	RootCmd.PersistentFlags().StringVarP(&configFilename, "config", "c", "", "Help for config flag")
}

func rootCmd(cmd *cobra.Command, args []string) {
	//fmt.Printf("Config value = \"%s\"\n", configFilename)
	// err := conf.ReadConfig(configFilename, conf.GetConfig())
	// if err != nil {
	// 	//fmt.Printf("%v\n", err)
	// 	os.Exit(-1)
	// }

	//configFileName := DefaultConfigFilename
	c, err := config.ReadConfig(configFilename)
	if err != nil {
		os.Exit(-1)
	}

	// set up templates
	conf.SetUpTemplates()
	server.Start(c)
}

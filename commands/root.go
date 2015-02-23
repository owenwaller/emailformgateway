package commands

import (
	"fmt"
	"os"

	conf "github.com/owenwaller/emailformgateway/config"
	"github.com/owenwaller/emailformgateway/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultConfigFilename = "config"
)

var RootCmd = &cobra.Command{
	Use:   "emailformgateway <[flags]>",
	Short: "A simple gateway that sends emails",
	Long: `Longer description..
            feel free to use a few lines here.
            `,
	Run: rootCmd,
}

var config conf.Config
var configFilename string

func init() {
	// set up the flags
	RootCmd.PersistentFlags().StringVarP(&configFilename, "config", "c", "", "Help for config flag")
}

func rootCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("Config value = \"%s\"\n", configFilename)
	setConfigFile(configFilename)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Could not read in config. Error \"%s\"\n", err)
		os.Exit(-1)
	}

	fmt.Printf("All Keys: \"%#v\"\n", viper.AllKeys())
	err = viper.Marshal(&config)
	if err != nil {
		fmt.Printf("Failed to marshal config. Error: \"%s\"\n", err)
		os.Exit(-1)
	}

	server.Start(&config)
	//viper.Debug()
}

func setConfigFile(configFilename string) {
	if configFilename == "" {
		configFilename = DefaultConfigFilename
		// viper automatically searches the CWD - see findConfigFile
		viper.AddConfigPath("/etc/emailformgateway/")
		viper.SetConfigName(configFilename)
	} else {
		// check of the file exists
		_, err := os.Lstat(configFilename)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(-1)
		}
		viper.SetConfigFile(configFilename)
	}
}

package commands

import (
	"fmt"
	"os"

	conf "github.com/owenwaller/emailformgateway/config"
	"github.com/owenwaller/emailformgateway/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ConfigReadError struct {
	msg string
	err string
}

type ConfigMarshalError struct {
	msg string
	err string
}

func NewConfigReadError(msg string, err string) ConfigReadError {
	return ConfigReadError{msg, err}
}

func (e ConfigReadError) Error() string {
	return fmt.Sprintf("%s, %s", e.msg, e.err)
}

func NewConfigMarshalError(msg string, err string) ConfigMarshalError {
	return ConfigMarshalError{msg, err}
}

func (e ConfigMarshalError) Error() string {
	return fmt.Sprintf("%s, %s", e.msg, e.err)
}

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
	err := ReadConfig(configFilename, &config)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(-1)
	}
	server.Start(&config)
	//viper.Debug()
}

func setConfigFile(configFilename string) error {
	if configFilename == "" {
		configFilename = DefaultConfigFilename
		// viper automatically searches the CWD - see findConfigFile
		// need to make the path work on Windows
		viper.AddConfigPath("/etc/emailformgateway/")
		viper.SetConfigName(configFilename)
	} else {
		// check of the file exists
		_, err := os.Lstat(configFilename)
		if err != nil {
			return err // this will be of type os.PathError
		}
		viper.SetConfigFile(configFilename)
	}
	return nil
}

func ReadConfig(filename string, c *conf.Config) error {
	err := setConfigFile(filename)
	if err != nil {
		return err // returns an os.PathError
	}
	err = viper.ReadInConfig()
	if err != nil {
		re := NewConfigReadError("Could not read in config.", err.Error())
		return re
	}

	fmt.Printf("All Keys: \"%#v\"\n", viper.AllKeys())
	err = viper.Marshal(c)
	if err != nil {
		me := NewConfigMarshalError("Failed to marshal config.", err.Error())
		return me
	}
	return err // will be nil
}

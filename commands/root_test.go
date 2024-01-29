package commands

// import (
// 	"os"
// 	"testing"

// 	conf "github.com/owenwaller/emailformgateway/config"
// 	"github.com/spf13/viper"
// )

// func TestReadConfig(t *testing.T) {
// 	var filename = "this/file/does/not/exist.toml"
// 	//var c conf.Config
// 	//err := conf.ReadConfig(filename, &c)
// 	_, err := conf.ReadConfigNew(filename)
// 	if err != nil {
// 		_, ok := err.(conf.ConfigReadError)
// 		if !ok {
// 			t.Fatalf("Error expected error was not of type config.ConfigReadError. Got %#v\n", err)
// 		}
// 	}
// 	// now try with a real file - via the ENV
// 	// filename = os.Getenv("TEST_CONFIG_FILE")
// 	// if filename == "" {
// 	// 	t.Fatalf("Required enviromental variable \"TEST_CONFIG_FILE\" not set.\nIt should be the absolute path of the config file.")
// 	// }

// 	dir, ok := os.LookupEnv("TEST_CONFIG_DIR")
// 	if ok {
// 		viper.AddConfigPath(dir)
// 	}
// 	filename = os.Getenv("TEST_CONFIG_FILE")
// 	// if filename is unset then we'll use the default so this is safe
// 	_, err = conf.ReadConfigNew(filename)

// 	if err != nil {
// 		switch errtype := err.(type) {
// 		case conf.ConfigMarshalError:
// 			// this is the error we expected so ok.
// 			//t.Fatalf("Error expected error was not of type ConfigMarshalError. Got %#v\n", err)
// 		case conf.ConfigReadError:
// 			// this is the error we expected so ok.
// 			//t.Fatalf("Error expected error was not of type ConfigReadError. Got %#v\n", err)
// 		case *os.PathError:
// 			// this is an error we expect so ok.
// 		default:
// 			// or we got an error we did not expect
// 			t.Fatalf("Error expected error was not of type %#v. Got %#v\n", errtype, err)
// 		}
// 	}
// 	// else we got the error we where expecting
// }

// func TestSetConfigFilename(t *testing.T) {
// 	var filename = "this/file/does/not/exist.toml"
// 	err := conf.SetConfigFile(filename)
// 	if err == nil {
// 		t.Fatalf("Managed to find a file that does not exist %v!\n", filename)
// 	}
// 	_, ok := err.(*os.PathError)
// 	if !ok {
// 		t.Fatalf("Error expected error was not of type PathError. Got %#v\n", err)
// 	}
// 	// now try with a real file - via the ENV
// 	filename = os.Getenv("TEST_CONFIG_FILE")
// 	if filename == "" {
// 		t.Fatalf("Required enviromental variable \"TEST_CONFIG_FILE\" not set.\nIt should be the absolute path of the config file.")
// 	}

// 	err = conf.SetConfigFile(filename)
// 	if err != nil {
// 		t.Fatalf("Could not find the config file. Error %v\n", err)
// 	}
// }

// func TestVerifyConfig(t *testing.T) {
// 	//var c conf.Config
// 	var filename string
// 	// filename = os.Getenv("TEST_CONFIG_FILE")
// 	// if filename == "" {
// 	// 	t.Fatalf("Required enviromental variable \"TEST_CONFIG_FILE\" not set.\nIt should be the absolute path of the config file.")
// 	// }
// 	//	err := conf.ReadConfig(filename, &c)

// 	dir, ok := os.LookupEnv("TEST_CONFIG_DIR")
// 	if ok {
// 		viper.AddConfigPath(dir)
// 	}
// 	filename = os.Getenv("TEST_CONFIG_FILE")
// 	// if filename is unset then we'll use the default so this is safe

// 	c, err := conf.ReadConfigNew(filename)

// 	if err != nil {
// 		t.Fatalf("Could not read the config file. Error %v\n", err)
// 	}
// 	// verify
// 	if c.Server.Host != "localhost" {
// 		t.Fatalf("hostname not set. Expected: localhost Got: %v\n", c.Server.Host)
// 	}
// 	if c.Templates.Dir != "/home/owen/go/src/github.com/owenwaller/emailformgateway" {
// 		t.Fatalf("teampltes dir not set. Expected: /home/owen/go/src/github.com/owenwaller/emailformgateway Got: %v\n", c.Templates.Dir)
// 	}

// }

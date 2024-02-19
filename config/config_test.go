package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestMarshal(t *testing.T) {
	c, err := ReadConfig("config") // read the config.toml file in the package
	if err != nil {
		t.Fatalf("Failed to read config. Error: \"%s\"\n", err)
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		t.Fatalf("Failed to marshal config. Error: \"%s\"\n", err)
	}
}

func TestFieldAccess(t *testing.T) {
	var c Config
	viper.SetConfigName("config")
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not find cwd. Error \"%s\"\n", err)
	}
	viper.AddConfigPath(cwd)
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config. Error: \"%s\"\n", err)
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		t.Fatalf("Failed to marshal config. Error: \"%s\"\n", err)
	}

	// now print the fields
	var m map[string]string
	m = make(map[string]string)
	for _, v := range c.Fields {
		//fmt.Printf("Key:%v=%v\nKey[%v]\n", k, v, v.Name)
		m[v.Name] = v.Name
	}
	//fmt.Printf("Map=%v\n", m)
}

func TestReadConfig(t *testing.T) {
	configFileName := DefaultConfigFilename
	c, err := ReadConfig(configFileName)
	if err != nil {
		t.Fatal(err)
	}

	ec := newDefaultTestConfig()
	if err := verifyConfigs(c, ec); err != nil {
		t.Fatal(err)
	}
}

func TestReadConfigWithEnvVarOverload(t *testing.T) {
	// As an example set the Auth>Password from an Env Var - this will take precedence over any config file value.
	err := viper.BindEnv("Auth.Password", "TEST_PASSWORD")
	if err != nil {
		t.Fatalf("Can't bind to TEST_PASSWORD")
	}
	configFileName := DefaultConfigFilename
	// this sill replace any config file default for Auth.Password with the Env Var value
	c, err := ReadConfig(configFileName)
	if err != nil {
		t.Fatal(err)
	}
	ec := newDefaultTestConfig()
	// set the Auth.Password to the expected value in the Env Var.
	envVarPassword := os.Getenv("TEST_PASSWORD")
	if envVarPassword != "" {
		// there is a value set, so we need to use that  value as the expected value in the test.
		ec.Auth.Password = envVarPassword
		// otherwise we use the value from the config
	}

	// compare the configs.
	if err := verifyConfigs(c, ec); err != nil {
		t.Fatal(err)
	}
}

func newDefaultTestConfig() *Config {
	ec := new(Config)

	ec.LogFile.Filename = "access.log"
	ec.LogFile.Path = "/var/log/emailformgateway"
	ec.LogFile.Level = "INFO"

	ec.Smtp.Host = "smtp.localhost"
	ec.Smtp.Port = 25

	ec.Auth.Username = "local.user@localhost"
	ec.Auth.Password = "password123"

	ec.Addresses.CustomerFrom = "do-not-reply@localhost"
	ec.Addresses.CustomerFromName = "Localhost Contact Us"
	ec.Addresses.CustomerReplyTo = "do-not-reply@localhost"
	ec.Addresses.SystemTo = "to@localhost"
	ec.Addresses.SystemToName = "Localhost Contact Us Form"
	ec.Addresses.SystemFrom = "do-not-reply@localhost"
	ec.Addresses.SystemFromName = "Localhost Contact Us Form"
	ec.Addresses.SystemReplyTo = "do-not-reply@localhost.com"

	ec.Subjects.Customer = "Thank you for contacting localhost!"
	ec.Subjects.System = "Localhost Contact Us Form Message:"

	ec.Templates.Dir = "/template/dir"
	ec.Templates.CustomerText = "customer-email-text.template"
	ec.Templates.CustomerHtml = "customer-email-html.template"
	ec.Templates.SystemText = "system-email-text.template"
	ec.Templates.SystemHtml = "system-email-html.template"

	ec.Fields = make(map[string]FieldData)

	ec.Fields["field1"] = FieldData{Name: "name", Type: "textRestricted"}
	ec.Fields["field2"] = FieldData{Name: "email", Type: "email"}
	ec.Fields["field3"] = FieldData{Name: "subject", Type: "textRestricted"}
	ec.Fields["field4"] = FieldData{Name: "feedback", Type: "textUnrestricted"}

	return ec
}

func verifyConfigs(c, ec *Config) error {
	if c.LogFile != ec.LogFile {
		return fmt.Errorf("Logfile\nGot\n%+v\nExpected\n%+v\n", c.LogFile, ec.LogFile)
	}
	if c.Smtp != ec.Smtp {
		return fmt.Errorf("Smtp\nGot\n%+v\nExpected\n%+v\n", c.Smtp, ec.Smtp)
	}
	if c.Auth != ec.Auth {
		return fmt.Errorf("Auth\nGot\n%+v\nExpected\n%+v\n", c.Auth, ec.Auth)
	}
	if c.Addresses != ec.Addresses {
		return fmt.Errorf("Addresses\nGot\n%+v\nExpected\n%+v\n", c.Addresses, ec.Addresses)
	}
	if c.Subjects != ec.Subjects {
		return fmt.Errorf("Subjects\nGot\n%+v\nExpected\n%+v\n", c.Subjects, ec.Subjects)
	}
	if c.Templates != ec.Templates {
		return fmt.Errorf("Templates\nGot\n%+v\nExpected\n%+v\n", c.Templates, ec.Templates)
	}
	for k, f := range c.Fields {
		value, found := ec.Fields[k]
		if !found {
			return fmt.Errorf("Fields[%s]: Can't find key: %q", k, k)
		}
		if value != f {
			return fmt.Errorf("Fields\nGot: Fields[%s]=%+v\nExpected: Fields[%s]=%+v. ", k, value, k, f)
		}
	}
	return nil
}

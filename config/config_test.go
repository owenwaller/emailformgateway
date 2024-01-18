package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestMarshal(t *testing.T) {
	//viper.SetConfigName("config")
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	t.Fatalf("Could not find cwd. Error \"%s\"\n", err)
	// }
	//viper.AddConfigPath(cwd)
	c := new(Config)
	err := ReadConfig("config.toml", c)
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

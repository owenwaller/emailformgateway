package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestMarshal(t *testing.T) {
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
	err = viper.Marshal(&c)
	if err != nil {
		t.Fatalf("Failed to marshal config. Error: \"%s\"\n", err)
	}
	viper.Debug()
}
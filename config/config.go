package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerData
	LogFile   LogFileData
	Smtp      SmtpData
	Addresses EmailAddressData
	Subjects  EmailSubjectData
	Templates EmailTemplatesData
	Fields    map[string]FieldData
}

type ServerData struct {
	Host string
	Path string
	Port int
}

type LogFileData struct {
	Filename string
	Path     string
	Level    string
}

type SmtpData struct {
	Host     string
	Port     int
	Username string
	Password string
}

type EmailAddressData struct {
	CustomerFrom     string
	CustomerFromName string
	CustomerReplyTo  string
	SystemTo         string
	SystemToName     string
	SystemFrom       string
	SystemFromName   string
	SystemReplyTo    string
}

type EmailSubjectData struct {
	Customer string
	System   string
}

type EmailTemplatesData struct {
	Dir                  string
	CustomerText         string
	CustomerHtml         string
	SystemText           string
	SystemHtml           string
	CustomerTextFileName string
	CustomerHtmlFileName string
	SystemTextFileName   string
	SystemHtmlFileName   string
}

type FieldData struct {
	Name string
	Type string
}

type EmailTemplateData struct {
	FormData      map[string]string
	UserAgent     string
	RemoteIp      string
	XForwardedFor string
}

const (
	DefaultConfigFilename = "config"
)

var c Config

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

func GetConfig() *Config {
	return &c
}

func SetConfigFile(configFilename string) error {
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

func ReadConfig(filename string, c *Config) error {
	err := SetConfigFile(filename)
	if err != nil {
		return err // returns an os.PathError
	}
	err = viper.ReadInConfig()
	if err != nil {
		re := NewConfigReadError("Could not read in config.", err.Error())
		return re
	}

	//fmt.Printf("All Keys: \"%#v\"\n", viper.AllKeys())
	err = viper.Marshal(c)
	if err != nil {
		me := NewConfigMarshalError("Failed to marshal config.", err.Error())
		return me
	}
	return err // will be nil
}

func SetUpTemplates() {
	// first get all the templates
	c.Templates.CustomerTextFileName = BuildTemplateFilename(c.Templates.Dir, c.Templates.CustomerText)
	c.Templates.CustomerHtmlFileName = BuildTemplateFilename(c.Templates.Dir, c.Templates.CustomerHtml)
	c.Templates.SystemTextFileName = BuildTemplateFilename(c.Templates.Dir, c.Templates.SystemText)
	c.Templates.SystemHtmlFileName = BuildTemplateFilename(c.Templates.Dir, c.Templates.SystemHtml)
}

func BuildTemplateFilename(dir, filename string) string {
	return filepath.Join(dir, filename)
}

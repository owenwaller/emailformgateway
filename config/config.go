package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	LogFile   LogFileData
	Smtp      SmtpData
	Auth      AuthData
	Addresses EmailAddressData
	Subjects  EmailSubjectData
	Templates EmailTemplatesData
	Fields    map[string]FieldData
}

type LogFileData struct {
	Filename string
	Path     string
	Level    string
}

type SmtpData struct {
	Host string
	Port int
}

type AuthData struct {
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
	DefaultConfigType     = "toml"
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
	viper.AddConfigPath("/etc/emailformgateway/")
	viper.AddConfigPath(".")
	viper.SetConfigType(DefaultConfigType)

	if configFilename == "" {
		configFilename = DefaultConfigFilename
		// viper automatically searches the CWD - see findConfigFile
		// need to make the path work on Windows
		viper.SetConfigName(configFilename)
		return nil
	}
	// The intended us for this block is for testing, against a different or non-existent config file
	// The block is not a general solution to pulling the filename from an absolute or relative path
	// In particular it will not handle hidden files or directories correctly i.e. were the filename start with "."
	// or files that have multiple extensions e.g. .tar.gz
	path := filepath.Dir(configFilename)
	ext := filepath.Ext(configFilename)
	base := filepath.Base(configFilename) // will be filename.ext or filename or last dir name
	viper.AddConfigPath(path)
	viper.SetConfigType(ext)
	configFilename = base[:len(base)-len(ext)] // if this filename does not exit we will fail when we try to open it for reading,
	viper.SetConfigName(configFilename)        // just in case the slicing results in an empty string

	return nil
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

func ReadConfig(filename string) (*Config, error) {
	c := new(Config)
	SetConfigFile(filename)
	err := viper.ReadInConfig()
	if err != nil {
		re := NewConfigReadError("Could not read in config.", err.Error())
		return nil, re
	}
	err = viper.Unmarshal(c)
	if err != nil {
		me := NewConfigMarshalError("Failed to marshal config.", err.Error())
		return nil, me
	}
	return c, err // will be nil
}

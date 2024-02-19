// Copyright (c) 2024 Owen Waller. All rights reserved.
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"net"
	"net/http"
	"strings"

	"github.com/owenwaller/emailformgateway/config"
	"github.com/owenwaller/emailformgateway/emailer"
	"github.com/owenwaller/emailformgateway/validation"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Field struct {
	Name  string `json:"name"` // the `json:` maps the json field names to the struct names
	Value string `json:"value"`
}

type formResponse struct {
	Valid     bool
	BadFields []string
}

type Server struct {
	mu         sync.Mutex
	configName string
	mux        *http.ServeMux
	corsMux    http.Handler
	host       string
}

func NewServer(host, port, route string) *Server {
	s := new(Server)
	s.mux = http.NewServeMux()
	s.mux.HandleFunc(route, s.gatewayHandler)
	s.corsMux = cors.Default().Handler(s.mux)
	s.host = host + ":" + port
	return s
}

func (s *Server) setConfigName(configName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.configName = configName
}

func (s *Server) getConfigName() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.configName
}

func (s *Server) Start(configFileName string) {
	// check that the config file exists
	_, err := config.ReadConfig(configFileName) // if configFileName is an empty string the default name "config" will be used.
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not find a config file called %q: %s\n", configFileName, err)
		os.Exit(-1)
	}
	// set the filename in the Server - this will be referenced by multiple Go routines so must be synchronised
	s.setConfigName(configFileName)
	http.ListenAndServe(s.host, s.corsMux)
}

func (s *Server) gatewayHandler(w http.ResponseWriter, r *http.Request) {
	// The web form sends a JSON array of key value encoded pairs like this:
	// [
	// 	{
	// 		"name": "name",
	// 		"value": "Me"
	// 	},
	// 	{
	// 		"name": "email",
	// 		"value": "Me@example.com"
	// 	},
	// 	{
	// 		"name": "subject",
	// 		"value": "The subject"
	// 	},
	// 	{
	// 		"name": "feedback",
	// 		"value": "The feedback"
	// 	}
	// ]
	//
	// read the json and decode it
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Could not read the http request body: %s", err)
	}
	var fields []Field
	err = json.Unmarshal(body, &fields)
	if err != nil {
		fmt.Printf("Error could not decode JSON - \"%s\"\n", err)
	}

	// validate and write the http response.
	var fr formResponse
	scrubFields(fields, &fr)
	// The server always writes HTTP 200 OK back to the client along with the form response.
	// The form response always sets the formResponse.Valid field to true or false. The browser based client
	// then looks at the value of the valid field to determine if the form data was rejected or not.
	// This isn't very RESTful, but it is the way it works ATM
	writeResponse(w, &fr)

	c, err := config.ReadConfig(s.getConfigName())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("SystemTo: %q\n", viper.GetString("Addresses.SystemTo"))
	log.Printf("SystemToName: %q\n", viper.GetString("Addresses.SystemToName"))
	log.Printf("Templates Dir: %q\n", viper.GetString("Templates.Dir"))
	log.Printf("Config File: %q\n", viper.ConfigFileUsed())

	// set the full path to the templates - this should really be done dynamically by the emailer package...
	c.Templates.CustomerTextFileName = config.BuildTemplateFilename(c.Templates.Dir, c.Templates.CustomerText)
	c.Templates.CustomerHtmlFileName = config.BuildTemplateFilename(c.Templates.Dir, c.Templates.CustomerHtml)
	c.Templates.SystemTextFileName = config.BuildTemplateFilename(c.Templates.Dir, c.Templates.SystemText)
	c.Templates.SystemHtmlFileName = config.BuildTemplateFilename(c.Templates.Dir, c.Templates.SystemHtml)

	// build the EmailTemplateData that we pass to emailer.SendMail. This holds the info we want to add to the email messages.
	var etd config.EmailTemplateData
	etd.FormData = createFormDataMap(fields)
	var ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	var xForwardedFor = r.Header.Get("X-FORWARDED-FOR")
	var ua = r.UserAgent()
	etd.UserAgent = ua
	etd.RemoteIp = ip
	etd.XForwardedFor = xForwardedFor

	// try to send the email
	err = emailer.SendEmail(etd, c.Smtp, c.Auth, c.Addresses, c.Subjects, c.Templates)
	if err != nil {
		log.Fatalf("Failed to send email; %s", err)
	}
	//fmt.Printf("SENT!\n")
}

func scrubFields(fields []Field, fr *formResponse) {
	//fmt.Printf("formResponse.Valid=%v\n", fr.Valid)
	// look in the config to see what fields we should expect
	var c = config.GetConfig()
	for _, v := range c.Fields {
		// find the type of the fields in the fields map we were sent that has the same name
		match, err := find(v.Name, fields)
		if err != nil {
			//fmt.Printf("Error: %s\n", err)
		}
		// else we have a match
		//fmt.Printf("Match found: \"%#v\"\n", match)
		// now we have a match we need to validate it according to its type
		validateField(match, v.Type, fr)
	}
}

func find(name string, fields []Field) (*Field, error) {
	// the slice is small so just iterate....
	for i := 0; i < len(fields); i++ {
		if strings.EqualFold(fields[i].Name, name) {
			return &fields[i], nil
		}
	}
	// ok no match so ths is an error
	return nil, errors.New("Could not find a field named \"" + name + "\" in the json block.")
}

func validateField(match *Field, requiredType string, fr *formResponse) {
	match.Value = strings.TrimSpace(match.Value)
	match.Value = validation.RemoveEmailHeaders(match.Value)
	match.Value = validation.RemoveScriptTagsAndContents(match.Value)
	match.Value = validation.EscapeHTML(match.Value)
	valid := false
	requiredType = strings.ToLower(requiredType)
	switch requiredType {
	case "email":
		valid = validation.ValidateAsEmail(match.Value)
	case "textrestricted":
		valid = validation.ValidateAsRestrictedText(match.Value)
	case "textunrestricted":
		valid = validation.ValidateAsUnrestrictedText(match.Value)
	default:
		//fmt.Printf("Unknown imput type: \"%s\"\n", requiredType)
	}

	if !valid {
		fr.setBadFields(match)
	}
}

func writeResponse(w http.ResponseWriter, fr *formResponse) {
	body, err := fr.marshal()
	if err != nil {
		log.Printf("Error could not create JSON response \"%s\"\n", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		log.Printf("Error: Could not write response \"%s\"\n", err)
	}
	// wipe the bad
	fr.clearBadFields()
}

func createFormDataMap(formFields []Field) map[string]string {
	// now print the fields
	titler := cases.Title(language.English)
	var m map[string]string
	m = make(map[string]string)
	for _, v := range formFields {
		m[titler.String(v.Name)] = v.Value
	}
	return m
}

func (fr *formResponse) setBadFields(match *Field) {
	fr.BadFields = append(fr.BadFields, match.Name)
}

func (fr *formResponse) setValidity() {
	if len(fr.BadFields) == 0 { // have not yet appended so still good.
		fr.Valid = true
	} else {
		fr.Valid = false
	}
}

func (fr *formResponse) clearBadFields() {
	fr.Valid = false
	fr.BadFields = nil
}

func (fr *formResponse) marshal() (buf []byte, err error) {
	fr.setValidity()
	return json.Marshal(fr)
}

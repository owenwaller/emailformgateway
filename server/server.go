package server

import (
	"encoding/json"
	"errors"
	"io"

	//"fmt"

	"net"
	"net/http"
	"strconv"
	"strings"

	conf "github.com/owenwaller/emailformgateway/config"
	"github.com/owenwaller/emailformgateway/emailer"
	"github.com/owenwaller/emailformgateway/validation"
	"github.com/rs/cors"
)

type Field struct {
	Name  string `json:"name"` // the `json:` maps the json field names to the struct names
	Value string `json:"value"`
}

type formResponse struct {
	Valid     bool
	BadFields []string
}

func Start(c *conf.Config) {
	//fmt.Println("Creating new serve mux")
	mux := http.NewServeMux()
	mux.HandleFunc(c.Server.Path, gatewayHandler)

	corsMux := cors.Default().Handler(mux)
	//fmt.Printf("Registering %s => errorHandler\n", c.Server.Path)
	//fmt.Println("Listening on localhost:1314")
	host := c.Server.Host + ":" + strconv.Itoa(c.Server.Port)
	http.ListenAndServe(host, corsMux)
}

func gatewayHandler(w http.ResponseWriter, r *http.Request) {
	// the formResponse must be local - the handler runs in its own go routine
	// and we cannot share the form response across different requests.
	var fr formResponse
	// read the json and print it
	body, err := io.ReadAll(r.Body)

	var fields []Field
	err = json.Unmarshal(body, &fields)
	if err != nil {
		//fmt.Printf("Error could not decode JSON - \"%s\"\n", err)
	}
	//fmt.Printf("Decoded as \"%#v\"\n", fields)
	//fmt.Printf("fields[0].Value: %s\n", fields[0].Value)

	scrubFields(fields, &fr)
	writeResponse(w, &fr)

	//fmt.Printf("GetConfig()\n")
	var c = conf.GetConfig()
	var etd conf.EmailTemplateData
	//fmt.Printf("createFormDataMap()\n")
	etd.FormData = createFormDataMap(fields)
	//fmt.Printf("SplitHostPort()\n")
	var ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	//fmt.Printf("r.Header.Get()\n")
	var xForwardedFor = r.Header.Get("X-FORWARDED-FOR")
	//fmt.Printf("r.UserAgent()\n")
	var ua = r.UserAgent()
	etd.UserAgent = ua
	etd.RemoteIp = ip
	etd.XForwardedFor = xForwardedFor
	//fmt.Printf("SendEmail()\n")
	//fmt.Printf("%v\n", c)
	//fmt.Printf("%v\n", etd)
	// need to add in the user agent, remote IP and XForwardedFor IP
	emailer.SendEmail(etd, c.Smtp, c.Addresses, c.Subjects, c.Templates)
	//fmt.Printf("SENT!\n")
}

func scrubFields(fields []Field, fr *formResponse) {
	//fmt.Printf("formResponse.Valid=%v\n", fr.Valid)
	// look in the config to see what fields we should expect
	var c = conf.GetConfig()
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
		//fmt.Printf("Error could not create JSON respnse \"%s\"\n", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)

	if err != nil {
		//fmt.Printf("Error: Could not write response \"%s\"\n", err)
	}
	// wipe the bad
	fr.clearBadFields()
}

func createFormDataMap(formFields []Field) map[string]string {
	// now print the fields
	var m map[string]string
	m = make(map[string]string)
	for _, v := range formFields {
		m[strings.Title(v.Name)] = v.Value
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

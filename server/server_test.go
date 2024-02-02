package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFormResponse(t *testing.T) {
	var fr formResponse
	var f = Field{Name: "email", Value: "true"}

	fr.setValidity()

	var expected = true
	if fr.Valid != expected {
		t.Fatalf("A empty formResponse should be valid but is not.\n")
	}
	fr.setBadFields(&f)
	fr.setValidity()
	expected = false
	if fr.Valid != expected {
		t.Fatalf("A formResponse with one or more bad fields should not be valid, but is\n")
	}

	fr.clearBadFields()
	fr.setValidity()
	expected = true
	if fr.Valid != expected {
		t.Fatalf("After clearing badFields the response should be valid but is not\n")
	}
}

func TestFormResponseMarshal(t *testing.T) {
	var fr formResponse
	var f = Field{Name: "email", Value: "true"}
	var expected = "{\"Valid\":true,\"BadFields\":null}"
	var b []byte
	var result string

	b, err := fr.marshal()
	if err != nil {
		t.Fatalf("Could not marshal formRequest. Err: %v\n", err)
	}
	result = string(b)
	if result != expected {
		t.Fatalf("Did not marshal to JSON as expected. Expected %v but got %v\n", expected, result)
	}

	fr.setBadFields(&f)
	expected = "{\"Valid\":false,\"BadFields\":[\"email\"]}"
	b, err = fr.marshal()
	if err != nil {
		t.Fatalf("Could not marshal non-empty formRequest. Err: %v\n", err)
	}
	result = string(b)
	if result != expected {
		t.Fatalf("Did not marshal non-empty formRequest to JSON as expected. Expected %v but got %v\n", expected, result)
	}

	f.Name = "subject" // make the subejct bad as well as the email
	fr.setBadFields(&f)
	expected = "{\"Valid\":false,\"BadFields\":[\"email\",\"subject\"]}"
	b, err = fr.marshal()
	if err != nil {
		t.Fatalf("Could not marshal non-empty formRequest. Err: %v\n", err)
	}
	result = string(b)
	if result != expected {
		t.Fatalf("Did not marshal non-empty formRequest to JSON as expected. Expected %v but got %v\n", expected, result)
	}

}

func TestServerResponseBodyWithBadFields(t *testing.T) {
	var fr formResponse
	var f = Field{Name: "email", Value: "true"}
	fr.setBadFields(&f)

	w := httptest.NewRecorder()
	writeResponse(w, &fr)

	var expectedReturnCode = http.StatusOK
	var expectedBody = "{\"Valid\":false,\"BadFields\":[\"email\"]}"
	var responseCode = w.Code
	if expectedReturnCode != responseCode {
		t.Fatalf("Did not get the correct HTTP response code. Expected %v got %v\n", expectedReturnCode, responseCode)
	}

	var body = w.Body.String()
	if expectedBody != body {
		t.Fatalf("Did not get the correct HTTP response body. Expected %v got %v\n", expectedBody, body)
	}
}

func TestServerResponeBodyWithoutBadFields(t *testing.T) {
	var fr formResponse

	w := httptest.NewRecorder()
	writeResponse(w, &fr)

	var expectedReturnCode = http.StatusOK
	var expectedBody = "{\"Valid\":true,\"BadFields\":null}"
	var responseCode = w.Code
	if expectedReturnCode != responseCode {
		t.Fatalf("Did not get the correct HTTP response code. Expected %v got %v\n", expectedReturnCode, responseCode)
	}

	var body = w.Body.String()
	if expectedBody != body {
		t.Fatalf("Did not get the correct HTTP response body. Expected %v got %v\n", expectedBody, body)
	}

}

func TestServerResponseHeader(t *testing.T) {
	var fr formResponse
	var f = Field{Name: "email", Value: "true"}
	fr.setBadFields(&f)

	w := httptest.NewRecorder()
	writeResponse(w, &fr)

	var expectedKey = "Content-Type"
	var expectedValue = "application/json; charset=utf-8"
	// headeris a map[string][]string
	var header http.Header = w.Header()
	for k, v := range header {
		// check that the key exists in the expectedheader
		var keyExists = header.Get(expectedKey)
		if keyExists == "" {
			t.Fatalf("Expected a key called \"%v\", but did not find one.\n", expectedKey)
		}
		var i = 0
		var val string
		for i = 0; i < len(v); i++ {
			val = v[i]
			if val != expectedValue {
				t.Fatalf("Expecting the \"%v\" header to contain \"%v\" but found \"%v\"\n", k, expectedValue, val)
			}
		}
	}

}

// Test sending an email using an HTTP POST as the browser does.
func TestServerSendEmail(t *testing.T) {
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
	// These will then map into a []Fields
	fields := make([]Field, 0)

	// Create the slice of Field, pulling the To address from the env var
	fields = append(fields, Field{Name: "name", Value: "Me"})
	fields = append(fields, Field{Name: "subject", Value: "The subject"})
	fields = append(fields, Field{Name: "feedback", Value: "The feedback"})

	var to = os.Getenv("TEST_CUSTOMER_TO_EMAIL")
	if to == "" {
		t.Fatalf("Required environmental variable \"TEST_CUSTOMER_TO_EMAIL\" not set.")
	}
	fields = append(fields, Field{Name: "email", Value: to})

	// now encode the slice as JSON, as the Client side javascript does.
	b, err := json.Marshal(fields)
	jsonReader := bytes.NewReader(b)
	// create the server
	s := httptest.NewServer(http.HandlerFunc(gatewayHandler))
	defer s.Close()

	// use the default http client to POST to the server
	resp, err := http.Post(s.URL, "application/json; charset=utf-8", jsonReader)

	// we expect HTTP 200 OK back
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Got a status code of %d not %d. Error: %s", resp.StatusCode, http.StatusOK, err)
	}
	// We don't expect an error
	if err != nil {
		t.Fatalf("Failed to post to server: %s", err)
	}
	// we do expect an empty body on success
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to Read the response body: %s", err)
	}
	// the response should be a JSON encoded formResponse, with the Valid field set to true
	var fr formResponse
	err = json.Unmarshal(body, &fr)
	if err != nil {
		t.Fatalf("Could not decode form response JSON: %s", err)
	}
	// we expect Valid to be true and BadFields to be nil/empty
	if fr.Valid != true {
		t.Fatalf("Expected Valid to be %t but got %t.", true, fr.Valid)
	}
	if fr.BadFields != nil {
		t.Fatalf("Expected BadFields to be %v but got %v.", nil, fr.BadFields)
	}
}

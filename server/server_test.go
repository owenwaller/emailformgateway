package server

import (
	"net/http"
	"net/http/httptest"
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

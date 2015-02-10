package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/owenwaller/cors"
)

type formResponse struct {
	Valid     bool
	BadFields []string
}

func main() {
	fmt.Println("Creating new serve mux")
	corsMux := cors.NewServeMux()
	fmt.Println("Registering /success => successHandler")
	corsMux.HandleFunc("/success", successHandler)
	fmt.Println("Registering /error => errorHandler")
	corsMux.HandleFunc("/error", errorHandler)
	fmt.Println("Listening on localhost:1314")
	http.ListenAndServe("localhost:1314", corsMux)
}

func successHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered success handler")
	if !checkRequestTypeIsPost(w, r) {
		return
	}
	body := generateSuccessBody(w)
	returnResponse(w, body, http.StatusAccepted)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered error handler - banana")
	if !checkRequestTypeIsPost(w, r) {
		return
	}
	body := generateErrorBody()
	returnResponse(w, body, http.StatusOK)
}

func checkRequestTypeIsPost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed.Only POST requests are accepted", http.StatusMethodNotAllowed)
		fmt.Printf("Error: method was: %s\n", r.Method)
		return false
	}
	return true
}

func generateSuccessBody(w http.ResponseWriter) []byte {
	body := formResponse{Valid: true, BadFields: []string{""}}
	fmt.Printf("Success: body: \"%#v\"\n", body)
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println("error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Printf("Success: Responding with \"%#v\"\n", b)
	return b
}

func generateErrorBody() []byte {
	body := formResponse{Valid: false, BadFields: []string{"subject", "email"}}
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("Error: Responding with \"%s\"", string(b))
	return b
}

func returnResponse(w http.ResponseWriter, body []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err := w.Write(body)

	if err != nil {
		fmt.Println(err)
	}
}

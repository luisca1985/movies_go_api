package main

import (
	"fmt"
	"net/http"
)

func HandleRoot(w http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func HandleHome(w http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(w, "This is the API Endpoint")
}

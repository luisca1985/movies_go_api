package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	initTime := time.Now()
	servers := []string{
		"http://google.com",
		"http://facebook.com",
		"http://instagram.com",
	}

	for _, server := range servers {
		checkServer(server)
	}
	elapsedTime := time.Since(initTime)
	fmt.Printf("Execution time %s", elapsedTime)
}

func checkServer(server string) {
	_, err := http.Get(server)
	if err != nil {
		fmt.Println(server, "does not available.")
	} else {
		fmt.Println(server, "is working normally.")
	}
}

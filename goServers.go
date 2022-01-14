package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	initTime := time.Now()
	channel := make(chan string)
	servers := []string{
		"http://google.com",
		"http://facebook.com",
		"http://instagram.com",
	}

	for _, server := range servers {
		go checkServer(server, channel)
	}

	for i := 0; i < len(servers); i++ {
		fmt.Println(<-channel)
	}

	elapsedTime := time.Since(initTime)
	fmt.Printf("Execution time %s", elapsedTime)
}

func checkServer(server string, channel chan string) {
	_, err := http.Get(server)
	if err != nil {
		fmt.Println(server, "is not available.")
		channel <- server + " is not available."
	} else {
		fmt.Println(server, "is working normally.")
		channel <- server + " is working normally."
	}
}

package main

import (
	"fmt"
	"io"
	"net/http"
)

type webWritter struct {
}

func (webWritter) Write(p []byte) (int, error) {
	fmt.Println(string(p))
	return len(p), nil
}

func main() {
	response, err := http.Get("http://google.com")
	if err != nil {
		fmt.Println(err)
	}
	e := webWritter{}
	io.Copy(e, response.Body)
}

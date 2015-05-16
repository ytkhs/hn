package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func getJSON(url string) io.ReadCloser {
	req, _ := http.NewRequest("GET", url, nil)
	req.Close = true // connection reset

	client := new(http.Client)
	response, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return response.Body
}

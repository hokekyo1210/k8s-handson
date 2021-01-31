package main

import (
	"log"
	"net/http"

	"git.dmm.com/tsuchida-yuki1/cloud-functions-go/function"
)

func main() {
	// http.HandleFunc("/time", function.TimeUTC)
	http.HandleFunc("/time", function.TimeJST)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

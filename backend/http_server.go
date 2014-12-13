package main

import (
	"net/http"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
}

func initHttpServer() {
	http.HandleFunc("/bulk", httpHandler)
	http.ListenAndServe(httpDsn, nil)
}

package main

import (
	"net/http"
)

func Server() {
	// mux := http.NewServeMux()
	// mux.HandleFunc("POST /{$}", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprint(w, "Handler for POST request")
	// })

	mux := Handlers()

	http.ListenAndServe("localhost:8080", mux)
}

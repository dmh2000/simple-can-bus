package main

import (
	"net/http"
	"text/template"
)

func main() {
	var command = Command{}
	var telemetry = Telemetry{}
	var index = template.Must(template.ParseFiles("index.html"))

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			index.Execute(w, nil)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/command", &command)
	mux.Handle("/telemetry", &telemetry)

	print("Server is running on http://127.0.0.1:8000\n")
	http.ListenAndServe("127.0.0.1:8000", mux)
}

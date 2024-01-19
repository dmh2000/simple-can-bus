package main

import (
	"net/http"
	"text/template"
	"time"
)

func main() {
	var command = Command{}
	var telemetry = Telemetry{}
	var index = template.Must(template.ParseFiles("index.html"))

	// init the command data
	command.DioSet = 0
	command.DacSet = 0
	for i := 0; i < 10; i++ {
		err := command.Init()
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	print(command.DioSet)
	print(command.DacSet)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			index.Execute(w, command)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// handle command updates
	mux.Handle("/command", &command)

	// handle telemetry updates
	mux.Handle("/telemetry", &telemetry)

	print("Server is running on http://127.0.0.1:8000\n")
	err := http.ListenAndServe("127.0.0.1:8000", mux)
	if err != nil {
		panic(err)
	}
}

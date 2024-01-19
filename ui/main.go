package main

import (
	"net/http"
	"text/template"
)

func main() {
	var index = template.Must(template.ParseFiles("index.html"))

	var state = DeviceState{}
	err := state.fetch()
	if err != nil {
		print()
	}
	state.Tel.DioSet = state.Tel.DioOut
	state.Tel.DacSet = state.Tel.AdcOut
	state.Cmd.DioSet = state.Tel.DioOut
	state.Cmd.DacSet = state.Tel.AdcOut

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			index.Execute(w, state)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// handle parameter updates
	mux.Handle("/update", &state)

	print("Server is running on http://127.0.0.1:8000\n")
	err = http.ListenAndServe("127.0.0.1:8000", mux)
	if err != nil {
		panic(err)
	}
}

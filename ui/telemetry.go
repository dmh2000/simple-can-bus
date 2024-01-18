package main

import (
	"encoding/json"
	"io"
	"net/http"
	"text/template"
)

type Telemetry struct {
	DioSet int32
	DioOut int32
	DacSet int32
	AdcOut int32
}

var tmplTelementry = template.Must(template.ParseFiles("telemetry.html"))

func (telemetry *Telemetry) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// fetch the API data
		response, err := http.Get("http://localhost:6001/can/3")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// read the response body
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// unmarshal the response body
		err = json.Unmarshal(responseData, &telemetry)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the fragment
		tmplTelementry.Execute(w, telemetry)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

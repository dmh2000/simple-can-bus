package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"
)

// IMPORTANT! : The source code comments in this file were
// created using Github Copilot with the chat prompt
// "/explain annotate source"
// The generated comments are a bit wordy. Is this level of detail
// really necessary? I don't know. I'm just trying it out.

type Command struct {
	DioSet int32
	DacSet int32
}

type Telemetry struct {
	DioSet int32
	DioOut int32
	DacSet int32
	AdcOut int32
}

type DeviceState struct {
	Cmd Command
	Tel Telemetry
}

var tmplCommand = template.Must(template.ParseFiles("command.html"))
var tmplTelementry = template.Must(template.ParseFiles("telemetry.html"))

func (state *DeviceState) fetch() error {
	// fetch the API data
	response, err := http.Get("http://localhost:6001/can/3")
	if err != nil {
		return err
	}

	// read the response body
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// unmarshal the response body
	err = json.Unmarshal(responseData, &state.Tel)
	if err != nil {
		return err
	}

	state.Cmd.DioSet = state.Tel.DioSet
	state.Cmd.DacSet = state.Tel.DacSet

	return nil
}

// The put function sends a PUT request to the specified URL with the given payload.
func put(url string, payload string) error {
	// Create a new HTTP request with the PUT method, the specified URL, and the payload as the request body.
	// The payload is converted to a byte slice and wrapped in a bytes.Buffer to create an io.Reader, which is what http.NewRequest expects for the request body.
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(payload)))
	// If an error occurs while creating the request, panic.
	if err != nil {
		return err
	}

	// Set the Content-Type header of the request to "application/json".
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client.
	client := http.Client{}

	// Use the client to send the HTTP request.
	res, err := client.Do(req)

	// If an error occurs while sending the request, return the error.
	if err != nil {
		return err
	}

	// Ensure that the response body is closed when the function returns.
	defer res.Body.Close()

	// If no errors occurred, return nil.
	return nil
}

// ServeHTTP makes the Command type an http.Handler.
func (state *DeviceState) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Switch on the HTTP method of the request.
	switch r.Method {
	case "GET":

		err := state.fetch()
		if err != nil {
			// no connection to the API
			state.Tel.DioSet = -1
			state.Tel.DioOut = -1
			state.Tel.DacSet = -1
			state.Tel.AdcOut = -1
		}

		// update the fragments
		tmplTelementry.Execute(w, state)
	case "PUT":
		// Parse the form data from the request.
		err := r.ParseForm()
		// If an error occurs, respond with a 400 Bad Request status code and return.
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Get the "DioSet" and "DacSet" values from the form data and convert them to integers.
		dio, err := strconv.Atoi(r.FormValue("DioSet"))
		// If an error occurs, respond with a 400 Bad Request status code and return.
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dac, err := strconv.Atoi(r.FormValue("DacSet"))
		// If an error occurs, respond with a 400 Bad Request status code and return.
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Set the DioSet and DacSet fields of the Command to the parsed values.
		state.Cmd.DioSet = int32(dio)
		state.Cmd.DacSet = int32(dac)

		// Create a JSON string with the DioSet value.
		s := fmt.Sprintf("{\"Dio\": \"%d\"}", state.Cmd.DioSet)
		// Send a PUT request to the "/can/1" endpoint with the JSON string as the payload.
		err = put("http://localhost:6001/can/1", s)
		// If an error occurs, respond with a 500 Internal Server Error status code and return.
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Create a JSON string with the DacSet value.
		s = fmt.Sprintf("{\"Dac\": \"%d\"}", state.Cmd.DacSet)
		// Send a PUT request to the "/can/2" endpoint with the JSON string as the payload.
		err = put("http://localhost:6001/can/2", s)
		// If an error occurs, respond with a 500 Internal Server Error status code and return.
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the fragments
		tmplCommand.Execute(w, state)
	default:
		// If the HTTP method is not PUT, respond with a 405 Method Not Allowed status code.
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

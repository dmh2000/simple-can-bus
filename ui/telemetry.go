package main

// var tmplTelementry = template.Must(template.ParseFiles("telemetry.html"))

// func (state DeviceState) fetch() error {
// 	// fetch the API data
// 	response, err := http.Get("http://localhost:6001/can/3")
// 	if err != nil {
// 		return err
// 	}

// 	// read the response body
// 	responseData, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		return err
// 	}

// 	// unmarshal the response body
// 	err = json.Unmarshal(responseData, &state.Telemetry)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (state DeviceState) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":

// 		err := state.fetch()
// 		if err != nil {
// 			// no connection to the API
// 			state.Telemetry.DioSet = -1
// 			state.Telemetry.DioOut = -1
// 			state.Telemetry.DacSet = -1
// 			state.Telemetry.AdcOut = -1
// 		}

// 		// update the fragment
// 		tmplTelementry.Execute(w, state)
// 	default:
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 	}
// }

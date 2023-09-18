package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type CanDevice struct {
	DioIn  uint16
	DioOut uint16
	DacIn  int32
	AdcOut int32
}

var device = CanDevice{DioIn: 0, DioOut: 0, DacIn: 0, AdcOut: 0}

/*
	URLS
	CLIENT_DIO_IN  = /can/1
	CLIENT_DAC_IN  = /can/2
	CLIENT_DEVICE_OUT = /can/3
*/

// set the dio input commands
func DioIn(w http.ResponseWriter, r *http.Request) {
	var d CanDevice
	reqBody, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, &d)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("DioIn error: %s\n", reqBody)
		return
	}
	device.DioIn = d.DioIn
	device.DioOut = d.DioIn % 0x000f
	w.Header().Set("status", "200")
	fmt.Println(device)
}

func DacIn(w http.ResponseWriter, r *http.Request) {
	var d CanDevice
	reqBody, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, &d)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("DacIn error: %s\n", reqBody)
		return
	}
	device.DacIn = d.DacIn
	device.AdcOut = d.DacIn / 2
	w.Header().Set("status", "200")
	fmt.Println(device)
}

// get the current dio outputs
func DeviceOut(w http.ResponseWriter, r *http.Request) {
	jsondata, err := json.Marshal(device)
	fmt.Printf("jsondata: %s\n", jsondata)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("status", "200")
	fmt.Fprint(w, device)
}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/can/1", DioIn).Methods("PUT")
	r.HandleFunc("/can/2", DacIn).Methods("PUT")
	r.HandleFunc("/can/3", DeviceOut).Methods("GET")

	log.Fatal(http.ListenAndServe("127.0.0.1:6001", r))
}

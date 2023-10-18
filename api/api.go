package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CanDevice struct {
	DioIn  uint16
	DioOut uint16
	DacIn  int32
	AdcOut int32
}

type Dio struct {
	dio uint16
}

type Dac struct {
	dac int32
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
	setCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	setHeaders(&w)

	reqBody, _ := io.ReadAll(r.Body)
	v, err := unmarshalDio(string(reqBody))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("DacIn error: %s\n", reqBody)
		return
	}
	device.DioIn = v
	device.DioOut = device.DioIn % 0x55
}

func DacIn(w http.ResponseWriter, r *http.Request) {
	setCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	setHeaders(&w)

	reqBody, _ := io.ReadAll(r.Body)
	v, err := unmarshalDac(string(reqBody))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("DacIn error: %s\n", reqBody)
		return
	}
	device.DacIn = v
	device.AdcOut = device.DacIn / 2
}

// get the current dio outputs
func DeviceOut(w http.ResponseWriter, r *http.Request) {
	setCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	setHeaders(&w)
	jsondata, err := json.Marshal(device)
	if err != nil {
		fmt.Println(err)
		return
	}
	s := string(jsondata)
	fmt.Fprint(w, s)
}

func setCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func setHeaders(w *http.ResponseWriter) {
	// (*w).Header().Set("Access-Control-Allow-Origin", "*")
	// (*w).Header().Set("Access-Control-Allow-Methods", "*")
}

func main() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/can/1", DioIn).Methods(http.MethodPut, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	r.HandleFunc("/can/2", DacIn).Methods(http.MethodPut, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	r.HandleFunc("/can/3", DeviceOut).Methods(http.MethodGet, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	log.Fatal(http.ListenAndServe("127.0.0.1:6001", r))
}

// ======================================
// unmarshal utilities
// ======================================

type DioS struct {
	Dio string `json:"dio"`
}

type DacS struct {
	Dac string `json:"dac"`
}

func unmarshalDio(s string) (uint16, error) {
	var v DioS
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(v.Dio)
	if err != nil {
		return 0, err
	}

	return uint16(i), nil
}

func unmarshalDac(s string) (int32, error) {
	var v DacS
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(v.Dac)
	if err != nil {
		return 0, err
	}

	return int32(i), nil
}

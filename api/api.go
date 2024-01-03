package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"sqirvy.xyz/client"
	"sqirvy.xyz/types"
)

type CanDevice struct {
	DioSet uint16
	DioOut uint16
	DacSet int32
	AdcOut int32
}

type Dio struct {
	dio uint16
}

type Dac struct {
	dac int32
}

var device = CanDevice{DioSet: 0, DioOut: 0, DacSet: 0, AdcOut: 0}

/*
	URLS
	CLIENT_DIO_IN  = /can/1
	CLIENT_DAC_IN  = /can/2
	CLIENT_DEVICE_OUT = /can/3
*/

// set the dio input commands
func DioSet(w http.ResponseWriter, r *http.Request) {
	setCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	setHeaders(&w)

	reqBody, _ := io.ReadAll(r.Body)
	v, err := unmarshalDio(reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("DioSet error: %s\n", reqBody)
		return
	}
	device.DioSet = v
	fmt.Printf("Dio Set: %d\n", v)
	// forward to the CAN bus
	client.PutCanUint16(types.ID_DIO_SET, v)
}

func DacSet(w http.ResponseWriter, r *http.Request) {
	setCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	setHeaders(&w)

	reqBody, _ := io.ReadAll(r.Body)
	v, err := unmarshalDac(reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("DacSet error: %s\n", reqBody)
		return
	}
	device.DacSet = v
	fmt.Printf("Dac Set: %d\n", v)
	// forward to the CAN bus
	client.PutCanInt32(types.ID_DAC_SET, v)
}

// get the current dio outputs
func DeviceOut(w http.ResponseWriter, r *http.Request) {
	var err error
	var adc int32
	var dio uint16

	setCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	setHeaders(&w)

	// update the outputs
	adc, err = client.GetCanInt32(types.ID_ADC_OUT)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("AdcIn error: %d\n", err)
		return
	}
	device.AdcOut = adc

	dio, err = client.GetCanUint16(types.ID_DIO_OUT)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("AdcIn error: %d\n", err)
		return
	}
	device.DioOut = dio

	fmt.Printf("%d %d %d %d\n", device.DioSet, device.DioOut, device.DacSet, device.AdcOut)

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
	// start the client
	go client.Run()

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/can/1", DioSet).Methods(http.MethodPut, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	r.HandleFunc("/can/2", DacSet).Methods(http.MethodPut, http.MethodOptions)
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

func unmarshalDio(b []byte) (uint16, error) {
	var v DioS
	err := json.Unmarshal(b, &v)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(v.Dio)
	if err != nil {
		return 0, err
	}

	return uint16(i), nil
}

func unmarshalDac(b []byte) (int32, error) {
	var v DacS
	err := json.Unmarshal(b, &v)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(v.Dac)
	if err != nil {
		return 0, err
	}

	return int32(i), nil
}

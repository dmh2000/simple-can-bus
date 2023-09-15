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

/*
*

	CLIENT_DIO_IN  = 1
	CLIENT_DIO_OUT = 2
	CLIENT_DAC_IN  = 3
	CLIENT_ADC_OUT = 4
*/

// set the dio input commands
func DioIn(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	fmt.Printf("reqBody: %s\n", reqBody)

	//client.PutCanUint16(1,11)
	w.WriteHeader(http.StatusOK)
}

// get the current dio outputs
func DioOut(w http.ResponseWriter, r *http.Request) {
	jsondata, err := json.Marshal(CanDevice{DioIn: 0xffff, DioOut: 0, DacIn: 0, AdcOut: 0})
	fmt.Printf("jsondata: %s\n", jsondata)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprint(w, string(jsondata))
}

func DacIn(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	fmt.Printf("reqBody: %s\n", reqBody)
	fmt.Fprint(w, "DacIn")
}

// get the current dio outputs
func AdcOut(w http.ResponseWriter, r *http.Request) {
	jsondata, err := json.Marshal(CanDevice{DioIn: 0, DioOut: 0, DacIn: 0, AdcOut: 9876})
	fmt.Printf("jsondata: %s\n", jsondata)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprint(w, string(jsondata))

}

func main() {

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/can/1", DioIn).Methods("PUT")
	r.HandleFunc("/can/2", DioOut).Methods("GET")
	r.HandleFunc("/can/3", DacIn).Methods("PUT")
	r.HandleFunc("/can/4", AdcOut).Methods("GET")

	log.Fatal(http.ListenAndServe("127.0.0.1:6000", r))
}

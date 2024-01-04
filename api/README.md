# SQIRVY API

This api is a 'rest'-ish api that provides remote access to two simulations : a can bus simulator and a modbus simulator

## API

/can/send/1 : PUT for CAN ID 1
/can/send/2 : PUT for CAN ID 2
/can/recv/3 : GET for CAN ID 3
/can/recv/4 : GET for CAN ID 4

## ./api/api.go

### ./api/api.go::func main()

The provided Go code is the main function of a web server that provides an API for interacting with a CAN (Controller Area Network) device.

The client.Run() function is started as a goroutine, which means it will run concurrently with the rest of the main function. This function is responsible for managing the communication with the CAN device.

Next, a new router is created using the mux.NewRouter() function from the Gorilla Mux package, a powerful HTTP router and URL matcher for building Go web servers. The StrictSlash(true) option is set, which means the router will automatically redirect paths with trailing slashes to the same path without the trailing slash.

The router is then configured with several routes, each associated with a specific handler function and HTTP methods:

The /can/1 route is associated with the DioSet handler function and is configured to respond to HTTP PUT and OPTIONS methods. This route is used to set the DIO (Digital Input/Output) value of the CAN device.
The /can/2 route is associated with the DacSet handler function and is configured to respond to HTTP PUT and OPTIONS methods. This route is used to set the DAC (Digital-to-Analog Converter) value of the CAN device.
The /can/3 route is associated with the DeviceOutput handler function and is configured to respond to HTTP GET and OPTIONS methods. This route is used to retrieve the output value from the CAN device.
The mux.CORSMethodMiddleware(r) function is used to apply Cross-Origin Resource Sharing (CORS) middleware to the router. This allows the server to respond to requests from different origins, which is necessary if the API is to be accessed from a web browser on a different domain.

Finally, the server is started on localhost (127.0.0.1) port 6001 using the http.ListenAndServe function. If this function returns an error (for example, if the port is already in use), the program will log the error message and exit.

### ./api/api.go::func DioSet(w http.ResponseWriter, r \*http.Request)

The provided Go code is part of a web server that provides an API for interacting with a CAN (Controller Area Network) device. The DioSet function is a handler function that is called when a HTTP request is made to the /can/1 route.

The DioSet function first calls the setCors function to set the CORS (Cross-Origin Resource Sharing) headers in the HTTP response. This allows the server to respond to requests from different origins, which is necessary if the API is to be accessed from a web browser on a different domain.

If the HTTP method of the request is OPTIONS, the function immediately returns. This is because an OPTIONS request is a preflight request in CORS and is used by the browser to check if the actual request is safe to send.

Next, the setHeaders function is called to set additional headers in the HTTP response. However, the code inside this function is commented out, so it doesn't actually do anything at the moment.

The function then reads the body of the HTTP request and attempts to unmarshal it into a DioS struct using the unmarshalDio function. If an error occurs during unmarshalling, the function writes a 400 Bad Request status code to the HTTP response and prints an error message.

If the unmarshalling is successful, the function updates the DioSet field of the device global variable with the unmarshalled value and prints a message. It then calls the client.PutCanUint16 function to send a CAN frame with the types.ID_DIO_SET ID and the unmarshalled value to the CAN bus.

The unmarshalDio function takes a byte slice, attempts to unmarshal it into a DioS struct using the json.Unmarshal function, converts the Dio field of the struct to an integer using the strconv.Atoi function, and returns the integer as a uint16. If an error occurs during unmarshalling or conversion, the function returns the error.

### ./api/api.go::func DacSet(w http.ResponseWriter, r \*http.Request)

The provided Go code is part of a web server that provides an API for interacting with a CAN (Controller Area Network) device. The DacSet function is a handler function that is called when a HTTP request is made to a specific route (presumably /can/2 based on previous context).

The DacSet function first calls the setCors function to set the CORS (Cross-Origin Resource Sharing) headers in the HTTP response. This allows the server to respond to requests from different origins, which is necessary if the API is to be accessed from a web browser on a different domain.

If the HTTP method of the request is OPTIONS, the function immediately returns. This is because an OPTIONS request is a preflight request in CORS and is used by the browser to check if the actual request is safe to send.

Next, the setHeaders function is called to set additional headers in the HTTP response. However, the code inside this function is commented out, so it doesn't actually do anything at the moment.

The function then reads the body of the HTTP request and attempts to unmarshal it into an integer using the unmarshalDac function. If an error occurs during unmarshalling, the function writes a 400 Bad Request status code to the HTTP response and prints an error message.

If the unmarshalling is successful, the function updates the DacSet field of the device global variable with the unmarshalled value and prints a message. It then calls the client.PutCanInt32 function to send a CAN frame with the types.ID_DAC_SET ID and the unmarshalled value to the CAN bus.

The unmarshalDac function takes a byte slice, attempts to unmarshal it into a DacS struct using the json.Unmarshal function, converts the Dac field of the struct to an integer using the strconv.Atoi function, and returns the integer as an int32. If an error occurs during unmarshalling or conversion, the function returns the error.

### ./api/api.go::func DeviceOutput(w http.ResponseWriter, r \*http.Request)

The provided Go code is part of a web server that provides an API for interacting with a CAN (Controller Area Network) device. The DeviceOutput function is a handler function that is called when a HTTP request is made to a specific route (presumably /can/3 based on previous context).

The DeviceOutput function first calls the setCors function to set the CORS (Cross-Origin Resource Sharing) headers in the HTTP response. This allows the server to respond to requests from different origins, which is necessary if the API is to be accessed from a web browser on a different domain.

If the HTTP method of the request is OPTIONS, the function immediately returns. This is because an OPTIONS request is a preflight request in CORS and is used by the browser to check if the actual request is safe to send.

Next, the setHeaders function is called to set additional headers in the HTTP response. However, the code inside this function is commented out, so it doesn't actually do anything at the moment.

The function then calls the client.GetCanInt32 and client.GetCanUint16 functions to retrieve the ADC output and DIO output values from the CAN device, respectively. If an error occurs during these operations, the function writes a 500 Internal Server Error status code to the HTTP response and prints an error message.

If the operations are successful, the function updates the AdcOut and DioOut fields of the device global variable with the retrieved values and prints a message.

The function then marshals the device global variable into a JSON string using the json.Marshal function. If an error occurs during marshalling, the function prints the error message and returns.

Finally, the function writes the JSON string to the HTTP response. This allows the client to retrieve the current state of the CAN device.

### ./api/api.go::func unmarshalDio(b []byte) (uint16, error)

The provided Go code defines a function unmarshalDio that takes a byte slice as input and returns a uint16 and an error. This function is used to unmarshal a JSON object into a DioS struct and convert a string to a uint16.

The DioS struct is defined with a single field Dio of type string. The json:"dio" tag indicates that when this struct is unmarshalled from a JSON object, the Dio field should be populated with the value associated with the "dio" key in the JSON object.

Inside the unmarshalDio function, a DioS struct v is declared. The json.Unmarshal function is then called with the byte slice and a pointer to v. This function attempts to unmarshal the byte slice into v. If an error occurs during unmarshalling (for example, if the byte slice is not valid JSON), the function returns 0 and the error.

Next, the strconv.Atoi function is called with v.Dio. This function attempts to convert the Dio string to an integer. If an error occurs during conversion (for example, if the string contains non-numeric characters), the function returns 0 and the error.

If the unmarshalling and conversion are successful, the function returns the integer as a uint16 and nil for the error. This allows the caller to use the Dio value as a uint16.

### ./api/api.go::func unmarshalDac(b []byte) (int32, error)

The provided Go code defines a function unmarshalDac that takes a byte slice as input and returns an int32 and an error. This function is used to unmarshal a JSON object into a DacS struct and convert a string to an int32.

The DacS struct is defined with a single field Dac of type string. The json:"dac" tag indicates that when this struct is unmarshalled from a JSON object, the Dac field should be populated with the value associated with the "dac" key in the JSON object.

Inside the unmarshalDac function, a DacS struct v is declared. The json.Unmarshal function is then called with the byte slice and a pointer to v. This function attempts to unmarshal the byte slice into v. If an error occurs during unmarshalling (for example, if the byte slice is not valid JSON), the function returns 0 and the error.

Next, the strconv.Atoi function is called with v.Dac. This function attempts to convert the Dac string to an integer. If an error occurs during conversion (for example, if the string contains non-numeric characters), the function returns 0 and the error.

If the unmarshalling and conversion are successful, the function returns the integer as an int32 and nil for the error. This allows the caller to use the Dac value as an int32.

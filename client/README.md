# Can Client

This module implements the specific interface between a higher level client and the simulated CAN device. All communication takes place over the CAN bus.

## Devices

The interface includes the following devices:

- ADC : 32 bit analog to digital input
- DAC : 32 bit digital to analog output
- DIO OUT : 32 bit digital output
- DIO IN : 32 bit digital input

### Analog

The client sets the value of the DAC output on demand.
The device outputs the current value of the ADC at 10 hz.
The device DAC output value set by the client loops back to the ADC output value.

### Digital

The client sets the value of the DIO input on demand.
The device outputs the current value of the 32 bit DIO input.
The device DIO input value set by the client loops back to the DIO output value.

## Messages

| Message Id | Sender    | Receiver  | Description                 |
| ---------- | --------- | --------- | --------------------------- |
| 1          | simulator | any       | periodic value of DAC 10 Hz |
| 2          | client    | simulator | set DAC value               |
| 3          | simulator | any       | periodic value of DIO IN    |
| 4          | client    | simulator | set DIO OUT value           |

## ./client/client.go

### ./client/client.go::func Run()

The provided Go code is part of a client program that communicates with a device over a CAN (Controller Area Network) interface.

The CanData struct is used to represent the data of a CAN device. It has four fields, all of which use atomic types to ensure safe concurrent access:

adc_in: an atomic integer (32 bits) representing the ADC (Analog-to-Digital Converter) input value.
dio_in: an atomic unsigned integer (32 bits) representing the DIO (Digital Input/Output) input value.
dio_out: an atomic unsigned integer (32 bits) representing the DIO output value.
adc_out: an atomic integer (32 bits) representing the ADC output value.
The init method of the CanData struct initializes the state of the client. It sets the dio_in, dio_out, adc_in, and adc_out fields to 0.

The Run function is the main function of the client program. It initializes the CAN interface and the client state, creates channels for receiving CAN frames and signaling when to quit, and enters a loop to continuously receive CAN frames from the device.

Inside the loop, the can.CanRecv function is called to receive a CAN frame from the device. If an error or timeout occurs during the receive operation, an error message is printed and the loop continues with the next iteration.

If a CAN frame is successfully received, the switch statement checks the CAN ID of the frame. If the ID is types.ID_DIO_OUT, the DIO output value is updated in the client state. If the ID is types.ID_ADC_OUT, the ADC output value is updated in the client state. The loop then continues with the next iteration, waiting for the next CAN frame to be received.

### ./client/client.go::exported functions for API use

The provided Go code defines a set of functions that interact with a CAN (Controller Area Network) device. These functions allow you to send data to the device (PutCanUint16 and PutCanInt32) and retrieve data from the device (GetCanUint16 and GetCanInt32).

The PutCanUint16 function takes a CAN ID and a 16-bit unsigned integer value. It checks if the ID is types.ID_DIO_SET. If it is, the function updates the DIO set value in the CAN state and sends a CAN frame with the DIO set ID and value to the CAN bus. If the ID is not recognized, the function returns an error.

The PutCanInt32 function works similarly to PutCanUint16, but it takes a 32-bit integer value and checks for the types.ID_DAC_SET ID. If the ID matches, the function updates the DAC set value in the CAN state and sends a CAN frame with the DAC set ID and value to the CAN bus.

The GetCanUint16 function retrieves a 16-bit unsigned integer value from the CAN state. It checks if the requested ID is types.ID_DIO_OUT. If it is, the function retrieves the DIO output value from the CAN state. If the ID is not recognized, the function returns an error.

The GetCanInt32 function works similarly to GetCanUint16, but it retrieves a 32-bit integer value and checks for the types.ID_ADC_OUT ID. If the ID matches, the function retrieves the ADC output value from the CAN state.

The CanData struct methods (putDioSet, getDioOut, putDacSet, getAdcOut) are used to update and retrieve the DIO and ADC values in the CAN state. These methods use atomic operations to ensure safe concurrent access to the values.

# CAN DEVICE

This module implements the remote device that the client communicates with over the can bus.

## ./device/device.go

### ./device/device.go::data structures

The code defines two struct types: deviceState and CANFrame.

The deviceState struct is used to represent the state of a device. It has four fields: dio_in and dio_out of type uint16, and adc_in and adc_out of type int32. The dio_in and dio_out fields presumably represent digital input and output values, while the adc_in and adc_out fields presumably represent analog-to-digital converter input and output values.

The CANFrame struct is used to represent a CAN (Controller Area Network) frame. It also has four fields: adc_in, dio_in, dio_out, and adc_out, all of type atomic.Value. The atomic.Value type is used for storing and retrieving values in a way that is safe for concurrent use, which is important in a multi-threaded environment like a CAN network. The adc_in and dio_in fields presumably represent the values of the ADC and DIO inputs in the CAN frame, while the adc_out and dio_out fields presumably represent the values of the ADC and DIO outputs in the CAN frame.

### ./device/device.go::func main()

This code defines a main function that simulates a device communicating over a CAN (Controller Area Network) interface.

sockfd := can.CanInit("vcan0") initializes the CAN interface and returns a socket file descriptor. The "vcan0" argument is the name of the virtual CAN interface.

defer can.CanClose(sockfd) ensures that the CAN interface is closed when the function returns.

state := deviceState{0, 0, 0, 0} initializes the device state, which presumably includes four fields, all set to 0.

fch := make(chan can.CanFrame) and quit := make(chan bool) create channels for receiving CAN frames and signaling when to quit, respectively.

go recvDevice(sockfd, fch, quit) starts a goroutine to receive CAN frames from the device.

prev_dio_out := uint16(0) and prev_adc_out := int32(0) initialize the previous DIO (Digital Input/Output) and ADC (Analog-to-Digital Converter) outputs.

The for loop runs until q is true. Inside the loop, a select statement waits for either a CAN frame to be received or a timeout of 1 second.

If a CAN frame is received, the switch statement checks the CAN ID of the frame. If the ID is ID_DIO_SET, the DIO input is updated and printed. If the ID is ID_DAC_SET, the ADC input is updated and printed.

If a timeout occurs, the DIO and ADC outputs are updated to match the inputs. If the outputs have changed, they are printed and sent as CAN frames to the client.

The can.CanSend function sends a CAN frame to the CAN interface. The can.Uint16Frame and can.Int32Frame functions create CAN frames with the specified CAN ID and data.

### ./device/device.go::func recvDevice(sockfd int, fch chan<- can.CanFrame, quit <-chan bool)

The function recvDevice continuously receives CAN (Controller Area Network) frames from a device. The function takes three parameters: sockfd, which is an integer representing the socket file descriptor of the CAN interface; fch, which is a send-only channel for CAN frames; and quit, which is a receive-only channel for boolean values.

Inside the function, a can.CanFrame variable frame is declared to hold the received CAN frames. A for loop is then started, which will run indefinitely until explicitly broken.

Within the loop, the can.CanRecv function is called to receive a CAN frame from the bus. This function takes three parameters: the socket file descriptor, a pointer to the frame variable, and a timeout value. The return value of the function is stored in ret.

If ret is less than 0, this indicates an error occurred during the receive operation. In this case, an error message is printed and the loop continues with the next iteration. If ret is 0, this indicates a timeout occurred during the receive operation. In this case, the loop also continues with the next iteration.

Next, a select statement is used to check if a value has been sent on the quit channel. If a value has been sent and it is true, the break statement is executed, which will exit the for loop. If no value has been sent on the quit channel, the default case is executed, which does nothing.

Finally, the received CAN frame is sent on the fch channel. This allows the frame to be processed by another part of the program. The loop then continues with the next iteration, waiting for the next CAN frame to be received.

# A simple library for sending and receiving CAN message on Linux

Standard CAN bus on linux uses a network/socket paradigm.

canlib_framt_t is a simplified version of struct can_frame from /usr/include/linux/can.h. It is used here to keep it simple for the Go interface. If you only want to use the C functions you could substitute struct can_frame.

```c

// initialize the can bus and the client networks for sending
// @return a socket
int canlib_init(const char *can_dev);

// receive a single can frame message and populate the can_frame
// @return number of bytes received or < 0 if failed
int canlib_receive(int can_sock, canlib_frame_t *can_frame, int timeout_ms);

// send a single can frame message from the input can_frame
// @returns number of bytes sent or < 0 if failed
int canlib_send(int can_sock, canlib_frame_t *can_frame);

// close the can network interface
int canlib_close(int can_sock)

```

## ./c/canlib.c

### ./g/canlib.c::int canlib_init(const char \*can_dev)t

The function **canlib_init** initializes a CAN (Controller Area Network) interface. The function takes one parameter, can_dev, which is a string representing the name of the CAN device to be initialized.

The function starts by declaring two integer variables, sock and status. sock is used to store the socket descriptor returned by the socket function, and status is used to store the status of the bind function.

The socket function is called with PF_CAN as the domain (indicating a CAN protocol), SOCK_RAW as the type (indicating a raw socket), and CAN_RAW as the protocol (indicating a raw CAN protocol). The return value is stored in sock. If sock is less than 0, this indicates an error occurred when creating the socket, and the function returns -1.

Next, a struct ifreq named ifr is declared and its ifr_name field is populated with the can_dev string using the strncpy function. The ifr_name field is then null-terminated to ensure it is a valid C string. The ifr_ifindex field of ifr is populated with the index of the network interface corresponding to can_dev using the if_nametoindex function. If ifr_ifindex is 0, this indicates an error occurred when getting the interface index, and the function returns -2.

A struct sockaddr_can named addr is then declared and its fields are initialized. The can_family field is set to AF_CAN to indicate a CAN protocol, and the can_ifindex field is set to the interface index stored in ifr.ifr_ifindex.

The bind function is then called to bind the socket to the CAN interface. The bind function takes sock as the first argument, the address of addr cast to a struct sockaddr \* as the second argument, and the size of addr as the third argument. The return value is stored in status. If status is less than 0, this indicates an error occurred when binding the socket, and the function returns -3.

Finally, if no errors occurred, the function returns sock. This is likely a file descriptor that can be used to interact with the initialized CAN interface.

### ./c/canlib.c::int canlib_send(int can_sock, canlib_frame_t \*can_frame)

The function **canlib_send** sends a CAN (Controller Area Network) frame to a CAN interface. The function takes two parameters: can_sock, which is an integer representing the socket descriptor of the CAN interface, and can_frame, which is a pointer to a canlib_frame_t struct.

The canlib_frame_t struct presumably has at least three fields: can_id, can_dlc, and data. can_id is likely an identifier for the CAN frame. can_dlc is probably the data length code, which specifies the number of bytes in the data field of the frame. data is likely an array of bytes representing the data of the CAN frame.

In the canlib_send function, a local variable frame of type struct can_frame is declared. This struct is a standard struct used in Linux for representing CAN frames. Its fields are populated with the values from the can_frame parameter. The can_id and can_dlc fields are directly assigned the values of can_id and can_dlc from can_frame, respectively. The data field of frame is populated with the values from the data field of can_frame using the memcpy function.

The write function is then called to send the CAN frame to the CAN interface. The write function takes can_sock as the first argument, the address of frame as the second argument, and the size of frame as the third argument. The number of bytes written is returned by the write function and stored in the bytes variable.

If bytes is less than 0, this indicates an error occurred during the write operation, and the function returns bytes. Otherwise, the function returns bytes, which represents the number of bytes that were successfully written to the CAN interface.

### ./g/canlib.c::int canlib_receive(int can_sock, canlib_frame_t \*can_frame, int timeout_ms)

The function **canlib_receive** receives a CAN (Controller Area Network) frame from a CAN interface. The function takes three parameters: can_sock, which is an integer representing the socket descriptor of the CAN interface; can_frame, which is a pointer to a canlib_frame_t struct; and timeout_ms, which is an integer representing the timeout for the receive operation in milliseconds.

The function starts by declaring several local variables. frame is a struct can_frame that will store the received CAN frame. readfds is a fd_set that will be used with the select function to wait for the CAN interface to be ready for reading. timeout is a struct timeval that will be used to specify the timeout for the select function. status and bytes are integers that will store the return values of the select and read functions, respectively.

The FD_ZERO function is called to initialize readfds to an empty set. The timeout struct is then populated with the timeout value converted to seconds and microseconds.

The FD_SET function is called to add can_sock to readfds. The select function is then called to wait for can_sock to be ready for reading, with a timeout specified by timeout. If status is less than 0, this indicates an error occurred during the select operation, and the function returns status. If status is 0, this indicates a timeout occurred, and the function returns 0.

If status is not less than 0 or 0, the FD_ISSET function is called to check if can_sock is in readfds. If not, this indicates an error, and the function returns -1.

The read function is then called to read a CAN frame from can_sock into frame. If bytes is less than or equal to 0, this indicates an error occurred during the read operation, and the function returns bytes.

If bytes is greater than 0, the fields of can_frame are populated with the values from frame. The can_id and can_dlc fields are directly assigned the values of can_id and can_dlc from frame, respectively. The data field of can_frame is populated with the values from the data field of frame using the memcpy function.

Finally, the function returns bytes, which represents the number of bytes that were successfully read from the CAN interface.

### ./c/canlib.c::int canlib_close(int can_sock)

The function **canlib_close** closes a CAN (Controller Area Network) interface. The function takes one parameter: can_sock, which is an integer representing the socket descriptor of the CAN interface.

The function starts by declaring a local variable status and initializing it to 0. This variable is used to store the return value of the close function.

Next, there's an if statement that checks if can_sock is greater than 0. If it is, the close function is called with can_sock as an argument, and its return value is stored in status. The close function is a standard C function that closes a file descriptor, in this case, the socket descriptor for the CAN interface.

After the close function is called, there's another if statement that checks if status is less than 0. If it is, this indicates an error occurred when closing the socket, and the function returns errno. errno is a global variable in C that is set by system calls and some library functions in the event of an error to indicate what went wrong.

If status is not less than 0, this means the socket was successfully closed, and the function returns 0. This return value is likely used as an error code, with a value of 0 indicating success and other values indicating specific error conditions.

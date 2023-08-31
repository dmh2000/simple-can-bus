# A simple library for sending and receiving CAN message on Linux

Standard CAN bus on linux uses a network/socket paradigm.

canlib_framt_t is a simplified version of struct can_frame from /usr/include/linux/can.h. It is used here to keep it simple for the Go interface. If you only want to use the C functions you could substitute struct can_frame.

```c

// initialize the can bus and the client networks for sending
// @return a socket
int canlib_init(const char *can_dev);

// receive a single can frame message and populate the can_frame
// @return number of bytes received or < 0 if failed
int canlib_receive(int can_sock, canlib_frame_t *can_frame);

// send a single can frame message from the input can_frame
// @returns number of bytes sent or < 0 if failed
int canlib_send(int can_sock, canlib_frame_t *can_frame);

// close the can network interface
int canlib_close(int can_sock)

```

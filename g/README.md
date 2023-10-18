# A simple library for sending and receiving CAN message on Linux

Standard CAN bus on linux uses a network/socket paradigm.

The functions in this library rely on the C library in "../c".

```go
type CanFrame struct {
	CanId  uint32	// CAN message ID
	CanDlc byte		// payload length (0..8)
	Data   [8]byte  // data (valid up to data length)
}

// initialize a connection to a can device
// returns a linux socket
func CanInit(dev string) int

// send a single can message
// returns number of bytes sent or < 0 if failed
// CanSend always sends sizeof CanFrame = 16
// CanDlc specifies the number of valid data bytes
func CanSend(sock int, frame *CanFrame) int

// receive a single can message
// returns number of bytes received  or < 0 if failed
// CanReceive always receives sizeof CanFrame = 16,
// CanDlc specifies the number of valid data bytes
func CanRecv(sock C.int, frame *CanFrame) int

// close the can socket
func CanClose(sock int) int
```

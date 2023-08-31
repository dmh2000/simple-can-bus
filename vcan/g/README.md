# A simple library for sending and receiving CAN message on Linux

Standard CAN bus on linux uses a network/socket paradigm.

```go
type CanFrame struct {
	CanId  uint32
	CanDlc byte
	Data   [8]byte
}

// initialize a connection to a can device
// returns a linux socket
func CanInit(dev string) int

// send a single can message
// returns status < 0 = fail. >= 0 means success
func CanSend(sock int, frame *CanFrame) int

// receive a single can message
// returns status < 0 = fail. >= 0 means success
func CanRecv(sock C.int, frame *CanFrame) int

// close the can socket
func CanClose(sock int) int
```

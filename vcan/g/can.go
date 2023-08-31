package can

/*
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -L./ -lcan
#include "canlib.h"
struct canlib_frame;
*/
import "C"

// note : sizeof canlib_frame is 16 bytes due to padding alignment
type CanFrame struct {
	CanId  uint32
	CanDlc byte
	Data   [8]byte
}

func CanInit(dev string) int {
	var sock C.int
	sock = C.canlib_init(C.CString(dev))
	return int(sock)
}

func CanSend(sock int, frame *CanFrame) int {
	var ret C.int
	var cframe C.struct_canlib_frame
	cframe.can_id = C.uint(frame.CanId)
	cframe.can_dlc = C.uchar(frame.CanDlc)

	// always copy 8 bytes
	for i := 0; i < 8; i++ {
		cframe.data[i] = C.uchar(frame.Data[i])
	}

	ret = C.canlib_send(C.int(sock), &cframe)
	return int(ret)
}

// CanRecv receives a frame from the socket, blocking
func CanRecv(sock int, frame *CanFrame) int {
	var ret C.int
	var cframe C.struct_canlib_frame
	ret = C.canlib_receive(C.int(sock), &cframe)
	frame.CanId = uint32(cframe.can_id)
	frame.CanDlc = byte(cframe.can_dlc)

	// always copy 8 bytes
	for i := 0; i < 8; i++ {
		frame.Data[i] = byte(cframe.data[i])
	}

	return int(ret)
}

func CanClose(sock int) int {
	var ret = C.canlib_close(C.int(sock))
	return int(ret)
}

package can

import (
	"fmt"
	"syscall"
)

//#cgo CFLAGS: -g -Wall -I../c
//#cgo LDFLAGS: -L${SRCDIR}/../c -lcan
//#include "canlib.h"
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
func CanRecv(sock int, frame *CanFrame, timeout int) int {
	var ret C.int
	// var errno _Ctype_uint
	var cframe C.struct_canlib_frame
	ret = C.canlib_receive(C.int(sock), &cframe, C.int(timeout))
	frame.CanId = uint32(cframe.can_id)
	frame.CanDlc = byte(cframe.can_dlc)

	if ret < 0 {
		return int(ret)
	}

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

func CanPrint(frame *CanFrame) {
	fmt.Printf("%d\n", frame.CanId)
	fmt.Printf("%d\n", frame.CanDlc)
	for i := 0; i < 8; i++ {
		fmt.Printf("%02x", frame.Data[i])
	}
	fmt.Printf("\n")
}

func CanErrno() int {
	var errno = C.canlib_status()
	return int(errno)
}

func CanErrnoString() string {
	var errno = C.canlib_status()
	return syscall.Errno(errno).Error()
}

func Uint16Frame(id uint32, v uint16) CanFrame {
	frame := CanFrame{}
	frame.CanId = id
	frame.CanDlc = 2
	b := Uint16ToBytes(v)
	frame.Data[0] = b[0]
	frame.Data[1] = b[1]

	return frame
}

func Int32Frame(id uint32, v int32) CanFrame {
	frame := CanFrame{}
	frame.CanId = id
	frame.CanDlc = 4
	b := Int32ToBytes(v)
	frame.Data[0] = b[0]
	frame.Data[1] = b[1]
	frame.Data[2] = b[2]
	frame.Data[3] = b[3]

	return frame
}

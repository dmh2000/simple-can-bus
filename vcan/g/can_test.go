package can_test

import (
	"testing"

	. "sqirvy.xyz/can"
)

func TestCanInit(t *testing.T) {
	var sock int
	sock = CanInit("vcan0")
	if sock < 0 {
		t.Error("CanInit failed")
	}
}

func TestCanClose(t *testing.T) {
	var sock int
	sock = CanInit("vcan0")
	if sock < 0 {
		t.Error("CanInit failed")
	}
	ret := CanClose(sock)
	if ret < 0 {
		t.Error("CanClose failed")
	}
}

func TestCanSend(t *testing.T) {
	var sock int
	var ret int
	sock = CanInit("vcan0")
	if sock < 0 {
		t.Error("CanInit failed")
	}

	var frame CanFrame
	frame.CanId = 99
	frame.CanDlc = 8
	for i := 0; i < 8; i++ {
		frame.Data[i] = byte(i)
	}

	// run candump vcan0 in another terminal to see the frame
	ret = CanSend(sock, &frame)
	if ret < 0 {
		t.Error("CanSend failed")
	}

	// run candump vcan0 in another terminal to see the frame
	frame.CanDlc = 4
	ret = CanSend(sock, &frame)
	if ret < 0 {
		t.Error("CanSend failed")
	}

	ret = CanClose(sock)
	if ret < 0 {
		t.Error("CanClose failed")
	}
}

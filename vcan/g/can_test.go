package can_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
	"unsafe"

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

// TestCanSend sends frames to a receiver
// use candump vcan0 in another terminal to see the frames
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

	ret = CanSend(sock, &frame)
	if ret < 0 {
		t.Error("CanSend failed")
	}

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

func TestCanReceive1(t *testing.T) {

	var send = func(wg *sync.WaitGroup) {
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

		ret = CanSend(sock, &frame)
		if ret < 0 {
			t.Error("CanSend failed")
		}

		ret = CanClose(sock)

		wg.Done()
	}

	var recv = func(wg *sync.WaitGroup) {
		var sock int
		var ret int
		sock = CanInit("vcan0")
		if sock < 0 {
			t.Error("CanInit failed")
		}

		var frame CanFrame

		// run cansend vcan0 in another terminal to see the frame
		ret = CanRecv(sock, &frame)
		if ret < 0 {
			t.Error("CanSend failed")
		}

		if frame.CanId != 99 {
			t.Error("frame.CanId != 99")
		}

		if frame.CanDlc != 8 {
			t.Error("frame.CanDlc != 8")
		}

		if ret != int(unsafe.Sizeof(frame)) {
			t.Error(fmt.Sprintf("did not receive sizeof frame bytes"))
		}

		for i := 0; i < 8; i++ {
			if frame.Data[i] != byte(i) {
				t.Error(fmt.Sprintf("frame.Data[%d] != %d", i, i))
			}
		}

		CanClose(sock)

		wg.Done()
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go recv(&wg) // blocks
	time.Sleep(2 * time.Second)
	go send(&wg)
	wg.Wait()

}

func TestCanReceive2(t *testing.T) {

	var send = func(wg *sync.WaitGroup) {
		var sock int
		var ret int
		sock = CanInit("vcan0")
		if sock < 0 {
			t.Error("CanInit failed")
		}

		var frame CanFrame
		frame.CanId = 1
		frame.CanDlc = 4
		for i := 0; i < 4; i++ {
			frame.Data[i] = byte(i)
		}

		ret = CanSend(sock, &frame)
		if ret < 0 {
			t.Error("CanSend failed")
		}

		ret = CanClose(sock)

		wg.Done()
	}

	var recv = func(wg *sync.WaitGroup) {
		var sock int
		var ret int
		sock = CanInit("vcan0")
		if sock < 0 {
			t.Error("CanInit failed")
		}

		var frame CanFrame

		// run cansend vcan0 in another terminal to see the frame
		ret = CanRecv(sock, &frame)
		if ret < 0 {
			t.Error("CanSend failed")
		}

		if frame.CanId != 1 {
			t.Error("frame.CanId != 99")
		}

		if frame.CanDlc != 4 {
			t.Error("frame.CanDlc != 8")
		}

		if ret != int(unsafe.Sizeof(frame)) {
			t.Error(fmt.Sprintf("did not receive sizeof frame bytes"))
		}

		for i := 0; i < int(frame.CanDlc); i++ {
			if frame.Data[i] != byte(i) {
				t.Error(fmt.Sprintf("frame.Data[%d] != %d", i, i))
			}
		}

		CanClose(sock)

		wg.Done()
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go recv(&wg) // blocks
	time.Sleep(2 * time.Second)
	go send(&wg)
	wg.Wait()

}

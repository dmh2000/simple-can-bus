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

	ret, err := CanSend(sock, &frame)
	if ret < 0 || err != nil {
		t.Error("CanSend failed")
	}

	frame.CanDlc = 4
	ret, err = CanSend(sock, &frame)
	if ret < 0 || err != nil {
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

		ret, err := CanSend(sock, &frame)
		t.Log("sent frame")
		if ret < 0 || err != nil {
			t.Error("CanSend failed")
		}

		ret = CanClose(sock)

		wg.Done()
	}

	var recv = func(wg *sync.WaitGroup) {
		var sock int
		sock = CanInit("vcan0")
		if sock < 0 {
			t.Error("CanInit failed")
		}

		var frame CanFrame

		// run cansend vcan0 in another terminal to see the frame
		ret, err := CanRecv(sock, &frame, 10000)
		t.Log("received frame or timeout")
		if ret < 0 || err != nil {
			t.Error("CanSend failed")
		}

		if ret == 0 {
			t.Error("CanRecv timed out")
		}

		if ret != 16 {
			t.Error(fmt.Sprintf("CanRecv did not receive 16 bytes %d", ret))
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

		ret, err := CanSend(sock, &frame)
		if ret < 0 || err != nil {
			t.Error("CanSend failed")
		}
		t.Log("sent frame")

		ret = CanClose(sock)

		wg.Done()
	}

	var recv = func(wg *sync.WaitGroup) {
		var sock int
		sock = CanInit("vcan0")
		if sock < 0 {
			t.Error("CanInit failed")
		}

		var frame CanFrame

		ret, err := CanRecv(sock, &frame, 10000)
		t.Log("received frame or timeout")
		if ret < 0 || err != nil {
			t.Error("CanSend failed")
		}
		if ret == 0 {
			t.Error("CanRecv timed out")
		}

		if ret != 16 {
			t.Error(fmt.Sprintf("CanRecv did not receive 16 bytes %d", ret))
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
	time.Sleep(5 * time.Second)
	go send(&wg)
	wg.Wait()

}

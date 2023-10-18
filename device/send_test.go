package main_test

import (
	"testing"
	"time"

	"sqirvy.xyz/can"
)

func TestSend1(t *testing.T) {
	// TODO
	sockfd := can.CanInit("vcan0")
	defer can.CanClose(sockfd)

	var count byte = 0
	frame := can.CanFrame{}
	for count < 10 {
		frame.CanId = 1
		frame.CanDlc = 2
		frame.Data[0] = 0
		frame.Data[1] = count
		count++
		can.CanSend(sockfd, &frame)

		time.Sleep(1 * time.Second)
	}

}

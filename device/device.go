package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"sqirvy.xyz/can"
	"sqirvy.xyz/types"
)

type deviceState struct {
	dio_in  uint16
	dio_out uint16
	adc_in  int32
	adc_out int32
}

type CANFrame struct {
	adc_in atomic.Value
	dio_in atomic.Value

	dio_out atomic.Value
	adc_out atomic.Value
}

// the can bus recv function is blocking so run it
// in a goroutine and send the frame back to main for processing
func recvDevice(sockfd int, fch chan<- can.CanFrame, quit <-chan bool) {
	var frame can.CanFrame
	for {
		// receive with timeout
		ret := 2
		ret = can.CanRecv(sockfd, &frame, types.DEVICE_RECV_TIMEOUT)
		if ret < 0 {
			fmt.Printf("can.CanRecv() failed: %d %d %s\n", ret, can.CanErrno(), can.CanErrnoString())
			continue
		}
		if ret == 0 {
			// fmt.Printf("can.CanRecv() timeout: %d\n", ret)
			continue
		}
		select {
		case q := <-quit:
			if q {
				break
			}
		default:
		}
		fch <- frame
	}
}

func main() {
	sockfd := can.CanInit("vcan0")
	defer can.CanClose(sockfd)

	state := deviceState{0, 0, 0, 0}

	fch := make(chan can.CanFrame)
	defer close(fch)
	quit := make(chan bool)
	defer close(quit)

	go recvDevice(sockfd, fch, quit)

	prev_dio_out := uint16(0)
	prev_adc_out := int32(0)

	q := false
	for q == false {
		select {
		case frame := <-fch:
			// receive from client
			switch frame.CanId {
			case types.ID_DIO_SET:
				// DIO is uint16
				state.dio_in = can.BytesToUint16(frame.Data[0:2])
				print("DIO SET: ", state.dio_in, "\n")
			case types.ID_DAC_SET:
				// DAC is int32
				state.adc_in = can.BytesToInt32(frame.Data[0:4])
				print("ADC SET: ", state.adc_in, "\n")
			default:
			}
		case <-time.After(1000 * time.Millisecond):
			// update the simulator
			state.dio_out = state.dio_in
			state.adc_out = state.adc_in

			// send to client
			if state.dio_out != prev_dio_out {
				print("DIO OUT: ", state.dio_out, "\n")
				prev_dio_out = state.dio_out
			}
			frame16 := can.Uint16Frame(types.ID_DIO_OUT, state.dio_out)
			can.CanSend(sockfd, &frame16)

			if state.adc_out != prev_adc_out {
				print("ADC OUT: ", state.adc_out, "\n")
				prev_adc_out = state.adc_out
			}
			frame32 := can.Int32Frame(types.ID_ADC_OUT, state.adc_out)
			can.CanSend(sockfd, &frame32)
		}
	}
}

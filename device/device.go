package main

import (
	"log"
	"os"
	"time"

	"sqirvy.xyz/can"
	"sqirvy.xyz/types"
)

type deviceState struct {
	dio_set uint16
	dio_out uint16
	dac_set int32
	adc_out int32
}

// type CANFrame struct {
// 	adc_in atomic.Value
// 	dio_in atomic.Value

// 	dio_out atomic.Value
// 	adc_out atomic.Value
// }

// the can bus recv function is blocking so run it
// in a goroutine and send the frame back to main for processing
func recvDevice(sockfd int, fch chan<- can.CanFrame, quit <-chan bool) {
	var frame can.CanFrame
	var err error
loop:
	for {
		// receive with timeout
		ret := 2
		ret, err = can.CanRecv(sockfd, &frame, types.DEVICE_RECV_TIMEOUT)
		if ret < 0 || err != nil {
			log.Printf("can.CanRecv() failed: %d %d %s\n", ret, can.CanErrno(), can.CanErrnoString())
			continue
		}
		if ret == 0 {
			continue
		}
		select {
		case <-quit:
			break loop
		default:
		}
		fch <- frame
	}
}

func main() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	print("Starting Simulated Device\n")

	sockfd := can.CanInit("vcan0")
	defer can.CanClose(sockfd)

	state := deviceState{0, 0, 0, 0}

	fch := make(chan can.CanFrame)
	defer close(fch)
	quit := make(chan bool)
	defer close(quit)

	go recvDevice(sockfd, fch, quit)

	// prev_dio_out := uint16(0)
	// prev_adc_out := int32(0)

	q := false
	for !q {
		select {
		case frame := <-fch:
			// receive from client
			switch frame.CanId {
			case types.ID_DIO_SET:
				// DIO is uint16
				state.dio_set = can.BytesToUint16(frame.Data[0:2])
				print("set DIO SET: ", state.dio_set, "\n")
			case types.ID_DAC_SET:
				// DAC is int32
				state.dac_set = can.BytesToInt32(frame.Data[0:4])
				print("set ADC SET: ", state.dac_set, "\n")
			default:
			}
		case <-time.After(1000 * time.Millisecond):
			// update the simulator
			state.dio_out = state.dio_set
			state.adc_out = state.dac_set

			frame16 := can.Uint16Frame(types.ID_DIO_OUT, state.dio_out)
			ret, err := can.CanSend(sockfd, &frame16)
			if ret < 0 || err != nil {
				log.Printf("can.CanSend() failed: %d %d %s\n", ret, can.CanErrno(), can.CanErrnoString())
			}

			frame16 = can.Uint16Frame(types.ID_DIO_SET, state.dio_set)
			ret, err = can.CanSend(sockfd, &frame16)
			if ret < 0 || err != nil {
				log.Printf("can.CanSend() failed: %d %d %s\n", ret, can.CanErrno(), can.CanErrnoString())
			}

			frame32 := can.Int32Frame(types.ID_ADC_OUT, state.adc_out)
			ret, err = can.CanSend(sockfd, &frame32)
			if ret < 0 || err != nil {
				log.Printf("can.CanSend() failed: %d %d %s\n", ret, can.CanErrno(), can.CanErrnoString())
			}

			frame32 = can.Int32Frame(types.ID_DAC_SET, state.dac_set)
			ret, err = can.CanSend(sockfd, &frame32)
			if ret < 0 || err != nil {
				log.Printf("can.CanSend() failed: %d %d %s\n", ret, can.CanErrno(), can.CanErrnoString())
			}

			print(state.dio_out, " ", state.dio_set, " ", state.adc_out, " ", state.dac_set, "\n")
		}
	}
}

package device

import (
	"encoding/binary"
	"fmt"
	"time"

	"sqirvy.xyz/can"
)

type deviceState struct {
	dio_in  uint16
	dio_out uint16
	adc_in  int32
	adc_out int32
}

func BytesToUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func BytesToInt32(b []byte) int32 {
	return int32(binary.BigEndian.Uint32(b))
}

func Uint16ToBytes(v uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, v)
	return b
}

func Int32ToBytes(v int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v))
	return b
}

func uint16Frame(id uint32, v uint16) can.CanFrame {
	frame := can.CanFrame{}
	frame.CanId = id
	frame.CanDlc = 2
	b := Uint16ToBytes(v)
	frame.Data[0] = b[0]
	frame.Data[1] = b[1]

	return frame
}

func int32Frame(id uint32, v int32) can.CanFrame {
	frame := can.CanFrame{}
	frame.CanId = id
	frame.CanDlc = 4
	b := Int32ToBytes(v)
	frame.Data[0] = b[0]
	frame.Data[1] = b[1]
	frame.Data[2] = b[2]
	frame.Data[3] = b[3]

	return frame
}

// the can bus recv function is blocking so run it
// in a goroutine and send the frame back to main for processing
func recvDevice(sockfd int, fch chan<- can.CanFrame, quit <-chan bool) {
	var frame can.CanFrame
	for {
		// receive with timeout
		ret := can.CanRecv(sockfd, &frame, SIM_RECV_TIMEOUT)
		if ret < 0 {
			fmt.Printf("can.CanRecv() failed: %d\n", ret)
			continue
		}
		if ret == 0 {
			fmt.Printf("can.CanRecv() timeout: %d\n", ret)
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

	q := false
	for q == false {
		select {
		case frame := <-fch:
			// receive from client
			switch frame.CanId {
			case SIM_DIO_IN:
				// DIO is uint16
				state.dio_in = BytesToUint16(frame.Data[0:2])
			case SIM_DAC_IN:
				// DAC is int32
				state.adc_in = BytesToInt32(frame.Data[0:4])
			default:
			}
		case <-time.After(100 * time.Millisecond):
			// update the simulator
			state.dio_out = state.dio_in
			state.adc_out = state.adc_in
			// send to client
			frame16 := uint16Frame(SIM_DIO_OUT, state.dio_out)
			can.CanSend(sockfd, &frame16)

			frame32 := int32Frame(SIM_ADC_OUT, state.adc_out)
			can.CanSend(sockfd, &frame32)
		default:
		}
	}
}

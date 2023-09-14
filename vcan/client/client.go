package client

// "../g/sqirvy.xyz/can"

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"sqirvy.xyz/can"
)

type CanData struct {
	adc_in atomic.Value
	dio_in atomic.Value

	dio_out atomic.Value
	adc_out atomic.Value
}

var clientState CanSim

type CanSim struct {
	CanData
	t      time.Time
	d      time.Duration
	sockfd int
	mtx    sync.Mutex
}

func (s *CanSim) init(sock int) {
	s.dio_in.Store(0)
	s.dio_out.Store(0)
	s.adc_in.Store(0)
	s.adc_out.Store(0)
	s.t = time.Now()
	s.d = time.Duration(0)
	s.sockfd = sock
	s.mtx = sync.Mutex{}
}

func (s *CanSim) putDioIn(v uint16) {
	s.dio_in.Store(v)

	// send to can sim
}

func (s *CanSim) getDioOut() uint16 {
	v := s.dio_out.Load()
	return v.(uint16)
}

func (s *CanSim) putDacIn(v int32) {
	s.adc_in.Store(v)

	// send to can sim
}

func (s *CanSim) getAdcOut() int32 {
	v := s.adc_out.Load()
	return v.(int32)
}

// conversions (big endian)
func bytesToUint16(b []byte) uint16 {
	return (uint16(b[0]) << 8) | uint16(b[1])
}

func bytesToInt32(b []byte) int32 {
	return int32((uint32(b[0]) << 24) | (uint32(b[1]) << 16) | (uint32(b[2]) << 8) | uint32(b[3]))
}

func uint16ToBytes(v uint16) []byte {
	b := make([]byte, 2)
	b[0] = byte(v >> 8)
	b[1] = byte(v)
	return b
}

func int32ToBytes(v int32) []byte {
	b := make([]byte, 4)
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
	return b
}

// the can bus recv function is blocking so run it
// in a goroutine and send the frame back to main for processing
func receiver(sockfd int, fch chan<- can.CanFrame, quit <-chan bool) {
	var frame can.CanFrame
	for {
		// receive with timeout
		ret := can.CanRecv(sockfd, &frame, CLIENT_RECV_TIMEOUT)
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

// singleton socket, init to -1 to indicate not initialized
var sockfd int = -1

func main() {
	var frame can.CanFrame

	sockfd = can.CanInit("vcan0")
	defer can.CanClose(sockfd)

	simState := new(CanSim)
	simState.init(sockfd)

	fch := make(chan can.CanFrame)
	defer close(fch)
	quit := make(chan bool)
	defer close(quit)

	for {
		// receive with timeout
		ret := can.CanRecv(sockfd, &frame, CLIENT_RECV_TIMEOUT)
		if ret < 0 {
			fmt.Printf("can.CanRecv() failed: %d\n", ret)
			continue
		}
		if ret == 0 {
			fmt.Printf("can.CanRecv() timeout: %d\n", ret)
			continue
		}

		switch frame.CanId {
		case CLIENT_DIO_OUT:
			v := uint16(bytesToUint16(frame.Data[:]))
			simState.dio_out.Store(v)
		case CLIENT_ADC_OUT:
			v := int32(bytesToInt32(frame.Data[:]))
			simState.adc_out.Store(v)
		default:
		}
	}
}

// =========================
// exports to client
// =========================

func PutCanUint16(id int, v uint16) error {
	var err error
	var frame can.CanFrame
	switch id {
	case CLIENT_DIO_IN:
		clientState.putDioIn(v)
		// send to CAN bus
		frame.CanId = CLIENT_DIO_IN
		frame.CanDlc = 2
		b := uint16ToBytes(v)
		copy(frame.Data[:], b)
		can.CanSend(sockfd, &frame)
		err = nil
		break
	default:
		err = fmt.Errorf("invalid id")
		break
	}
	return err
}

func PutCanInt32(id int, v int32) error {
	var err error
	var frame can.CanFrame
	switch id {
	case CLIENT_DAC_IN:
		clientState.putDacIn(v)
		// send to CAN bus
		frame.CanId = CLIENT_DIO_IN
		frame.CanDlc = 2
		b := int32ToBytes(v)
		copy(frame.Data[:], b)
		can.CanSend(sockfd, &frame)
		err = nil
		break
	default:
		err = fmt.Errorf("invalid id")
		break
	}
	return err
}

func GetCanUint16(id int) (uint16, error) {
	var v uint16
	var err error
	switch id {
	case CLIENT_DIO_OUT:
		v = clientState.getDioOut()
		err = nil
	default:
		v = 0
		err = fmt.Errorf("invalid id")
	}
	return v, err
}

func GetCanInt32(id int) (int32, error) {
	var v int32
	var err error

	switch id {
	case CLIENT_ADC_OUT:
		v = clientState.getAdcOut()
		err = nil
	default:
		v = 0
		err = fmt.Errorf("invalid id")
	}

	return v, err
}

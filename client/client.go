package client

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"sqirvy.xyz/can"
	"sqirvy.xyz/types"
)

type CanData struct {
	adc_in atomic.Int32
	dio_in atomic.Uint32

	dio_out atomic.Uint32
	adc_out atomic.Int32
}

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

func (s *CanSim) putDioSet(v uint16) {
	s.dio_in.Store(uint32(v))
}

func (s *CanSim) getDioOut() uint16 {
	v := s.dio_out.Load()
	return uint16(v)
}

func (s *CanSim) putDacSet(v int32) {
	s.adc_in.Store(int32(v))
}

func (s *CanSim) getAdcOut() int32 {
	v := s.adc_out.Load()
	return v
}

// // conversions (big endian)
// func bytesToUint16(b []byte) uint16 {
// 	return (uint16(b[0]) << 8) | uint16(b[1])
// }

// func bytesToInt32(b []byte) int32 {
// 	return int32((uint32(b[0]) << 24) | (uint32(b[1]) << 16) | (uint32(b[2]) << 8) | uint32(b[3]))
// }

// func uint16ToBytes(v uint16) []byte {
// 	b := make([]byte, 2)
// 	b[0] = byte(v >> 8)
// 	b[1] = byte(v)
// 	return b
// }

// func int32ToBytes(v int32) []byte {
// 	b := make([]byte, 4)
// 	b[0] = byte(v >> 24)
// 	b[1] = byte(v >> 16)
// 	b[2] = byte(v >> 8)
// 	b[3] = byte(v)
// 	return b
// }

// the can bus recv function is blocking so run it
// in a goroutine and send the frame back to main for processing
func receiver(sockfd int, fch chan<- can.CanFrame, quit <-chan bool) {
	var frame can.CanFrame
	for {
		// receive with timeout
		ret := can.CanRecv(sockfd, &frame, types.CLIENT_RECV_TIMEOUT)
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
var simState = new(CanSim)

func Run() {
	var frame can.CanFrame

	sockfd = can.CanInit("vcan0")
	if sockfd < 0 {
		fmt.Printf("can.CanInit() failed: %d %s\n", sockfd, can.CanErrnoString())
		return
	}
	defer can.CanClose(sockfd)

	simState.init(sockfd)

	fch := make(chan can.CanFrame)
	defer close(fch)
	quit := make(chan bool)
	defer close(quit)

	for {
		// receive with timeout
		ret := can.CanRecv(sockfd, &frame, types.CLIENT_RECV_TIMEOUT)
		if ret < 0 {
			fmt.Printf("can.CanRecv() failed: %d %s\n", ret, can.CanErrnoString())
			continue
		}
		if ret == 0 {
			fmt.Printf("can.CanRecv() timeout: %d %s\n", ret, can.CanErrnoString())
			continue
		}

		switch frame.CanId {
		case types.ID_DIO_OUT:
			v := uint16(can.BytesToUint16(frame.Data[:]))
			simState.dio_out.Store(uint32(v))
		case types.ID_ADC_OUT:
			v := int32(can.BytesToInt32(frame.Data[:]))
			simState.adc_out.Store(v)
		default:
		}
	}
}

// =========================
// exports to api
// =========================

func PutCanUint16(id int, v uint16) error {
	var err error
	switch id {
	case types.ID_DIO_SET:
		simState.
			putDioSet(v)
		// send to CAN bus
		frame16 := can.Uint16Frame(types.ID_DIO_SET, v)
		can.CanSend(sockfd, &frame16)
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
	switch id {
	case types.ID_DAC_SET:
		simState.putDacSet(v)
		// send to CAN bus
		frame32 := can.Int32Frame(types.ID_DAC_SET, v)
		can.CanSend(sockfd, &frame32)
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
	case types.ID_DIO_OUT:
		v = simState.getDioOut()
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
	case types.ID_ADC_OUT:
		v = simState.getAdcOut()
		err = nil
	default:
		v = 0
		err = fmt.Errorf("invalid id")
	}

	return v, err
}

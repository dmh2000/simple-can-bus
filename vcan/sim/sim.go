package sim

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

type CanSim struct {
	CanData
	t      time.Time
	d      time.Duration
	sockfd int
	mtx    sync.Mutex
}

func (s *CanSim) Init(sock int) {
	s.dio_in.Store(0)
	s.dio_out.Store(0)
	s.adc_in.Store(0)
	s.adc_out.Store(0)
	s.t = time.Now()
	s.d = time.Duration(0)
	s.sockfd = sock
	s.mtx = sync.Mutex{}
}

func (s *CanSim) PutDioIn(v uint16) {
	s.dio_in.Store(v)

	// send to can sim
}

func (s *CanSim) GetDioOut() uint16 {
	v := s.dio_out.Load()
	return v.(uint16)
}

func (s *CanSim) PutDacIn(v int32) {
	s.adc_in.Store(v)

	// send to can sim
}

func (s *CanSim) GetAdcOut() int32 {
	v := s.adc_out.Load()
	return v.(int32)
}

// the can bus recv function is blocking so run it
// in a goroutine and send the frame back to main for processing
func receiver(sockfd int, fch chan<- can.CanFrame, quit <-chan bool) {
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

// exports to client
var simState CanSim

func PutCanUint16(id int, v uint16) {
	switch id {
	case SIM_DIO_IN:
		simState.PutDioIn(v)
		break
	default:
		break
	}
	simState.PutDioIn(v)
}

func PutCanInt32(id int, v int32) {
	switch id {
	case SIM_DAC_IN:
		simState.PutDacIn(v)
		break
	default:
		break
	}
}

func recv() {

	sockfd := can.CanInit("vcan0")
	defer can.CanClose(sockfd)

	simState := new(CanSim)
	simState.Init(sockfd)

	fch := make(chan can.CanFrame)
	defer close(fch)
	quit := make(chan bool)
	defer close(quit)

	go receiver(sockfd, fch, quit)

	q := false
	for q == false {
		select {
		case frame := <-fch:
			// receive from client
			switch frame.CanId {
			case SIM_DIO_IN:
				// DIO is uint16
				//u64, _ = binary.Uvarint(frame.Data[0:4])
				//simState.PutDioIn(uint16(u64))
			case SIM_DAC_IN:
				// DAC is int32
				//i64, _ = binary.Varint(frame.Data[:])
				//simState.PutAd(int32(i64))
			default:
			}
		default:
		}
	}
}

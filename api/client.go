package main

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"sqirvy.xyz/can"
	"sqirvy.xyz/types"
)

type CanData struct {
	DioSet atomic.Uint32
	DioOut atomic.Uint32
	DacSet atomic.Int32
	AdcOut atomic.Int32
}

func (s *CanData) init(sock int) {
	s.DioSet.Store(0)
	s.DioOut.Store(0)
	s.DacSet.Store(0)
	s.AdcOut.Store(0)
}

func (s *CanData) putDioSet(v uint16) {
	s.DioSet.Store(uint32(v))
}

func (s *CanData) getDioSet() uint16 {
	v := s.DioSet.Load()
	return uint16(v)
}

func (s *CanData) getDioOut() uint16 {
	v := s.DioOut.Load()
	return uint16(v)
}

func (s *CanData) putDacSet(v int32) {
	s.DacSet.Store(int32(v))
}

func (s *CanData) getDacSet() int32 {
	v := s.DacSet.Load()
	return v
}

func (s *CanData) getAdcOut() int32 {
	v := s.AdcOut.Load()
	return v
}

// state object
var canState = new(CanData)
var sockfd int = -1

func Run() {
	var frame types.CanFrame

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	sockfd = can.CanInit("vcan0")
	if sockfd < 0 {
		log.Printf("can.CanInit() failed: %d %s\n", sockfd, can.CanErrnoString())
		return
	}
	defer can.CanClose(sockfd)

	canState.init(sockfd)

	for {
		// receive with timeout
		ret, err := can.CanRecv(sockfd, &frame, types.CLIENT_RECV_TIMEOUT)
		if ret < 0 || err != nil {
			log.Printf("can.CanRecv() failed: %d %s\n", ret, can.CanErrnoString())
			continue
		}
		if ret == 0 {
			log.Printf("can.CanRecv() timeout: %d %s\n", ret, can.CanErrnoString())
			continue
		}

		switch frame.CanId {
		case types.ID_DIO_OUT:
			v := uint16(can.BytesToUint16(frame.Data[:]))
			canState.DioOut.Store(uint32(v))
		case types.ID_DIO_SET:
			v := uint16(can.BytesToUint16(frame.Data[:]))
			canState.DioSet.Store(uint32(v))
		case types.ID_ADC_OUT:
			v := int32(can.BytesToInt32(frame.Data[:]))
			canState.AdcOut.Store(v)
		case types.ID_DAC_SET:
			v := int32(can.BytesToInt32(frame.Data[:]))
			canState.DacSet.Store(v)
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
		canState.
			putDioSet(v)
		// send to CAN bus
		frame16 := can.Uint16Frame(types.ID_DIO_SET, v)
		ret, err := can.CanSend(sockfd, &frame16)
		if ret < 0 || err != nil {
			log.Printf("can.CanSend() failed: %d %s\n", ret, can.CanErrnoString())
		}
	default:
		err = fmt.Errorf("invalid id")
	}
	return err
}

func PutCanInt32(id int, v int32) error {
	var err error
	switch id {
	case types.ID_DAC_SET:
		canState.putDacSet(v)
		// send to CAN bus
		frame32 := can.Int32Frame(types.ID_DAC_SET, v)
		ret, err := can.CanSend(sockfd, &frame32)
		if ret < 0 || err != nil {
			log.Printf("can.CanSend() failed: %d %s\n", ret, can.CanErrnoString())
		}
	default:
		err = fmt.Errorf("invalid id")
	}
	return err
}

func GetCanUint16(id int) (uint16, error) {
	var v uint16
	var err error
	switch id {
	case types.ID_DIO_OUT:
		v = canState.getDioOut()
		err = nil
	case types.ID_DIO_SET:
		v = canState.getDioSet()
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
		v = canState.getAdcOut()
		err = nil
	case types.ID_DAC_SET:
		v = canState.getDacSet()
	default:
		v = 0
		err = fmt.Errorf("invalid id")
	}

	return v, err
}

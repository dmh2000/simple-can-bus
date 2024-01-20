package types

const (
	ID_DIO_SET = 1
	ID_DIO_OUT = 2
	ID_DAC_SET = 3
	ID_ADC_OUT = 4
)

type CanDevice struct {
	DioSet uint16
	DioOut uint16
	DacSet int32
	AdcOut int32
}

// note : sizeof canlib_frame is 16 bytes due to padding alignment
type CanFrame struct {
	CanId  uint32
	CanDlc byte
	Data   [8]byte
}

const (
	CLIENT_RECV_TIMEOUT = 1000
	CLIENT_SEND_TIMEOUT = 100
	CLIENT_TIMEOUT      = 10000
	DEVICE_RECV_TIMEOUT = 1000
	DEVICE_SEND_TIMEOUT = 100
	DEVICE_TIMEOUT      = 10000
)

package modbus

import (
	"testing"
	"unsafe"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestAduSize(t *testing.T) {
	var padu AduB
	if MODBUS_MBAP_SIZE+MODBUS_MAX_PDU_SIZE != len(padu) {
		t.Fatalf("adu size error")
	}
	// ModbusAdu size may be different from ModbusAduB size
	// due to structure padding. so sizeof(adu) will be
	// equal to or larger than sizeof(padu)
	var adu ADU
	if unsafe.Sizeof(adu) < unsafe.Sizeof(padu) {
		t.Fatalf("adu size error")
	}
}

func TestMbap(t *testing.T) {
	var adu ADU

	adu.Mbap.Tid = 1
	adu.Mbap.Pid = 0
	adu.Mbap.Len = 10
	adu.Mbap.Uid = 2
	adu.Pdu.Fc = 3
	adu.Pdu.Data[0] = 4
	adu.Pdu.Data[1] = 5
	adu.Pdu.Data[2] = 6
	adu.Pdu.Data[3] = 7
	adu.Pdu.Data[4] = 8
	adu.Pdu.Data[5] = 9
	adu.Pdu.Data[6] = 10
	adu.Pdu.Data[7] = 11

	b := AduPack(&adu)

	var adu2 ADU
	adu2.Mbap = MbapUnpack(b)
	adu2.Pdu = PduUnpack(int(adu2.Mbap.Len-2), b[MODBUS_MBAP_SIZE:])

	if adu2.Mbap.Tid != 1 {
		t.Fatalf("mbap tid error")
	}
	if adu2.Mbap.Pid != 0 {
		t.Fatalf("mbap pid error")
	}
	if adu2.Mbap.Len != 10 {
		t.Fatalf("mbap len error")
	}
	if adu2.Mbap.Uid != 2 {
		t.Fatalf("mbap uid error")
	}
	if adu2.Pdu.Fc != 3 {
		t.Fatalf("pdu fc error")
	}
	j := 4
	i := 0
	for ; i < 8; i++ {
		if adu2.Pdu.Data[i] != byte(j) {
			t.Fatalf("pdu data error")
		}
		j++
	}
}

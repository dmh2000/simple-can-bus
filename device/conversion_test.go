package main_test

import (
	"testing"

	can "sqirvy.xyz/can"
)

func TestBytesToUint16(t *testing.T) {
	b := make([]byte, 2)
	b[0] = 0
	b[1] = 1

	u := can.BytesToUint16(b)
	if u != 1 {
		t.Errorf("BytesToUint16: expected 1, got %d", u)
	}

	b[0] = 1
	b[1] = 0
	u = can.BytesToUint16(b)
	if u != 0x0100 {
		t.Errorf("BytesToUint16: expected 0xffff, got %d", u)
	}

	b[0] = 0xff
	b[1] = 0xff
	u = can.BytesToUint16(b)
	if u != 0xffff {
		t.Errorf("BytesToUint16: expected 0xffff, got %d", u)
	}
}

func TestBytesToInt32(t *testing.T) {
	var v int32

	b := make([]byte, 4)
	b[0] = 0
	b[1] = 0
	b[2] = 0
	b[3] = 1

	v = can.BytesToInt32(b)

	if v != 1 {
		t.Errorf("BytesToInt32: expected 1, got %d", v)
	}

	b[0] = 0x00
	b[1] = 0x00
	b[2] = 0xff
	b[3] = 0x00

	v = can.BytesToInt32(b)
	if v != 0xff00 {
		t.Errorf("BytesToInt32: expected 0xff00, got %d", v)
	}

	b[0] = 0x00
	b[1] = 0xff
	b[2] = 0x00
	b[3] = 0x00

	v = can.BytesToInt32(b)
	if v != 0xff0000 {
		t.Errorf("BytesToInt32: expected 0xff0000, got %d", v)
	}

	b[0] = 0xff
	b[1] = 0xff
	b[2] = 0xff
	b[3] = 0xff

	v = can.BytesToInt32(b)
	if v != -1 {
		t.Errorf("BytesToInt32: expected -1, got %d", v)
	}
}

func TestUint16ToBytes(t *testing.T) {
	var b []byte

	b = can.Uint16ToBytes(0)
	if b[0] != 0 || b[1] != 0 {
		t.Errorf("Uint16ToBytes: expected 0x0000, got 0x%x%x", b[0], b[1])
	}

	b = can.Uint16ToBytes(1)
	if b[0] != 0 || b[1] != 1 {
		t.Errorf("Uint16ToBytes: expected 0x0001, got 0x%x%x", b[0], b[1])
	}

	b = can.Uint16ToBytes(0x0100)
	if b[0] != 1 || b[1] != 0 {
		t.Errorf("Uint16ToBytes: expected 0x0100, got 0x%x%x", b[0], b[1])
	}

	b = can.Uint16ToBytes(0xffff)
	if b[0] != 0xff || b[1] != 0xff {
		t.Errorf("Uint16ToBytes: expected 0xffff, got 0x%x%x", b[0], b[1])
	}
}

func TestInt32ToBytes(t *testing.T) {
	b := can.Int32ToBytes(0)
	if b[0] != 0 || b[1] != 0 || b[2] != 0 || b[3] != 0 {
		t.Errorf("Int32ToBytes: expected 0x00000000, got 0x%x%x%x%x", b[0], b[1], b[2], b[3])
	}

	b = can.Int32ToBytes(1)
	if b[0] != 0 || b[1] != 0 || b[2] != 0 || b[3] != 1 {
		t.Errorf("Int32ToBytes: expected 0x00000001, got 0x%x%x%x%x", b[0], b[1], b[2], b[3])
	}

	b = can.Int32ToBytes(0xff00)
	if b[0] != 0 || b[1] != 0 || b[2] != 0xff || b[3] != 0 {
		t.Errorf("Int32ToBytes: expected 0x0000ff00, got 0x%x%x%x%x", b[0], b[1], b[2], b[3])
	}

	b = can.Int32ToBytes(0xff0000)
	if b[0] != 0 || b[1] != 0xff || b[2] != 0 || b[3] != 0 {
		t.Errorf("Int32ToBytes: expected 0x00ff0000, got 0x%x%x%x%x", b[0], b[1], b[2], b[3])
	}

	b = can.Int32ToBytes(-1)
	if b[0] != 0xff || b[1] != 0xff || b[2] != 0xff || b[3] != 0xff {
		t.Errorf("Int32ToBytes: expected 0xffffffff, got 0x%x%x%x%x", b[0], b[1], b[2], b[3])
	}
}

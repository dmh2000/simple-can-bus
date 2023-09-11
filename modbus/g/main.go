package main

import (
	"net"

	"sqirvy.xyz/modbus/modbus"
)

type modbusServerSim struct {
	value int
}

func (sim *modbusServerSim) Server(conn *net.TCPConn) {
	var adu modbus.ADU
	b := make([]byte, modbus.MODBUS_MBAP_SIZE)
	v, err := conn.Read(b)
	if err != nil || v != modbus.MODBUS_MBAP_SIZE {
		print("error reading mbap\n")
		conn.Close()
		return
	}
	adu.Mbap = modbus.MbapUnpack(b)

	b = make([]byte, adu.Mbap.Len-1)
	v, err = conn.Read(b)
	if err != nil || v != int(adu.Mbap.Len-1) {
		print("error reading pdu\n")
		conn.Close()
		return
	}

	adu.Pdu = modbus.PduUnpack(int(adu.Mbap.Len-1), b)
}

func main() {
	var handler modbus.ServerHandler
	handler = &modbusServerSim{value: 0}

	modbus.Server("127.0.01", 5000, handler)
}

package modbus

import "net"

const (
	// packed sizes
	MODBUS_MBAP_SIZE         = 7
	MODBUS_MAX_PDU_DATA_SIZE = 65536 - 7
	MODBUS_MAX_PDU_SIZE      = 65536 - MODBUS_MBAP_SIZE
)

// MBAP : modbus application protocol
type MBAP struct {
	Tid uint16 // transaction id : echo back with response
	Pid uint16 // protocol id : 0 for modbus
	Len uint16 // length : uid + fc + data
	Uid uint8  // unit id : tcp rs485 gateway
}
type MbapB = [MODBUS_MBAP_SIZE]byte

// PDU : protocol data unit
type PDU struct {
	Fc   uint8                          // function code
	Data [MODBUS_MAX_PDU_DATA_SIZE]byte // data
}
type PduB = [MODBUS_MAX_PDU_SIZE]byte

type ADU struct {
	Mbap MBAP
	Pdu  PDU
}

// max message size
type AduB = [MODBUS_MBAP_SIZE + MODBUS_MAX_PDU_SIZE]byte

// pack an instance of ADU including MBAP and PDU
func AduPack(adu *ADU) []byte {
	padu := make([]byte, adu.Mbap.Len+MODBUS_MBAP_SIZE-1)
	padu[0] = byte(adu.Mbap.Tid >> 8)
	padu[1] = byte(adu.Mbap.Tid)
	padu[2] = 0
	padu[3] = 0
	padu[4] = byte(adu.Mbap.Len >> 8)
	padu[5] = byte(adu.Mbap.Len & 0xff)
	padu[6] = adu.Mbap.Uid
	padu[7] = adu.Pdu.Fc
	copy(padu[8:], adu.Pdu.Data[:adu.Mbap.Len-1])
	return padu
}

// unpacking an adu requires reading the mbap first to get the length

// unpack an instance of MBAP
func MbapUnpack(pmbap []byte) MBAP {
	if len(pmbap) < MODBUS_MBAP_SIZE {
		panic("mbap size error")
	}
	mbap := MBAP{}
	mbap.Tid = uint16(pmbap[0])<<8 | uint16(pmbap[1])
	mbap.Pid = uint16(pmbap[2])<<8 | uint16(pmbap[3])
	mbap.Len = uint16(pmbap[4])<<8 | uint16(pmbap[5])
	mbap.Uid = pmbap[6]
	return mbap
}

// unpack an instance of PDU
func PduUnpack(data_length int, ppdu []byte) PDU {
	if len(ppdu) < data_length {
		panic("pdu size error")
	}
	pdu := PDU{}
	pdu.Fc = ppdu[0]
	copy(pdu.Data[:], ppdu[1:data_length+1])
	return pdu
}

type ServerHandler interface {
	Server(conn *net.TCPConn)
}

type ClientHandler interface {
	Client(conn *net.TCPConn)
}

// server listens on ip:port and calls handler.Server()
// this version is blocking and intended to represent one modbus connection
// if the server needs to handle multiple connections, then it should
// separate the listener from the clients and
func Server(ip string, port int, handler ServerHandler) {
	addr := net.TCPAddr{Port: port, IP: net.ParseIP(ip)}
	listener, err := net.ListenTCP("tcp", &addr)
	defer listener.Close()

	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.AcceptTCP()
		defer conn.Close()
		if err != nil {
			panic(err)
		}
		go handler.Server(conn)
	}
}

/**
 * Modbus client
 * connects to server and calls handler.client
 * this is a blocking call intended to represent on modbus connection
 */

// server listens on ip:port and calls handler.Server() for each connection
func client(ip string, port int, handler ClientHandler) {
	remoteAddr := net.TCPAddr{Port: port, IP: net.ParseIP(ip)}
	conn, err := net.DialTCP("tcp", nil, &remoteAddr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	for {
		handler.Client(conn)
	}
}

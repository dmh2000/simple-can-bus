#ifndef MODBUS_H__
#define MODBUS_H__
#include <stdint.h>

//                           max adu size - (mbap + fc)
#define MODBUS_MAX_PDU_SIZE (65536 - (sizeof(mbtcp_mbap_t) + sizeof(uint8_t)))

// Important : the modbus protocol uses big-endian representation. This library
// allows the calling function to use local byte order. The library will
// convert to big-endian when sending and convert back to local byte order
#pragma pack(1)
typedef struct mbtcp_mbap_t {
    uint16_t tid; // transaction id : echo back with response
    uint16_t pid; // protocol id : 0 for modbus
    uint16_t len; // length : uid + fc + data
    uint8_t uid;  // unit id : tcp rs485 gateway
} mbtcp_mbap_t;

typedef struct mbtcp_pdu_s {
    uint8_t fc;                        // function code
    uint8_t data[MODBUS_MAX_PDU_SIZE]; // data
} mbtcp_pdu_t;

typedef struct mbtcp_adu_s {
    mbtcp_mbap_t mbap; // modbus application protocol header
    mbtcp_pdu_t pdu;   // modbus protocol data unit
} mbtcp_adu_t;
#pragma pack()

// TCP client socket connect
// @return socket descriptor or -1 on error
int mbtcp_client_connect(const char *server_ip, uint16_t server_portid);

// Important
// this implementation only services one client connection to one server
// at a time. If you want to service multiple clients, you will need to
// implement a select() loop or use threads.

// TCP server bind-listen-accept
// @return socket descriptor or -1 on error
int mbtcp_server_listen(const char *local_ip, uint16_t local_port);

// TCP send
// @return 0 on success or -1 on error
int mbtcp_send(int sockfd, mbtcp_adu_t *adu);

// TCP receive
// @return 0 if tiemout, >0 if success, -1 if fail
int mbtcp_recv(int sockfd, mbtcp_adu_t *adu, uint32_t timeout_ms);

// Print modbus adu
void mbtcp_print_adu(const mbtcp_adu_t *adu);

#endif // MODBUS_H__

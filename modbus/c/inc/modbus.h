#ifndef MODBUS_H__
#define MODBUS_H__
#include <stdint.h>

// modbus 485 supports 253 byte pdus, modbus tcp supports 65535 byte pdus
//                           max adu size - (mbap + fc)
#define MODBUS_MAX_PDU_DATA_SIZE                                               \
    (65536 - (sizeof(mbtcp_mbap_t) + sizeof(uint8_t)))

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

// for this implementation, the pdu struct is defined as supporting the full
// amount of data bytes. This is to keep it simple by not requiring dynamically
// allocate memory for the data. this does mean that if a lot of ADU or PDU
// structs are created, it may be less efficient for memory usage
typedef struct mbtcp_pdu_s {
    uint8_t fc;                             // function code
    uint8_t data[MODBUS_MAX_PDU_DATA_SIZE]; // data
} mbtcp_pdu_t;

// mbap and pdu are combined to form the adu. they are separated here because
// receiving an adu requires reading the mbap first to get the data length, and
// then the pdu.
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

// TCP server bind-listen
// @return server socket descriptor or -1 on error
int mbtcp_server_listen(const char *local_ip, uint16_t local_port);

// TCP server accept
// @return c;oemt socket descriptor or -1 on error
int mbtcp_server_accept(int sockfd);

// TCP close socket
// @return 0 on success or -1 on error
int mbtcp_close(int sockfd);

// TCP send
// @return 0 on success or -1 on error
int mbtcp_send(int sockfd, mbtcp_adu_t *adu);

// TCP receive
// @return 0 if tiemout, >0 if success, -1 if fail
int mbtcp_recv(int sockfd, mbtcp_adu_t *adu, uint32_t timeout_ms);

// Print modbus adu
void mbtcp_print_adu(const mbtcp_adu_t *adu);

// modbus function codes
/* function codes (from spec)
6.1 01 (0x01) Read Coils
6.2 02 (0x02) Read Discrete Inputs
6.3 03 (0x03) Read Holding Registers
6.4 04 (0x04) Read Input Registers
6.5 05 (0x05) Write Single Coil
6.6 06 (0x06) Write Single Register
*/
#define MODBUS_FC_READ_COILS 0x01
#define MODBUS_FC_READ_DISCRETE_INPUTS 0x02
#define MODBUS_FC_READ_HOLDING_REGISTERS 0x03
#define MODBUS_FC_READ_INPUT_REGISTERS 0x04
#define MODBUS_FC_WRITE_SINGLE_COIL 0x05
#define MODBUS_FC_WRITE_SINGLE_REGISTER 0x06

// modbus read coils request 0x01
typedef struct mbtcp_read_coils_req_s {
    uint16_t addr; // address
    uint16_t cnt;  // count
} mbtcp_read_coils_t;

// modbus read coils response
typedef struct mbtcp_read_coils_rsp_s {
    uint8_t byte_cnt; // byte count
    // addressed on/off bits flow from low to high
    uint8_t data[MODBUS_MAX_PDU_DATA_SIZE - sizeof(uint8_t)];
} mbtcp_read_coils_rsp_t;

// modbus read discrete inputs request 0x02
typedef struct mbtcp_read_discrete_inputs_req_s {
    uint16_t addr; // address
    uint16_t cnt;  // count
} mbtcp_read_discrete_inputs_t;

// modbus read discrete inputs response
typedef struct mbtcp_read_discrete_inputs_rsp_s {
    uint8_t byte_cnt; // byte count
    // addressed on/off bits flow from low to high
    uint8_t data[MODBUS_MAX_PDU_DATA_SIZE - sizeof(uint8_t)];
} mbtcp_read_discrete_inputs_rsp_t;

// modbus read holding registers request 0x03
typedef struct mbtcp_read_holding_registers_req_s {
    uint16_t addr; // address
    uint16_t cnt;  // count
} mbtcp_read_holding_registers_t;

// modbus read holding registers response
typedef struct mbtcp_read_holding_registers_rsp_s {
    uint8_t byte_cnt; // byte count
    uint16_t reg[(MODBUS_MAX_PDU_DATA_SIZE - sizeof(uint8_t)) / 2];
} mbtcp_read_holding_registers_rsp_t;

// modbus read input registers request 0x04
typedef struct mbtcp_read_input_registers_req_s {
    uint16_t addr; // address
    uint16_t cnt;  // count
} mbtcp_read_input_registers_t;

// modbus read input registers response
typedef struct mbtcp_read_input_registers_rsp_s {
    uint8_t byte_cnt; // byte count
    // each register is two bytes in big-endian order
    uint16_t reg[(MODBUS_MAX_PDU_DATA_SIZE - sizeof(uint8_t)) / 2];
} mbtcp_read_input_registers_rsp_t;

// modbus write single coil request 0x05
typedef struct mbtcp_write_single_coil_req_s {
    uint16_t addr; // address
    uint16_t val;  // value (0xff00 = on, 0x0000 = off)
} mbtcp_write_single_coil_t;

// modbus write single coil response
typedef struct mbtcp_write_single_coil_rsp_s {
    uint16_t addr; // address
    uint16_t val;  // value
} mbtcp_write_single_coil_rsp_t;

// modbus write single register request 0x06
typedef struct mbtcp_write_single_register_req_s {
    uint16_t addr; // address
    uint16_t val;
} mbtcp_write_single_register_t;

// modbus write single register response
typedef struct mbtcp_write_single_register_rsp_s {
    uint16_t addr; // address
    uint16_t val;  // value
} mbtcp_write_single_register_rsp_t;

#endif // MODBUS_H__

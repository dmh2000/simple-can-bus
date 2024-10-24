#ifndef VCANLIB_H__
#define VCANLIB_H__

#include <stdint.h>

/* Error codes (all negative) */
#define CANLIB_OK             0
#define CANLIB_ERR_PARAM     -1  /* Invalid parameter */
#define CANLIB_ERR_SOCKET    -2  /* Socket creation failed */
#define CANLIB_ERR_BIND      -3  /* Socket bind failed */
#define CANLIB_ERR_INTERFACE -4  /* CAN interface error */
#define CANLIB_ERR_IO        -5  /* I/O error */
#define CANLIB_ERR_TIMEOUT   -6  /* Operation timed out */

#define CAN_MAX_DATA_LEN 8

/**
 * CAN frame structure modeled after linux can.h/can_frame
 * Provides a simplified interface for CAN communication
 */
typedef struct canlib_frame
{
    uint32_t can_id;           /* 32 bit CAN_ID + EFF/RTR/ERR flags */
    uint8_t can_dlc;           /* frame payload length in byte (0 .. CAN_MAX_DATA_LEN) */
    unsigned char data[CAN_MAX_DATA_LEN];  /* frame payload data */
} canlib_frame_t;

/**
 * Initialize a CAN interface
 * 
 * @param can_dev Name of CAN interface (e.g. "can0")
 * @return Positive socket descriptor on success, negative error code on failure
 */
int canlib_init(const char *can_dev);

/**
 * Receive a CAN frame with timeout
 * 
 * @param can_sock Socket descriptor from canlib_init
 * @param frame Pointer to frame structure to store received data
 * @param timeout_ms Timeout in milliseconds
 * @return Number of bytes read on success, CANLIB_ERR_TIMEOUT on timeout, negative error code on failure
 */
int canlib_receive(int can_sock, canlib_frame_t *frame, int timeout_ms);

/**
 * Send a CAN frame
 * 
 * @param can_sock Socket descriptor from canlib_init
 * @param frame Pointer to frame structure containing data to send
 * @return Number of bytes sent on success, negative error code on failure
 */
int canlib_send(int can_sock, const canlib_frame_t *frame);

/**
 * Close a CAN socket
 * 
 * @param can_sock Socket descriptor to close
 * @return CANLIB_OK on success, negative error code on failure
 */
int canlib_close(int can_sock);

/**
 * Get the last system error status
 * 
 * @return System errno value as uint32_t
 */
uint32_t canlib_status(void);

#endif // VCANLIB_H__

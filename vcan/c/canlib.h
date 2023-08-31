#ifndef VCANLIB_H__
#define VCANLIB_H__

#include <stdint.h>

#define CAN_MAX_DATA_LEN 8

// CAN frame structure
// modeled after linux can.h/can_frame repeated for convenience
typedef struct canlib_frame
{
    uint32_t can_id; /* 32 bit CAN_ID + EFF/RTR/ERR flags */
    uint8_t can_dlc; /* frame payload length in byte (0 .. CAN_MAX_DATA_LEN) */
    unsigned char data[CAN_MAX_DATA_LEN];
} canlib_frame_t;

struct canlib_frame;

int canlib_init(const char *can_dev);
int canlib_receive(int can_sock, canlib_frame_t *frame);
int canlib_send(int can_sock, canlib_frame_t *frame);
int32_t canlib_close(int can_sock);

#endif // VCANLIB_H__

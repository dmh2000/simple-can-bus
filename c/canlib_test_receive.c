#include <stdio.h>
#include <errno.h>
#include "canlib.h"

int main(int argc, char *argv[])
{
    int can_sock;
    int status;
    canlib_frame_t frame;

    can_sock = canlib_init("vcan0");
    if (can_sock < 0) {
        fprintf(stderr, "Error initializing CAN interface: %d\n", can_sock);
        return 1;
    }

    while (1) {
        status = canlib_receive(can_sock, &frame, 20000);
        if (status == CANLIB_ERR_TIMEOUT) {
            printf("Receive timeout\n");
            continue;
        }
        if (status < 0) {
            fprintf(stderr, "Error receiving CAN frame: %d (errno: %u)\n", 
                    status, canlib_status());
            canlib_close(can_sock);
            return 1;
        }

        printf("Received frame: ID=%04x, DLC=%d, Data=", frame.can_id, frame.can_dlc);
        for (int i = 0; i < frame.can_dlc; i++) {
            printf("%02x", frame.data[i]);
        }
        printf("\n");
    }

    canlib_close(can_sock);
    return 0;
}

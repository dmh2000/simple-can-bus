#include <stdio.h>
#include <errno.h>
#include <unistd.h>
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

    frame.can_id = 1;
    frame.can_dlc = 8;
    for (int i = 0; i < CAN_MAX_DATA_LEN; i++) {
        frame.data[i] = i;
    }

    for (;;) {
        status = canlib_send(can_sock, &frame);
        if (status < 0) {
            fprintf(stderr, "Error sending CAN frame: %d (errno: %u)\n", 
                    status, canlib_status());
            canlib_close(can_sock);
            return 1;
        }
        printf("Sent frame: ID=%04x, DLC=%d\n", frame.can_id, frame.can_dlc);
        sleep(1);
    }

    canlib_close(can_sock);

    return 0;
}

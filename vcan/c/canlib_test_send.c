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
    if (can_sock < 0)
    {
        printf("Error initializing CAN interface %s\n", argv[1]);
        return 1;
    }

    frame.can_id = 1;
    frame.can_dlc = 8;
    for (int i = 0; i < 8; i++)
    {
        frame.data[i] = i;
    }

    for (;;)
    {
        status = canlib_send(can_sock, &frame);
        if (status < 0)
        {
            printf("Error reading CAN frame\n");
            return 1;
        }
        sleep(1);
    }

    return 0;
}

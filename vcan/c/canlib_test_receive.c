#include <stdio.h>
#include <errno.h>
#include "canlib.h"

int main(int argc, char *argv[])
{
    int can_sock;
    int status;
    canlib_frame_t frame;

    if (argc < 2)
    {
        printf("Usage: %s <CAN interface>\n", argv[0]);
        return 1;
    }
    printf("%s\n", argv[1]);
    can_sock = canlib_init(argv[1]);
    if (can_sock < 0)
    {
        printf("Error initializing CAN interface %d:%s\n", errno, argv[1]);
        return 1;
    }

    while (1)
    {
        status = canlib_receive(can_sock, &frame);
        if (status < 0)
        {
            perror("Error reading CAN frame\n");
            return 1;
        }
        printf("%d %04x %04x\n", status, frame.can_id, frame.can_dlc);
        for (int i = 0; i < frame.can_dlc; i++)
        {
            printf("%02x", frame.data[i]);
        }
        printf("\n");
    }

    return 0;
}
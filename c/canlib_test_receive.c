#include <stdio.h>
#include <errno.h>
#include <signal.h>
#include <string.h>
#include "canlib.h"

static volatile sig_atomic_t keep_running = 1;

static void sig_handler(int signo) {
    keep_running = 0;
}

static void setup_signals(void) {
    struct sigaction sa;
    memset(&sa, 0, sizeof(sa));
    sa.sa_handler = sig_handler;
    sigaction(SIGINT, &sa, NULL);
    sigaction(SIGTERM, &sa, NULL);
}

int main(int argc, char *argv[]) {
    int can_sock;
    int status;
    canlib_frame_t frame;
    const char *interface = "vcan0";

    if (argc > 1) {
        interface = argv[1];
    }

    setup_signals();

    can_sock = canlib_init(interface);
    if (can_sock < 0) {
        fprintf(stderr, "Error initializing CAN interface: %d\n", can_sock);
        return 1;
    }

    while (keep_running) {
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

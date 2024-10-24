#include <stdio.h>
#include <errno.h>
#include <unistd.h>
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

    frame.can_id = 1;
    frame.can_dlc = 8;
    for (int i = 0; i < CAN_MAX_DATA_LEN; i++) {
        frame.data[i] = i;
    }

    while (keep_running) {
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

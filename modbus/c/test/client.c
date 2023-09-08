#include "../inc/modbus.h"
#include <arpa/inet.h>
#include <errno.h>
#include <stdint.h>
#include <stdio.h>
#include <unistd.h>
int main(int argc, char *argv[]) {
    int status;
    const char *ip = "127.0.0.1";
    uint16_t port = 5000;

    int sockfd = mbtcp_client_connect(ip, port);
    if (sockfd == -1) {
        perror("mbtcp_client_connect error\n");
        return -1;
    }

    for (;;) {
        mbtcp_adu_t adu;
        adu.mbap.tid = 1;
        adu.mbap.pid = 0;
        adu.mbap.len = 5;
        adu.mbap.uid = 0;
        adu.pdu.fc = 3;
        adu.pdu.data[0] = 1;
        adu.pdu.data[1] = 2;
        adu.pdu.data[2] = 3;

        status = mbtcp_send(sockfd, &adu);
        if (status < 0) {
            perror("mbtcp_send error\n");
            return -1;
        }
        printf("sent %d bytes\n", status);

        status = mbtcp_recv(sockfd, &adu, 10000);
        if (status < 0) {
            perror("mbtcp_recv");
            return -1;
        }
        printf("recv %d bytes\n", status);

        if (status == 0) {
            printf("timeout\n");
            continue;
        }

        mbtcp_print_adu(&adu);
        sleep(1);
    }

    return 0;
}
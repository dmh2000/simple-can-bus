#include "../inc/modbus.h"
#include <stdio.h>

mbtcp_adu_t modbus_sim(mbtcp_adu_t *adu) {
    mbtcp_adu_t adu_out;
    adu_out.mbap.tid = adu->mbap.tid;
    adu_out.mbap.pid = 0;
    adu_out.mbap.len = 5;
    adu_out.pdu.fc = adu->pdu.fc;
    adu_out.pdu.data[0] = 3;
    adu_out.pdu.data[1] = 2;
    adu_out.pdu.data[2] = 1;
    return adu_out;
}

int main(int argc, char *argv[]) {
    int sockfd;
    int status;
    mbtcp_adu_t adu_in;
    mbtcp_adu_t adu_out;

    const char *ip = "127.0.0.1";
    uint16_t port = 5000;
    sockfd = mbtcp_server_listen(ip, port);
    if (sockfd < 0) {
        perror("mtcp_server_init");
        return -1;
    }

    for (;;) {
        int clientfd = mbtcp_server_accept(sockfd);
        if (clientfd < 0) {
            perror("mbtcp_server_accept");
            return -1;
        }

        for (;;) {
            status = mbtcp_recv(clientfd, &adu_in, 10000);
            if (status < 0) {
                perror("mbtcp_recv");
                break;
            }

            if (status == 0) {
                printf("timeout\n");
                break;
            }

            mbtcp_print_adu(&adu_in);

            adu_out = modbus_sim(&adu_in);

            status = mbtcp_send(clientfd, &adu_out);
            if (status < 0) {
                perror("mbtcp_send");
                break;
            }
        }
        mbtcp_close(clientfd);
    }
}
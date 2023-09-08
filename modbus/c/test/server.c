#include "../inc/modbus.h"
#include <arpa/inet.h>
#include <errno.h>
#include <netinet/in.h>
#include <netinet/tcp.h>
#include <stdio.h>
#include <sys/socket.h>
#include <sys/time.h>
#include <sys/types.h>
#include <unistd.h>

int main(int argc, char *argv[]) {
    int sockfd;
    int status;
    mbtcp_adu_t adu;

    const char *ip = "127.0.0.1";
    uint16_t port = 5000;
    sockfd = mbtcp_server_listen(ip, port);
    if (sockfd < 0) {
        perror("mtcp_server_init");
        return -1;
    }

    for (;;) {
        struct sockaddr_in client_addr;
        socklen_t client_addr_len = sizeof(client_addr);
        int clientfd =
            accept(sockfd, (struct sockaddr *)&client_addr, &client_addr_len);
        if (status < 0) {
            perror("accept");
            return -1;
        }
        printf("client connected %08x:%u\n", client_addr.sin_addr.s_addr,
               client_addr.sin_port);

        for (;;) {
            status = mbtcp_recv(clientfd, &adu, 10000);
            if (status < 0) {
                perror("mbtcp_recv");
                return -1;
            }

            if (status == 0) {
                printf("timeout\n");
                continue;
            }

            mbtcp_print_adu(&adu);

            adu.mbap.tid = adu.mbap.tid;
            adu.mbap.pid = 0;
            adu.mbap.len = 5;
            adu.pdu.fc = adu.pdu.fc;
            adu.pdu.data[0] = 3;
            adu.pdu.data[1] = 2;
            adu.pdu.data[2] = 1;
            status = mbtcp_send(clientfd, &adu);
            if (status < 0) {
                perror("mbtcp_send");
                return -1;
            }
        }
    }
}
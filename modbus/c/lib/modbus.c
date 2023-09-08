#include "../inc/modbus.h"
#include <arpa/inet.h>
#include <netinet/in.h>
#include <netinet/tcp.h>
#include <stdio.h>
#include <sys/socket.h>
#include <sys/time.h>
#include <sys/types.h>
#include <unistd.h>

// TCP client connect
// @return socket descriptor or -1 on error
int mbtcp_client_connect(const char *server_ip, uint16_t server_port) {
    int status;

    int sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd == -1) {
        return -1;
    }

    // server address
    struct sockaddr_in addr;
    addr.sin_family = AF_INET;
    addr.sin_port = htons(server_port);
    status = inet_pton(AF_INET, server_ip, &addr.sin_addr);
    if (status < 0) {
        return -1;
    }

    status = connect(sockfd, (struct sockaddr *)&addr, sizeof(addr));
    if (status < 0) {
        return -1;
    }

    return sockfd;
}

// TCP server bind-listen-accept
// @return socket descriptor or -1 on error
int mbtcp_server_listen(const char *local_ip, uint16_t local_port) {
    int status;

    // tcp socket
    int sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd < 0) {
        return -1;
    }

    // local address
    struct sockaddr_in addr;
    addr.sin_family = AF_INET;
    addr.sin_port = htons(local_port);
    status = inet_pton(AF_INET, local_ip, &addr.sin_addr);
    if (status < 0) {
        return -1;
    }

    // bind
    status = bind(sockfd, (struct sockaddr *)&addr, sizeof(addr));
    if (status < 0) {
        return -1;
    }

    // listen
    status = listen(sockfd, 1);
    if (status < 0) {
        return -1;
    }

    return sockfd;
}

// TCP send
// @return 0 on success or -1 on error
int mbtcp_send(int sockfd, mbtcp_adu_t *adu) {
    int status;

    // switch to network byte order
    size_t len = (sizeof(mbtcp_mbap_t) - sizeof(uint8_t)) + adu->mbap.len;
    adu->mbap.tid = htons(adu->mbap.tid);
    adu->mbap.pid = htons(adu->mbap.pid);
    adu->mbap.len = htons(adu->mbap.len);
    adu->mbap.uid = adu->mbap.uid;
    adu->pdu.fc = adu->pdu.fc;

    status = send(sockfd, adu, len, 0);
    if (status < 0) {
        return -1;
    }

    // switch back to local byte order
    adu->mbap.tid = htons(adu->mbap.tid);
    adu->mbap.pid = htons(adu->mbap.pid);
    adu->mbap.len = htons(adu->mbap.len);
    adu->mbap.uid = adu->mbap.uid;
    adu->pdu.fc = adu->pdu.fc;

    return status;
}

// TCP receive
// @return 0 on success or -1 on error
int mbtcp_recv(int sockfd, mbtcp_adu_t *adu, uint32_t timeout_ms) {
    int status;
    fd_set fdread;
    struct timeval tv;
    tv.tv_sec = timeout_ms / 1000;
    tv.tv_usec = (timeout_ms % 1000) * 1000;

    FD_ZERO(&fdread);
    FD_SET(sockfd, &fdread);
    status = select(sockfd + 1, &fdread, 0, 0, &tv);
    if (status < 0) {
        return -1;
    }

    if (status == 0) {
        return 0;
    }

    if (FD_ISSET(sockfd, &fdread)) {
        // read mbap first to get length
        ssize_t rlen;
        rlen = recv(sockfd, &adu->mbap, sizeof(adu->mbap), 0);
        if (status < 0) {
            rlen - 1;
        }

        // switch to local byte order
        uint16_t len = adu->mbap.len - 1;
        adu->mbap.tid = htons(adu->mbap.tid);
        adu->mbap.pid = htons(adu->mbap.pid);
        adu->mbap.len = htons(adu->mbap.len);
        adu->mbap.uid = adu->mbap.uid;

        // read pdu
        ssize_t plen;
        plen = recv(sockfd, &adu->pdu, len, 0);
        if (status < 0) {
            return plen;
        }

        return rlen + plen;
    }
}

void mbtcp_print_adu(const mbtcp_adu_t *adu) {
    printf("tid:%u pid;%u len:%u uid:%u fc:%u\n", adu->mbap.tid, adu->mbap.pid,
           adu->mbap.len, adu->mbap.uid, adu->pdu.fc);
    printf("data: ");
    for (int i = 0; i < adu->mbap.len - 2; i++) {
        printf("%u ", adu->pdu.data[i]);
    }
    printf("\n");
}

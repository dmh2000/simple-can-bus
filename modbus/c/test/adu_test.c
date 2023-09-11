#include "../inc/modbus.h"
#include <stddef.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
    printf("mbap %lu %lu %lu %lu %lu \n", sizeof(mbtcp_mbap_t),
           offsetof(mbtcp_mbap_t, tid), offsetof(mbtcp_mbap_t, pid),
           offsetof(mbtcp_mbap_t, len), offsetof(mbtcp_mbap_t, uid));

    printf("pdu  %lu %lu %lu\n", sizeof(mbtcp_pdu_t), offsetof(mbtcp_pdu_t, fc),
           offsetof(mbtcp_pdu_t, data));

    printf("adu  %lu %lu\n", sizeof(mbtcp_adu_t),
           sizeof(mbtcp_mbap_t) + sizeof(mbtcp_pdu_t));
    return 0;
}
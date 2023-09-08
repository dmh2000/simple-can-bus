#include <stdio.h>
#include "../inc/modbus.h"

int main(int argc, char *argv[])
{
    printf("mbap %lu\n", sizeof(mbtcp_mbap_t));
    printf("pdu  %lu\n", sizeof(mbtcp_pdu_t));
    printf("adu  %lu %lu\n", sizeof(mbtcp_adu_t), sizeof(mbtcp_mbap_t) + sizeof(mbtcp_pdu_t));
    return 0;
}
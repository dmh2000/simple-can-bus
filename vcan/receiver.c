#include <sys/socket.h>
#include <sys/ioctl.h>
#include <linux/if.h>
// #include <net/if.h>
#include <unistd.h>
#include <linux/can.h>
#include <string.h>
#include <stdio.h>
#include <errno.h>

int main(int argc, char *argv[])
{

	int s;
	int status;

	s = socket(PF_CAN, SOCK_RAW, CAN_RAW);
	if (s < 0)
	{
		perror("socket");
		return 1;
	}

	struct sockaddr_can addr;
	memset(&addr, 0, sizeof(addr));

	struct ifreq ifr;

	strcpy(ifr.ifr_name, "vcan0");

	status = ioctl(s, SIOCGIFINDEX, &ifr);
	if (status < 0)
	{
		perror("ioctl");
		return 1;
	}

	addr.can_family = AF_CAN;

	addr.can_ifindex = ifr.ifr_ifindex;

	status = bind(s, (struct sockaddr *)&addr, sizeof(addr));
	if (status < 0)
	{
		perror("bind");
		return 1;
	}

	struct can_frame frame;

	for (;;)
	{
		int bytes = read(s, &frame, sizeof(struct can_frame));
		if (bytes <= 0)
		{
			perror("write");
			return 1;
		}

		printf("%02x %02x %02x %02x\n", frame.can_id = 0x101,
			   frame.can_dlc = 2,
			   frame.data[0] = 0x41,
			   frame.data[1] = 0x42);
	}

	close(s);

	return 0;
}

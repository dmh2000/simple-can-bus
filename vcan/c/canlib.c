#include <sys/socket.h>
#include <sys/ioctl.h>
#include <sys/select.h>
#include <unistd.h>
#include <net/if.h>
#include <linux/if.h>
#include <sys/ioctl.h>
#include <sys/socket.h>
#include <linux/can.h>
#include <linux/can/raw.h>
#include <string.h>
#include <errno.h>
#include "canlib.h"

int canlib_status(void)
{
	return errno;
}

/**
 * initialize the can bus and the client networks for sending
 */
int canlib_init(const char *can_dev)
{
	int sock;
	int status;

	sock = socket(PF_CAN, SOCK_RAW, CAN_RAW);
	if (sock < 0)
	{
		return -1;
	}

	struct ifreq ifr;
	strncpy(ifr.ifr_name, can_dev, IFNAMSIZ - 1);
	ifr.ifr_name[IFNAMSIZ - 1] = '\0';
	ifr.ifr_ifindex = if_nametoindex(can_dev);
	if (ifr.ifr_ifindex == 0)
	{
		return -2;
	}

	struct sockaddr_can addr;
	memset(&addr, 0, sizeof(addr));
	addr.can_family = AF_CAN;
	addr.can_ifindex = ifr.ifr_ifindex;
	status = bind(sock, (struct sockaddr *)&addr, sizeof(addr));
	if (status < 0)
	{
		return -3;
	}

	return sock;
}

// @return bytes read, 0 if timeout, < 0 if error
int canlib_receive(int can_sock, canlib_frame_t *can_frame, int timeout_ms)
{
	struct can_frame frame;
	fd_set readfds;
	struct timeval timeout;
	int status;
	int bytes;
	FD_ZERO(&readfds);
	timeout.tv_sec = timeout_ms / 1000;
	timeout.tv_usec = (timeout_ms % 1000) * 1000;

	FD_SET(can_sock, &readfds);
	status = select(can_sock + 1, &readfds, NULL, NULL, &timeout);
	if (status < 0)
	{
		// error
		return status;
	}
	if (status == 0)
	{
		// timeout
		return 0;
	}

	// should be read is available
	if (!FD_ISSET(can_sock, &readfds))
	{
		// error
		return -1;
	}

	bytes = read(can_sock, &frame, sizeof(struct can_frame));
	if (bytes <= 0)
	{
		return bytes;
	}

	// transfer data to canlib_frame_t
	can_frame->can_dlc = frame.can_dlc;
	can_frame->can_id = frame.can_id;
	memcpy(can_frame->data, frame.data, CAN_MAX_DATA_LEN);

	return bytes;
}

int canlib_send(int can_sock, canlib_frame_t *can_frame)
{
	struct can_frame frame;

	// transfer data from canlib_frame_t
	frame.can_id = can_frame->can_id;
	frame.can_dlc = (uint8_t)can_frame->can_dlc;
	memcpy(frame.data, can_frame->data, CAN_MAX_DATA_LEN);

	int bytes = write(can_sock, &frame, sizeof(struct can_frame));
	if (bytes < 0)
	{
		return bytes;
	}

	return bytes;
}

int canlib_close(int can_sock)
{
	int status = 0;

	if (can_sock > 0)
	{
		status = close(can_sock);
	}

	if (status < 0)
	{
		return errno;
	}
	return 0;
}

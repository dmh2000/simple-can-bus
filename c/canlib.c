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

/**
 * Get the last system error status
 * Returns the system errno value as uint32_t
 */
uint32_t canlib_status(void)
{
	return (uint32_t)errno;
}

/**
 * Initialize a CAN interface for communication
 * 
 * Creates and configures a socket for CAN communication using the specified interface.
 * 
 * @param can_dev Name of CAN interface (e.g. "can0")
 * @return Positive socket descriptor on success, negative error code on failure:
 *         CANLIB_ERR_PARAM if can_dev is NULL
 *         CANLIB_ERR_SOCKET if socket creation fails
 *         CANLIB_ERR_INTERFACE if interface lookup fails
 *         CANLIB_ERR_BIND if binding to interface fails
 */
int canlib_init(const char *can_dev)
{
	int sock;
	int status;

	if (!can_dev) {
		errno = EINVAL;
		return CANLIB_ERR_PARAM;
	}

	sock = socket(PF_CAN, SOCK_RAW, CAN_RAW);
	if (sock < 0) {
		return CANLIB_ERR_SOCKET;
	}

	struct ifreq ifr;
	strncpy(ifr.ifr_name, can_dev, IFNAMSIZ - 1);
	ifr.ifr_name[IFNAMSIZ - 1] = '\0';
	ifr.ifr_ifindex = if_nametoindex(can_dev);
	if (ifr.ifr_ifindex == 0) {
		close(sock);
		return CANLIB_ERR_INTERFACE;
	}

	struct sockaddr_can addr;
	memset(&addr, 0, sizeof(addr));
	addr.can_family = AF_CAN;
	addr.can_ifindex = ifr.ifr_ifindex;
	status = bind(sock, (struct sockaddr *)&addr, sizeof(addr));
	if (status < 0) {
		close(sock);
		return CANLIB_ERR_BIND;
	}

	return sock;
}

/**
 * Receive a CAN frame with timeout
 * 
 * Waits for data on the CAN bus with a specified timeout period.
 * Uses select() for timeout functionality.
 * 
 * @param can_sock Socket descriptor from canlib_init
 * @param can_frame Pointer to frame structure to store received data
 * @param timeout_ms Timeout in milliseconds
 * @return Number of bytes read on success, or error code:
 *         CANLIB_ERR_PARAM if invalid parameters
 *         CANLIB_ERR_IO if I/O error occurs
 *         CANLIB_ERR_TIMEOUT if operation times out
 */
int canlib_receive(int can_sock, canlib_frame_t *can_frame, int timeout_ms)
{
	if (can_frame == NULL || can_sock < 0) {
		errno = EINVAL;
		return CANLIB_ERR_PARAM;
	}

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
		return CANLIB_ERR_IO;
	}
	if (status == 0)
	{
		// timeout
		return CANLIB_ERR_TIMEOUT;
	}

	// should be read is available
	if (!FD_ISSET(can_sock, &readfds))
	{
		// error
		return CANLIB_ERR_IO;
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

/**
 * Send a CAN frame
 * 
 * Transmits a CAN frame on the specified socket.
 * Validates frame parameters before sending.
 * 
 * @param can_sock Socket descriptor from canlib_init
 * @param can_frame Pointer to frame structure containing data to send
 * @return Number of bytes sent on success, or error code:
 *         CANLIB_ERR_PARAM if invalid parameters
 *         CANLIB_ERR_IO if write operation fails
 */
int canlib_send(int can_sock, canlib_frame_t *can_frame)
{
	if (can_frame == NULL || can_sock < 0) {
		errno = EINVAL;
		return CANLIB_ERR_PARAM;
	}

	if (can_frame->can_dlc > CAN_MAX_DATA_LEN) {
		errno = EINVAL;
		return CANLIB_ERR_PARAM;
	}

	struct can_frame frame;

	// transfer data from canlib_frame_t
	frame.can_id = can_frame->can_id;
	frame.can_dlc = (uint8_t)can_frame->can_dlc;
	memcpy(frame.data, can_frame->data, CAN_MAX_DATA_LEN);

	int bytes = write(can_sock, &frame, sizeof(struct can_frame));
	if (bytes < 0)
	{
		return CANLIB_ERR_IO;
	}

	return bytes;
}

/**
 * Close a CAN socket
 * 
 * Safely closes an open CAN socket.
 * 
 * @param can_sock Socket descriptor to close
 * @return CANLIB_OK on success, CANLIB_ERR_IO if close fails
 */
int canlib_close(int can_sock)
{
	int status = 0;

	if (can_sock > 0)
	{
		status = close(can_sock);
	}

	if (status < 0)
	{
		return CANLIB_ERR_IO;
	}
	return CANLIB_OK;
}

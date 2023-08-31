# can_iot

sample canbus end-to-end

## install vcan and can-utils

### VCAN

https://www.pragmaticlinux.com/2021/10/how-to-create-a-virtual-can-interface-on-linux/

```bash
#!/bin/bash
# Load the kernel module.
sudo modprobe vcan
# Create the virtual CAN interface.
sudo ip link add dev vcan0 type vcan
# Bring the virtual CAN interface online.
sudo ip link set up vcan0
```

### can-utils

https://github.com/linux-can/can-utils

```bash

# INSTALL
sudo apt install can-utils

# TEST
# terminal 1
candump -tz vcan0

# terminal 2
cansend vcan0 123#00FFAA5501020304
```

Using candump or cansend is a good way to test the complement functions in the c or go versions.

## Directory 'vcan' : CAN bus access with C and/or Go

CAN bus data is usually accessed using struct can_frame, defined in /usr/include/linux/can.h. So to read and write to a can bus interface, use sizeof(can_frame) as the payload size.

## Build

## C

Directory 'c' contains a set of very simple functions that can be used to send and receive data from a CAN bus interface.

Test programs can_test_receive.c and can_test_send.c exercise the interface.

## G

Directory 'g' provides a Go module with functions that use the C libcan.so to interface to a CAN bus inteface. This directory also contains unit tests.

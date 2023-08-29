# can_iot
sample canbus end-to-end

## install vcan and can-utils

### VCAN
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
```bash

# INSTALL
sudo apt install can-utils

# TEST
# terminal 1
candump -tz vcan0

# terminal 2
cansend vcan0 123#00FFAA5501020304
```

### C Test Programs

./vcan/sender.c
./vcan/receiver.c




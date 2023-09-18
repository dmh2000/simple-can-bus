# SQIRVY

This project is a collection of programs that implement a simulation of a set of devices. It includes a CAN bus simulation app.
(a modbus sim is in development)

The architecture looks something like this:

<img src="./sqirvy.drawio.svg"/>

## CAN BUS SIMULATION

The CAN bus simulator uses the Linux 'vcan' support, along with a simple 'device' that is accessible over the CAN bus. Itlistens for CAN frames from a controller and responds with data.
There is a controller that consists of a web app front end, through an API endpoint that connects to the can bus and forwards commands and returns response data from the 'device'.

### How To Run

- Install vcan support on Linux if its not already there. I used Ubuntu 20 which has the support. The instructions here are for debian based systems.
- Install build-essential, clang and golang.
- activate the VCAN module for use as a network device (see below)
- clone the repo at https://github.com/dmh2000/sqirvy (pronounced 'scurvy')
- cd into top level
- execute 'make'
  - the make process will try to install a new shared library into /usr/local/lib

#### Activating the VCAN IP device

- https://www.pragmaticlinux.com/2021/10/how-to-create-a-virtual-can-interface-on-linux/

```bash
    #!/bin/bash
    # Load the kernel module.
    sudo modprobe vcan
    # Create the virtual CAN interface.
    sudo ip link add dev vcan0 type vcan
    # Bring the virtual CAN interface online.
    sudo ip link set up vcan0
```

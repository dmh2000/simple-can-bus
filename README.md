# SQIRVY

This project is a collection of Go programs that implement a simulation of a CAN device. It includes a CAN bus simulation app.

The architecture looks something like this:

<img src="./sqirvy.drawio.svg"/>

## TLDR

If you have a Linux system and you want to use CAN bus with Go, take a look in directory vcan/g. This module provides a simplified Go API for sending and receiving on a CAN bus. It uses the C library generated in vcan/c. These two directories are what you need to do basic CAN bus in the simplest way possible.

To see an example of a very simple full featured application, start in the next section.

### References

- https://github.com/linux-can/can-utils
- https://docs.kernel.org/networking/can.html
- https://elinux.org/CAN_Bus
- https://www.pragmaticlinux.com/2021/10/how-to-create-a-virtual-can-interface-on-linux/

## CAN BUS SIMULATION

The CAN bus simulator uses the Linux 'vcan' support, along with a simple 'device' that is accessible over the CAN bus. Itlistens for CAN frames from a controller and responds with data.
There is a controller that consists of a web app front end, through an API endpoint that connects to the can bus and forwards commands and returns response data from the 'device'.

The CAN bus messages include:

- 1 : ID_DIO_IN : the simulated device listens for this message to set a digital IO register.
- 2 : ID_DIO_OUT : the simulated device sends data from a digital IO register at 10 Hz. Applications can listen for this message to get updates.
- 3 : ID_DAC_IN : the simulated device listens for this message to set a digital-to-analog input.
- 4 : ID_ADC_OUT : the simulated devices sends data from an analog-to-digital device at 10 hz.

### How To Run

- Install vcan support on Linux if its not already there. I used Ubuntu 20 which has the support. The instructions here are for debian based systems.
- Install build-essential, golang (18 or later) and nodejs.
- Activate the VCAN module for use as a network device (see below)
- Clone this repo at https://github.com/dmh2000/sqirvy (pronounced 'scurvy')
- cd into top level
- execute 'make'
  - the make process will install a new shared library into /usr/local/lib. The makefile will ask for 'sudo' privileges.
- in a terminal, run the vcan/device/device program
- in a terminal, run the vcan/api/api program
- in a terminal, run the vcan/client/client program
- in a terminal, start the vcan/can-ui web client
  - this requires some node setup. See below for details.

### Activating the VCAN IP device

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

### Build

This projects uses simple make instructions to build the components.

The build process creates 4 executables, 1 Go support module and 1 shared library:

- vcan/c/libcan.so
- vcan/device/device
- vcan/client/client
- vcan/api/api

Note on make process :

- **vcan/c/libcan.so** : a C library that performs the low level access to the local can bus. In this case to the linux 'vcan' simulator.
  - the source files are in ./vcan/c
  - The build process also installs the updated library in /usr/local/lib. The process will ask for 'sudo' privileges.
  - this directory includes some simple tests that can be used with the can-utils package such as candump and cansend.
  - this library uses the IP model for the CAN bus
- **vcan/g/can.go**
  - a Go module that implements a bridge between the 'c' libcan.so and the Go programs.
  - it provides a simple Go api that higher level apps can use. The intent was to minimize the network details for sending and receiving CAN daa.
- **vcan/device/**

  - a Go program source and executable 'device', that acts as a very simple simulated CAN device.
  - it has two inputs, a DIO input and a DAC.
  - it has two ouputs, a DIO output and an ADC.

  MOVE CLIENT to API!!!!!

- **vcan/client**
  - a Go program that is the interface between the CAN bus and the REST api.
  - it provides functions to send app specific CAN messages and an active listener to data update.
- **vcan/api**
  - a Go program source that exposes a REST api accessible by any web application.
  - it exposes 3 urls:
    - /can/1
      - POSTs a json payload that sends the ID_DIO_IN message to the device.
    - /can/2
      - POSTs a json payloadthat sends the ID_DAC_IN message device.
    - /can/3
      - GETs the current value of all inputs and outputs in a json payload.

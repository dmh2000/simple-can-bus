# Can Client

This module implements the specific interface between a higher level client and the simulated CAN device. All communication takes place over the CAN bus.

## Devices

The interface includes the following devices:

- ADC : 32 bit analog to digital input
- DAC : 32 bit digital to analog output
- DIO OUT : 32 bit digital output
- DIO IN : 32 bit digital input

### Analog

The client sets the value of the DAC output on demand.
The device outputs the current value of the ADC at 10 hz.
The device DAC output value set by the client loops back to the ADC output value.

### Digital

The client sets the value of the DIO input on demand.
The device outputs the current value of the 32 bit DIO input.
The device DIO input value set by the client loops back to the DIO output value.

## Messages

| Message Id | Sender    | Receiver  | Description                 |
| ---------- | --------- | --------- | --------------------------- |
| 1          | simulator | any       | periodic value of DAC 10 Hz |
| 2          | client    | simulator | set DAC value               |
| 3          | simulator | any       | periodic value of DIO IN    |
| 4          | client    | simulator | set DIO OUT value           |

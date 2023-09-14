# Can Device Simulator

This directory contains an application that attempts to simulate a CAN device. The 'device' simulates some typical analog and digital devicees. It uses an ad-hoc message spec, not a higher level protocol such as Can Open.

## Devices

The simulation includes the following devices:

- ADC : 32 bit analog to digital input
- DAC : 32 bit digital to analog output
- DIO OUT : 32 bit digital output
- DIO IN : 32 bit digital input

### Analog

The simulator outputs the current value of the ADC at 10 hz.
The client sets the value of the DAC on demand.
The DAC output value set by the client loops back to the ADC input value.

### Digital

The simulator outputs the current value of the 32 bit DIO OUT.
The client sets the value of the DIO OUT on demand.
The DIO OUT value set by the client loops back to the DIO IN value.

## Messages

| Message Id | Sender    | Receiver  | Description                 |
| ---------- | --------- | --------- | --------------------------- |
| 1          | simulator | any       | periodic value of DAC 10 Hz |
| 2          | client    | simulator | set DAC value               |
| 3          | simulator | any       | periodic value of DIO IN    |
| 4          | client    | simulator | set DIO OUT value           |

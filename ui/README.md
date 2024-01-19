# WEB UI FOR SIMPLE CAN

## main.go

- the main function for a simple Go web server
- creates a main page handler and two fragment handlers (command and telemetry)
- web page URL : http://localhost:8000

## command.go

- creates a Handler func that displays the telemetry HTML fragment.
- provides 2 text fields and a submit button for seting the DIO and DAC values int he CAN device

## telemetry.go

- creates a Handler func that renders thethe state of the Telemetry HTML fragment.
- provides 4 text fields that display the current state of the CAN input settings and output data fields

all: 
	$(MAKE) -C vcan
	$(MAKE) -C modbus

clean:
	$(MAKE) -C vcan clean
	$(MAKE) -C modbus clean

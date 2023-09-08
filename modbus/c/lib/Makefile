all: modbus.c
	gcc -g -c modbus.c -lpthread
	ar rcs libmodbus.a modbus.o

clean:
	rm -rf modbus.o libmodbus.a	
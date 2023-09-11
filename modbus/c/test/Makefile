all :
	gcc -g -o adu_test adu_test.c -lpthread ../lib/libmodbus.a
	gcc -g -o client client.c -lpthread ../lib/libmodbus.a
	gcc -g -o server server.c -lpthread ../lib/libmodbus.a

clean:
	rm -rf client server
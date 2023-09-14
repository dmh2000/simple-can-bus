C=clang

all: libcan.so canlib_test_receive canlib_test_send

libcan.so : canlib.c canlib.h
	$(C) -g -Wall -o libcan.so -shared -fPIC canlib.c

canlib_test_receive : canlib_test_receive.c libcan.so
	$(C) -g -Wall -o canlib_test_receive canlib_test_receive.c -lcan

canlib_test_send: canlib_test_send.c libcan.so
	$(C) -g -Wall -o canlib_test_send canlib_test_send.c -lcan
	
clean:
	rm -f canlib_test_receive canlib_test_send libcan.so
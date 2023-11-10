export C=gcc

all: libcan.so canlib_test_receive canlib_test_send

libcan.so : canlib.c canlib.h
	@echo "------------------------------------------------------"
	@echo Building libcan.so
	@echo "------------------------------------------------------"	
	$(C) -g -Wall -o libcan.so -shared -fPIC canlib.c
	@echo "------------------------------------------------------"
	# @echo "Installing libcan.so to /usr/local/lib, requires sudo"
	# @echo "------------------------------------------------------"
	# sudo cp libcan.so /usr/local/lib
	# sudo ldconfig

canlib_test_receive : canlib_test_receive.c libcan.so
	$(C) -g -Wall -o canlib_test_receive -L. canlib_test_receive.c -lcan

canlib_test_send: canlib_test_send.c libcan.so
	$(C) -g -Wall -o canlib_test_send -L. canlib_test_send.c -lcan
	
clean:
	rm -f canlib_test_receive canlib_test_send libcan.so

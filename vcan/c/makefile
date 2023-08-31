all:
	clang -g -Wall -o libcan.so -shared -fPIC canlib.c
	clang -g -Wall -o canlib_test_receive -L. canlib_test_receive.c -lcan
	clang -g -Wall -o canlib_test_send -L. canlib_test_send.c -lcan
	cp canlib.h ../g
	cp libcan.so ../g

clean:
	rm -f canlib_test_receive canlib_test_send canlib.so
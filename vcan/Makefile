all:
	$(MAKE) -C c
	cp c/canlib.h g
	cp c/libcan.so g
	$(MAKE) -C g test 

clean:
	$(MAKE) -C c clean
	$(MAKE) -C g clean	


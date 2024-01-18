all: staticcheck

test: 
	@echo "------------------------------------------------------"
	@echo Testing can.go -go wrapper over canlib.so
	@echo "------------------------------------------------------"	
	LD_LIBRARY_PATH=../c go test can_test.go

staticcheck:
	staticcheck *.go

clean:
	rm -f canlib.h libcan.so

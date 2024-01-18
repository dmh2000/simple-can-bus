all : staticcheck device.go
	@echo "------------------------------------------------------"
	@echo Building simulated device
	@echo "------------------------------------------------------"	
	go build -o device device.go 

staticcheck:
	staticcheck *.go

clean :
	rm -f device
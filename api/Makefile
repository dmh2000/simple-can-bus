all : staticcheck api.go client.go
	@echo "------------------------------------------------------"
	@echo Building web client backend API
	@echo "------------------------------------------------------"	
	go build 

staticcheck:
	staticcheck *.go

clean:
	rm -f api
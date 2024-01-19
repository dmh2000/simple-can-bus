all : staticcheck main.go
	@echo "------------------------------------------------------"
	@echo Building Web UI
	@echo "------------------------------------------------------"	
	go build -o can-ui *.go

staticcheck:
	staticcheck *.go

clean:
	rm -f can-ui
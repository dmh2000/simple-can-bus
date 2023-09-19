export C=gcc

all: 
	@$(MAKE) -C vcan

clean:
	@$(MAKE) -C vcan clean

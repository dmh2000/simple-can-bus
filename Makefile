
all:
	@$(MAKE) -C c
	@$(MAKE) -C g
	@$(MAKE) -C device
	@$(MAKE) -C api
	@$(MAKE) -C ui

clean:
	@$(MAKE) -C c clean
	@$(MAKE) -C g clean	
	@$(MAKE) -C device clean
	@$(MAKE) -C api clean
	@$(MAKE) -C ui clean


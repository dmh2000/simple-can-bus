
all:
	@$(MAKE) -C c
	@$(MAKE) -C g test 
	@$(MAKE) -C device
	# @$(MAKE) -C client
	@$(MAKE) -C api

clean:
	@$(MAKE) -C c clean
	@$(MAKE) -C g clean	
	@$(MAKE) -C device clean
	# @$(MAKE) -C client clean
	@$(MAKE) -C api clean


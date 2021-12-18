SHELL=/bin/bash

EXE = jsgo

all: $(EXE)

jsgo:
	@echo "building $@ ..."
	$(MAKE) -s -f make.inc s=static

clean:
	rm -f $(EXE)


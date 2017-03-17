GO	:= $(shell which go 2>/dev/null)

all: hmc

hmc: hexmontocarlo.go
ifeq (${GO},"")
	$(error Could not find go in your $$PATH)
else
	go build -o $@ $^
endif

clean:
	rm -f hmc

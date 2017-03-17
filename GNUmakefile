GO	:= $(shell which go 2>/dev/null)

all: hmc

hmc: hexmontecarlo.go
ifeq (${GO},"")
	$(error Could not find go in your $$PATH)
else
	go build -o $@ $^
endif

clean:
	rm -f hmc

deps:
	${GO} get -v -u github.com/ericlagergren/go-prng/xorshift

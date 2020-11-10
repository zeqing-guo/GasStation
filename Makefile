unamestr := $(shell uname)

ifeq ($(unamestr),Linux)
	platform := linux
else
  ifeq ($(unamestr),Darwin)
	platform := darwin
  endif
endif

tags := $(platform)

GOBUILD := go build -tags "$(tags)"

bin/gas:
	$(GOBUILD) \
		-o bin/gas \
		github.com/zeqing-guo/GasStation/cmd/gas

all: bin/gas

clean:
	rm -rf bin/*

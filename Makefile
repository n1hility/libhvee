SRC = $(shell find . -type f -name '*.go')

.PHONY: default
default: build

.PHONY: build 
build: kvpctl.exe hyperv_kvp

kvpctl.exe: export GOOS=windows
kvpctl.exe: export GOARCH=amd64
kvpctl.exe: $(SRC) go.mod go.sum
	go build ./kvpctl

hyperv_kvp:
	go build
clean:
	rm *.exe hyperv_kvp

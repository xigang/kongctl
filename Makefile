PROJECT_BINARY_NAME := kongctl
PLATFORMS := linux windows darwin

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p bin && cd cmd && GOOS=$@ GOARCH=amd64 go build -o ../bin/$(PROJECT_BINARY_NAME)

.PHONY: bin
bin: linux windows darwin



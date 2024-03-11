GOCMD=go
GOBUILD=$(GOCMD) build

build-server:
	$(GOBUILD) -o ./bin/server ./cmd/server

build-client:
	$(GOBUILD) -o ./bin/client ./cmd/client

test:
	$(GOTEST) ./...

clean:
	$(GOCLEAN)
	rm -f ./bin/*


.PHONY: build-server build-client test clean

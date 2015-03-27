.PHONY: build clean api cli

all: clean build

build: api cli

api:
	go install github.com/oleksandr/idp/cmd/idp-api

cli:
	go install github.com/oleksandr/idp/cmd/idp-cli

clean:
	rm -rf $(GOPATH)/pkg/darwin_amd64/github.com/oleksandr/idp*

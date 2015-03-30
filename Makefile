.PHONY: build clean api cli thrift gox

all: clean build

build: api cli

api:
	go install github.com/oleksandr/idp/cmd/idp-api

cli:
	go install github.com/oleksandr/idp/cmd/idp-cli

thrift:
	mkdir -p $(PWD)/rpc/generated
	thrift -r --gen go:thrift_import="git-wip-us.apache.org/repos/asf/thrift.git/lib/go/thrift" -out $(PWD)/rpc/generated spec/services.thrift

gox:
	mkdir -p $(PWD)/build
	gox -osarch="linux/amd64" -osarch="windows/amd64" -osarch="darwin/amd64" -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" github.com/oleksandr/idp/cmd/idp-cli
	gox -osarch="linux/amd64" -osarch="windows/amd64" -osarch="darwin/amd64" -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" github.com/oleksandr/idp/cmd/idp-api

clean:
	rm -rf $(GOPATH)/pkg/darwin_amd64/github.com/oleksandr/idp*

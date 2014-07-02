NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: deps
	@mkdir -p bin/
	@echo "$(OK_COLOR)==> Building$(NO_COLOR)"
	@cd server; go build
	@mv server/server bin/image-server
	# @cd cli; go build
	# @mv cli/cli bin/image-cli
	@go test -v ./...

devserver:
	@go run `ls server/*.go | grep -v _test.go`

# devcli:
	# @go run `ls cli/*.go | grep -v _test.go`

deps:
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@go get -d -v ./...
	@echo $(DEPS) | xargs -n1 go get -d

updatedeps:
	@echo "$(OK_COLOR)==> Updating all dependencies$(NO_COLOR)"
	@go get -d -v -u ./...
	@echo $(DEPS) | xargs -n1 go get -d -u

clean:
	@rm -rf bin/

format:
	go fmt ./...

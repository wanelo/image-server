NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
MANTA_USER := $(shell echo $(MANTA_USER))
VERSION=1.0.6

all: deps
	@mkdir -p bin/
	@echo "$(OK_COLOR)==> Building$(NO_COLOR)"
	@go build -o bin/images
	# @cd cli; go build
	# @mv cli/cli bin/image-cli
	@go test -race -v ./...

devserver:
	@go run `ls server/*.go | grep -v _test.go`

# devcli:
	# @go run `ls cli/*.go | grep -v _test.go`

deps:
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@go get -d -v ./...
	@echo $(DEPS) | xargs -n1 go get -d
	@go get code.google.com/p/go.tools/cmd/godoc
	@go get code.google.com/p/go.tools/cmd/vet

updatedeps:
	@echo "$(OK_COLOR)==> Updating all dependencies$(NO_COLOR)"
	@go get -d -v -u ./...
	@echo $(DEPS) | xargs -n1 go get -d -u

clean:
	@rm -rf bin/

format:
	go fmt ./...

buildtomanta:
	@echo "$(OK_COLOR)==> Building for solaris amd64$(NO_COLOR)"
	@mkdir -p bin/solaris
	@GOOS=solaris GOARCH=amd64 go build -o bin/solaris/images
	# @echo "$(OK_COLOR)==> Compressing$(NO_COLOR)"
	# @cd bin/solaris && tar -czvf images.tar.gz images
	@echo "$(OK_COLOR)==> Uploading to manta$(NO_COLOR)"
	@mput -f bin/solaris/images /$(MANTA_USER)/public/images/bin/images-solaris-$(VERSION)

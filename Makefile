NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
MANTA_USER := $(shell echo $(MANTA_USER))
VERSION = $(shell cat core/version.go | grep 'const VERSION' | egrep -o '\d+\.\d+\.\d+')
AWS_ACCESS_KEY_ID := $(shell echo $(AWS_ACCESS_KEY_ID))
AWS_SECRET_KEY := $(shell echo $(AWS_SECRET_KEY))
AWS_BUCKET := $(shell echo $(AWS_BUCKET))
IMG_OUTPUTS := $(shell echo $(IMG_OUTPUTS))

all: deps tests
	@mkdir -p bin/
	@echo "$(OK_COLOR)==> Building$(NO_COLOR)"
	@go build -o bin/images
	# @cd cli; go build
	# @mv cli/cli bin/image-cli

devserver:
	@go run main.go --outputs $(IMG_OUTPUTS) --aws_access_key_id $(AWS_ACCESS_KEY_ID) --aws_secret_key $(AWS_SECRET_KEY) --aws_bucket $(AWS_BUCKET) --listen 127.0.0.1 --remote_base_path images server

# devcli:
	# @go run `ls cli/*.go | grep -v _test.go`

tests:
	@go test -race -v ./...

version:
	@echo $(VERSION)

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

build:
	@echo "$(OK_COLOR)==> Building for solaris amd64$(NO_COLOR)"
	@mkdir -p bin/solaris
	@GOOS=solaris GOARCH=amd64 go build -o bin/solaris/images

release: build
	# @echo "$(OK_COLOR)==> Compressing$(NO_COLOR)"
	# @cd bin/solaris && tar -czvf images.tar.gz images
	@echo "$(OK_COLOR)==> Uploading to manta$(NO_COLOR)"
	@mput -f bin/solaris/images /$(MANTA_USER)/public/images/bin/images-solaris-$(VERSION)
	@echo "$(VERSION)" | mput -H 'content-type: text/plain' /$(MANTA_USER)/public/images/bin/images-solaris-version

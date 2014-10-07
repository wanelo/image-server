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
IMG_REMOTE_BASE_PATH := $(shell echo $(IMG_REMOTE_BASE_PATH))
IMG_REMOTE_BASE_URL := $(shell echo $(IMG_REMOTE_BASE_URL))
HOME := $(shell echo $(HOME))
# GO = go
GO := $(shell echo $(HOME)/go.trunk/bin/go)

all: format deps tests

devserver:
	@$(GO) run main.go --outputs $(IMG_OUTPUTS) --aws_access_key_id $(AWS_ACCESS_KEY_ID) --aws_secret_key $(AWS_SECRET_KEY) --aws_bucket $(AWS_BUCKET) --listen 127.0.0.1 --remote_base_path $(IMG_REMOTE_BASE_PATH) --remote_base_url $(IMG_REMOTE_BASE_URL) server

tests:
	@$(GO) test -race -v ./...

version:
	@echo $(VERSION)

deps:
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@$(GO) get -d -v ./...
	@echo $(DEPS) | xargs -n1 go get -d
	@$(GO) get code.google.com/p/go.tools/cmd/godoc
	@$(GO) get code.google.com/p/go.tools/cmd/vet

updatedeps:
	@echo "$(OK_COLOR)==> Updating all dependencies$(NO_COLOR)"
	@$(GO) get -d -v -u ./...
	@echo $(DEPS) | xargs -n1 go get -d -u

clean:
	@rm -rf bin/
	@rm -fr tmp
	@rm -fr public

format:
	@go fmt ./...

build:
	@mkdir -p bin/
	@echo "$(OK_COLOR)==> Building$(NO_COLOR)"
	@rm -f bin/images bin/solaris/images
	@$(GO) build -o bin/images-$(VERSION)
	@cd bin && ln -s images-$(VERSION) images
	@echo "$(OK_COLOR)==> Building for solaris amd64$(NO_COLOR)"
	@mkdir -p bin/solaris
	@GOOS=solaris GOARCH=amd64 $(GO) build -o bin/solaris/images-$(VERSION)
	@cd bin/solaris && ln -s images-$(VERSION) images
	@echo "$(OK_COLOR)==> Compressing$(NO_COLOR)"
	@cd bin/solaris && tar -czvf images-$(VERSION).tar.gz images-$(VERSION)
	@echo "$(OK_COLOR)==> Build OK$(NO_COLOR)"

release: tests build
	@echo "$(OK_COLOR)==> Uploading to manta$(NO_COLOR)"
	@mput -f bin/solaris/images-$(VERSION) /$(MANTA_USER)/public/images/bin/images-solaris-$(VERSION)
	@echo "$(VERSION)" | mput -H 'content-type: text/plain' /$(MANTA_USER)/public/images/bin/images-solaris-version

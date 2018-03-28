VERSION ?= latest
NOROOT=-u $$(id -u):$$(id -g)
SRCDIR=/go/src/github.com/gravitational/stolon-app/tool/stolonctl
DOCKERFLAGS=--rm=true $(NOROOT) -v $(CURDIR):$(SRCDIR) -w $(SRCDIR)
BUILDIMAGE=golang:1.9

IMAGE_TAG := quay.io/gravitational/stolonctl:$(VERSION)

GOMETALINTER := $(GOPATH)/bin/gometalinter
GO_FILES = $(shell find . -type f -name "*.go")
BUILD_DIR := build

.PHONY: all
all: build

$(GOMETALINTER):
	go get -u gopkg.in/alecthomas/gometalinter.v2
	ln -s $(GOPATH)/bin/gometalinter.v2 $(GOPATH)/bin/gometalinter
	gometalinter --install

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter --vendor --disable-all \
		--enable=deadcode \
		--enable=ineffassign \
		--enable=gosimple \
		--enable=staticcheck \
		--enable=gofmt \
		--enable=goimports \
		--enable=dupl \
		--enable=misspell \
		--enable=errcheck \
		--enable=vet \
		--enable=vetshadow \
		--deadline=1m \
		./...

.PHONY: build
build: gofiles
	docker run $(DOCKERFLAGS) $(BUILDIMAGE) make stolonctl

.PHONY: gofiles
gofiles: $(GO_FILES)

.PHONY: stolonctl
stolonctl:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$@ -a -installsuffix cgo cmd/*.go

.PHONY: clean
clean:
	rm -rf build

.PHONY: images
images: build
	docker build --pull -t $(IMAGE_TAG) .

.PHONY: publish-images
publish-images: images
	docker push $(IMAGE_TAG)

.PHONY: revendor
revendor:
	dep ensure

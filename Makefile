export VERSION ?= $(shell ./version.sh)
REPOSITORY := gravitational.io
NAME := stolon-app
OPS_URL ?=

RUNTIME_VERSION := $(shell tele version | awk '/version:/ {print $$2}')

SRCDIR=/go/src/github.com/gravitational/stolon-app
DOCKERFLAGS=--rm=true -v $(PWD):$(SRCDIR) -w $(SRCDIR)
BUILDIMAGE=golang:1.11

EXTRA_GRAVITY_OPTIONS ?=

CONTAINERS := stolon-bootstrap:$(VERSION) \
			  stolon-uninstall:$(VERSION) \
			  stolon-hook:$(VERSION) \
			  stolon-jobs:$(VERSION) \
			  stolon:$(VERSION) \
			  stolon-telegraf:$(VERSION) \
			  stolonctl:$(VERSION)

IMPORT_IMAGE_OPTIONS := --set-image=stolon-bootstrap:$(VERSION) \
		--set-image=stolon-uninstall:$(VERSION) \
		--set-image=stolon-hook:$(VERSION) \
		--set-image=stolon-jobs:$(VERSION) \
		--set-image=stolon:$(VERSION) \
		--set-image=stolon-telegraf:$(VERSION) \
		--set-image=stolonctl:$(VERSION)

IMPORT_OPTIONS := --vendor \
		--ops-url=$(OPS_URL) \
		--insecure \
		--repository=$(REPOSITORY) \
		--name=$(NAME) \
		--version=$(VERSION) \
		--glob=**/*.yaml \
		--include="resources" \
		--include="registry" \
		--ignore="images" \
		--ignore="vendor/**/*.yaml" \
		--registry-url=leader.telekube.local:5000 \
		$(IMPORT_IMAGE_OPTIONS)

TELE_BUILD_OPTIONS := --insecure \
		--repository=$(OPS_URL) \
		--name=$(NAME) \
		--version=$(VERSION) \
		--glob=**/*.yaml \
		--ignore=".git" \
		--ignore="images" \
		--ignore="cmd" \
		--ignore="vendor/**/*.yaml" \
		$(IMPORT_IMAGE_OPTIONS)

BUILD_DIR := build

.PHONY: all
all: clean images

.PHONY: what-version
what-version:
	@echo $(VERSION)

.PHONY: images
images:
	cd images && $(MAKE) -f Makefile VERSION=$(VERSION)

.PHONY: import
import: images
	-gravity app delete --ops-url=$(OPS_URL) $(REPOSITORY)/$(NAME):$(VERSION) --force --insecure $(EXTRA_GRAVITY_OPTIONS)
	gravity app import $(IMPORT_OPTIONS) $(EXTRA_GRAVITY_OPTIONS) .

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.PHONY: build-app
build-app: $(BUILD_DIR) images
	sed -i.bak "s/version: \"0.0.0+latest\"/version: \"$(RUNTIME_VERSION)\"/" resources/app.yaml
	tele build -f -o $(BUILD_DIR)/installer.tar $(TELE_BUILD_OPTIONS) $(EXTRA_GRAVITY_OPTIONS) resources/app.yaml
	if [ -f resources/app.yaml.bak ]; then mv resources/app.yaml.bak resources/app.yaml; fi

.PHONY: build-stolonboot
build-stolonboot: $(BUILD_DIR)
	docker run $(DOCKERFLAGS) $(BUILDIMAGE) make build/stolonboot

build/stolonboot:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o $@ cmd/stolonboot/*.go

.PHONY: build-stolonctl
build-stolonctl: $(BUILD_DIR)
	docker run $(DOCKERFLAGS) $(BUILDIMAGE) make build/stolonctl

build/stolonctl:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o $@ cmd/stolonctl/*.go

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	cd images && $(MAKE) clean

$(GOMETALINTER):
	go get -u gopkg.in/alecthomas/gometalinter.v2
	ln -s $(GOPATH)/bin/gometalinter.v2 $(GOPATH)/bin/gometalinter
	gometalinter --install

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter --vendor --skip images/stolon/stolon --disable-all \
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

.PHONY: fix-logrus
fix-logrus:
	find vendor -type f -print0 | xargs -0 sed -i 's/Sirupsen/sirupsen/g'

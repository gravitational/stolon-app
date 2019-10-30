export VERSION ?= $(shell ./version.sh)
REPOSITORY := gravitational.io
NAME := stolon-app
OPS_URL ?= https://opscenter.localhost.localdomain:33009
TELE ?= $(shell which tele)
GRAVITY ?= $(shell which gravity)
RUNTIME_VERSION ?= $(shell $(TELE) version | awk '/^[vV]ersion:/ {print $$2}')
INTERMEDIATE_RUNTIME_VERSION ?= 5.2.15
GRAVITY_VERSION ?= 5.5.21
CLUSTER_SSL_APP_VERSION ?= "0.0.0+latest"

SRCDIR=/go/src/github.com/gravitational/stolon-app
DOCKERFLAGS=--rm=true -v $(PWD):$(SRCDIR) -w $(SRCDIR)
BUILDIMAGE=golang:1.11

EXTRA_GRAVITY_OPTIONS ?=

CONTAINERS := stolon-bootstrap:$(VERSION) \
			  stolon-uninstall:$(VERSION) \
			  stolon-hook:$(VERSION) \
			  stolon:$(VERSION) \
			  stolon-telegraf:$(VERSION) \
			  stolonctl:$(VERSION) \
			  stolon-pgbouncer:$(VERSION) \
			  stolon-etcd:$(VERSION)

IMPORT_IMAGE_OPTIONS := --set-image=stolon-bootstrap:$(VERSION) \
	--set-image=stolon-uninstall:$(VERSION) \
	--set-image=stolon-hook:$(VERSION) \
	--set-image=stolon:$(VERSION) \
	--set-image=stolon-telegraf:$(VERSION) \
	--set-image=stolonctl:$(VERSION) \
	--set-image=stolon-pgbouncer:$(VERSION) \
	--set-image=stolon-etcd:$(VERSION)

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
		--upgrade-via=$(INTERMEDIATE_RUNTIME_VERSION) \
		$(IMPORT_IMAGE_OPTIONS)

BUILD_DIR := build
BINARIES_DIR := bin

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(BINARIES_DIR):
	mkdir -p $(BINARIES_DIR)

.PHONY: all
all: clean images

.PHONY: what-version
what-version:
	@echo $(VERSION)

.PHONY: images
images:
	-git submodule update --init
	-git submodule update --remote
	cd images && $(MAKE) -f Makefile VERSION=$(VERSION)

.PHONY: import
import: images
	sed -i "s/version: \"0.0.0+latest\"/version: \"$(RUNTIME_VERSION)\"/" resources/app.yaml
	sed -i "s#gravitational.io/cluster-ssl-app:0.0.0+latest#gravitational.io/cluster-ssl-app:$(CLUSTER_SSL_APP_VERSION)#" resources/app.yaml
	-$(GRAVITY) app delete --ops-url=$(OPS_URL) $(REPOSITORY)/$(NAME):$(VERSION) --force --insecure $(EXTRA_GRAVITY_OPTIONS)
	$(GRAVITY) app import $(IMPORT_OPTIONS) $(EXTRA_GRAVITY_OPTIONS) .
	sed -i "s/version: \"$(RUNTIME_VERSION)\"/version: \"0.0.0+latest\"/" resources/app.yaml
	sed -i "s#gravitational.io/cluster-ssl-app:$(CLUSTER_SSL_APP_VERSION)#gravitational.io/cluster-ssl-app:0.0.0+latest#" resources/app.yaml

.PHONY: build-app
build-app: images
	sed -i "s/version: \"0.0.0+latest\"/version: \"$(RUNTIME_VERSION)\"/" resources/app.yaml
	sed -i "s#gravitational.io/cluster-ssl-app:0.0.0+latest#gravitational.io/cluster-ssl-app:$(CLUSTER_SSL_APP_VERSION)#" resources/app.yaml
	$(TELE) build -f -o $(BUILD_DIR)/installer.tar $(TELE_BUILD_OPTIONS) $(EXTRA_GRAVITY_OPTIONS) resources/app.yaml
	sed -i "s/version: \"$(RUNTIME_VERSION)\"/version: \"0.0.0+latest\"/" resources/app.yaml
	sed -i "s#gravitational.io/cluster-ssl-app:$(CLUSTER_SSL_APP_VERSION)#gravitational.io/cluster-ssl-app:0.0.0+latest#" resources/app.yaml

.PHONY: build-stolonboot
build-stolonboot: $(BUILD_DIR)
	docker run $(DOCKERFLAGS) $(BUILDIMAGE) make build-stolonboot-docker

build-stolonboot-docker:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o build/stolonboot cmd/stolonboot/*.go

.PHONY: build-stolonctl
build-stolonctl: $(BUILD_DIR)
	docker run $(DOCKERFLAGS) $(BUILDIMAGE) make build-stolonctl-docker

build-stolonctl-docker:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o build/stolonctl cmd/stolonctl/*.go

# number of environment variables are expected to be set
# see https://github.com/gravitational/robotest/blob/master/suite/README.md
#
.PHONY: robotest-run-suite
robotest-run-suite:
	./scripts/robotest_run_suite.sh $(shell pwd)/upgrade_from

.PHONY: download-binaries
download-binaries: $(BINARIES_DIR)
	for name in gravity tele; \
	do \
		curl https://get.gravitational.io/telekube/bin/$(GRAVITY_VERSION)/linux/x86_64/$$name -o $(BINARIES_DIR)/$$name; \
		chmod +x $(BINARIES_DIR)/$$name; \
	done

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	cd images && $(MAKE) clean
	rm -rf wd_suite

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

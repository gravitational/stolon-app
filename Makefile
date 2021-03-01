ifeq ($(origin VERSION), undefined)
# avoid ?= lazily evaluating version.sh (and thus rerunning the shell command several times)
VERSION := $(shell ./version.sh)
endif
REPOSITORY := gravitational.io
NAME := stolon-app
OPS_URL ?=
TELE ?= $(shell which tele)
GRAVITY ?= $(shell which gravity)
INTERMEDIATE_RUNTIME_VERSION ?=
GRAVITY_VERSION ?= 5.5.57
TELE_VERSION ?= $(GRAVITY_VERSION)
CLUSTER_SSL_APP_VERSION ?= 0.8.5
CLUSTER_SSL_APP_URL ?= https://github.com/gravitational/cluster-ssl-app/releases/download/${CLUSTER_SSL_APP_VERSION}/cluster-ssl-app-${CLUSTER_SSL_APP_VERSION}.tar.gz
STATEDIR ?= state

SRCDIR=/go/src/github.com/gravitational/stolon-app
DOCKERFLAGS=--rm=true -u $$(id -u):$$(id -g) -e XDG_CACHE_HOME=/tmp/.cache -v $(PWD):$(SRCDIR) -v $(GOPATH)/pkg:/gopath/pkg -w $(SRCDIR)
BUILDBOX=stolon-app-buildbox:latest

EXTRA_GRAVITY_OPTIONS ?=
TELE_BUILD_EXTRA_OPTIONS ?=

# --skip-version-check to build 5.5.x images with 7.0.x binary
TELE_BUILD_EXTRA_OPTIONS += --skip-version-check

# if variable is not empty add an extra parameter to tele build
ifneq ($(INTERMEDIATE_RUNTIME_VERSION),)
	TELE_BUILD_EXTRA_OPTIONS +=  --upgrade-via=$(INTERMEDIATE_RUNTIME_VERSION)
endif

# add state directory to the commands if STATEDIR variable not empty
ifneq ($(STATEDIR),)
	EXTRA_GRAVITY_OPTIONS +=  --state-dir=$(STATEDIR)
endif

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

ifneq ($(OPS_URL),)
	TELE_BUILD_EXTRA_OPTIONS +=  --repository=$(OPS_URL)
endif

TELE_BUILD_OPTIONS := --name=$(NAME) \
		--version=$(VERSION) \
		--glob=**/*.yaml \
		$(TELE_BUILD_EXTRA_OPTIONS) \
		$(IMPORT_IMAGE_OPTIONS)

TELE_BUILD_APP_OPTIONS := --insecure \
		--version=$(VERSION) \
		--set registry="" \
		--set image.tag=$(VERSION) \
		--set etcdImage.tag=$(VERSION) \
		--set telegrafImage.tag=$(VERSION) \
		--set pgbouncerImage.tag=$(VERSION) \
		--set stolonctlImage.tag=$(VERSION) \
		--values resources/custom-values.yaml

BUILD_DIR := build
BINARIES_DIR := bin

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(BINARIES_DIR):
	mkdir -p $(BINARIES_DIR)

$(STATEDIR):
	mkdir -p $(STATEDIR)

.PHONY: all
all: clean images

.PHONY: buildbox
buildbox:
	cd images && $(MAKE) -f Makefile buildbox

.PHONY: what-version
what-version:
	@echo $(VERSION)

.PHONY: images
images: lint
	cd images && $(MAKE) -f Makefile VERSION=$(VERSION)

.PHONY: import
import: images
	sed -i "s#gravitational.io/cluster-ssl-app:0.0.0+latest#gravitational.io/cluster-ssl-app:$(CLUSTER_SSL_APP_VERSION)#" resources/app.yaml
	sed -i "s/tag: latest/tag: $(VERSION)/g" resources/charts/stolon/values.yaml
	sed -i "s/0.1.0/$(VERSION)/g" resources/charts/stolon/Chart.yaml
	$(GRAVITY) app delete --ops-url=$(OPS_URL) $(REPOSITORY)/$(NAME):$(VERSION) --force $(EXTRA_GRAVITY_OPTIONS)
	$(GRAVITY) app import $(IMPORT_OPTIONS) $(EXTRA_GRAVITY_OPTIONS) .
	sed -i "s#gravitational.io/cluster-ssl-app:$(CLUSTER_SSL_APP_VERSION)#gravitational.io/cluster-ssl-app:0.0.0+latest#" resources/app.yaml
	sed -i "s/tag: $(VERSION)/tag: latest/g" resources/charts/stolon/values.yaml
	sed -i "s/$(VERSION)/0.1.0/g" resources/charts/stolon/Chart.yaml

# .PHONY because VERSION is dynamic
.PHONY: $(BUILD_DIR)/resources/app.yaml
$(BUILD_DIR)/resources/app.yaml: | $(BUILD_DIR)
	cp --archive resources $(BUILD_DIR)
	sed -i "s/version: \"0.0.0+latest\"/version: \"$(GRAVITY_VERSION)\"/" $(BUILD_DIR)/resources/app.yaml
	sed -i "s#gravitational.io/cluster-ssl-app:0.0.0+latest#gravitational.io/cluster-ssl-app:$(CLUSTER_SSL_APP_VERSION)#" $(BUILD_DIR)/resources/app.yaml
	sed -i "s/tag: latest/tag: $(VERSION)/g" $(BUILD_DIR)/resources/charts/stolon/values.yaml
	sed -i "s/0.1.0/$(VERSION)/g" $(BUILD_DIR)/resources/charts/stolon/Chart.yaml

.PHONY: build-app
build-app: images $(BUILD_DIR)/resources/app.yaml
	$(GRAVITY) $(EXTRA_GRAVITY_OPTIONS) package list
	$(TELE) build -f -o $(BUILD_DIR)/installer.tar $(TELE_BUILD_OPTIONS) $(EXTRA_GRAVITY_OPTIONS) $(BUILD_DIR)/resources/app.yaml

.PHONY: build-gravity-app
build-gravity-app: images
	sed -i "s/0.1.0/$(VERSION)/g" resources/charts/stolon/Chart.yaml
	$(TELE) build $(TELE_BUILD_APP_OPTIONS) -f -o $(BUILD_DIR)/application.tar resources/charts/stolon
	sed -i "s/$(VERSION)/0.1.0/g" resources/charts/stolon/Chart.yaml

.PHONY: build-stolonboot
build-stolonboot: $(BUILD_DIR)
	docker run $(DOCKERFLAGS) $(BUILDBOX) make build-stolonboot-docker

build-stolonboot-docker:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o build/stolonboot cmd/stolonboot/*.go

.PHONY: build-stolonctl
build-stolonctl: $(BUILD_DIR)
	docker run $(DOCKERFLAGS) $(BUILDBOX) make build-stolonctl-docker

build-stolonctl-docker:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o build/stolonctl cmd/stolonctl/*.go

# number of environment variables are expected to be set
# see https://github.com/gravitational/robotest/blob/master/suite/README.md
#
.PHONY: robotest-run-suite
robotest-run-suite:
	./robotest/run.sh pr

.PHONY: download-binaries
download-binaries: $(BINARIES_DIR)
	for name in gravity tele; \
	do \
		curl https://get.gravitational.io/telekube/bin/$(TELE_VERSION)/linux/x86_64/$$name -o $(BINARIES_DIR)/$$name; \
		chmod +x $(BINARIES_DIR)/$$name; \
	done

.PHONY: install-dependent-packages
install-dependent-packages: clean-state-dir $(STATEDIR) $(BUILD_DIR)
	$(TELE) pull gravity:$(GRAVITY_VERSION) $(EXTRA_GRAVITY_OPTIONS) -o $(BUILD_DIR)/gravity.tar --force
	tar xf $(BUILD_DIR)/gravity.tar -C $(STATEDIR) gravity.db packages
	curl -L $(CLUSTER_SSL_APP_URL) -o $(BUILD_DIR)/cluster-ssl-app.tar.gz
	$(GRAVITY) $(EXTRA_GRAVITY_OPTIONS) app import $(BUILD_DIR)/cluster-ssl-app.tar.gz

.PHONY: clean
clean: clean-state-dir
	-rm -rf $(BUILD_DIR)
	cd images && $(MAKE) clean
	-rm -rf wd_suite

clean-state-dir:
	-rm -rf $(STATEDIR)

.PHONY: fix-logrus
fix-logrus:
	find vendor -type f -print0 | xargs -0 sed -i 's/Sirupsen/sirupsen/g'

.PHONY: lint
lint: buildbox
	docker run $(DOCKERFLAGS) $(BUILDBOX) golangci-lint run --timeout=5m --skip-dirs=vendor ./...

.PHONY: push
push:
	$(TELE) push -f $(EXTRA_GRAVITY_OPTIONS) $(BUILD_DIR)/installer.tar

.PHONY: get-version
get-version:
	@echo $(VERSION)

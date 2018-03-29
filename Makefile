export VERSION ?= $(shell git describe --long --tags --always|awk -F'[.-]' '{print $$1 "." $$2 "." $$4}')
REPOSITORY := gravitational.io
NAME := stolon-app
OPS_URL ?= https://opscenter.localhost.localdomain:33009

SRCDIR=/go/src/github.com/gravitational/stolon-app
DOCKERFLAGS=--rm=true -v $(PWD):$(SRCDIR) -w $(SRCDIR)
BUILDIMAGE=golang:1.9.4

EXTRA_GRAVITY_OPTIONS ?=

CONTAINERS := stolon-bootstrap:$(VERSION) \
			  stolon-uninstall:$(VERSION) \
			  stolon-hook:$(VERSION) \
			  stolon-jobs:$(VERSION) \
			  stolon:$(VERSION) \
			  stolon-telegraf:$(VERSION) \
			  stolon-telegraf-node:$(VERSION)

IMPORT_IMAGE_OPTIONS := --set-image=stolon-bootstrap:$(VERSION) \
	--set-image=stolon-uninstall:$(VERSION) \
	--set-image=stolon-hook:$(VERSION) \
	--set-image=stolon-jobs:$(VERSION) \
	--set-image=stolon:$(VERSION) \
	--set-image=stolon-telegraf:$(VERSION) \
	--set-image=stolon-telegraf-node:$(VERSION)

IMPORT_OPTIONS := --vendor \
		--ops-url=$(OPS_URL) \
		--insecure \
		--repository=$(REPOSITORY) \
		--name=$(NAME) \
		--version=$(VERSION) \
		--glob=**/*.yaml \
		--exclude="build" \
		--exclude=".git" \
		--exclude="cmd" \
		--exclude="vendor" \
		--exclude="Makefile" \
		--exclude="images" \
		--exclude="gravity.log" \
		--ignore=images \
		--registry-url=apiserver:5000 \
		$(IMPORT_IMAGE_OPTIONS)

TELE_BUILD_OPTIONS := --insecure \
                --repository=$(OPS_URL) \
                --name=$(NAME) \
                --version=$(VERSION) \
                --glob=**/*.yaml \
                --ignore=".git" \
                --ignore="images" \
                --ignore="cmd" \
                --ignore="vendor" \
                $(IMPORT_IMAGE_OPTIONS)

BUILD_DIR := build
TARBALL := $(BUILD_DIR)/stolon-app.tar.gz

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

.PHONY: export
export: $(TARBALL)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(TARBALL): import $(BUILD_DIR)
	gravity package export $(REPOSITORY)/$(NAME):$(VERSION) $(TARBALL) $(EXTRA_GRAVITY_OPTIONS)

.PHONY: build-app
build-app: images
	tele build -o build/installer.tar $(TELE_BUILD_OPTIONS) $(EXTRA_GRAVITY_OPTIONS) resources/app.yaml

.PHONY: build-stolonboot
build-stolonboot:
	docker run $(DOCKERFLAGS) $(BUILDIMAGE) make build/stolonboot

build/stolonboot:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o $@ cmd/stolonboot/*.go

.PHONY: clean
clean:
	cd images && $(MAKE) clean

DB_NAME ?= postgres
.PHONY: dev-createdb
dev-createdb:
	-kubectl delete -f resources/createdb.yaml
	sed 's/{{STOLON_CREATE_DB}}/$(DB_NAME)/' resources/createdb.yaml | kubectl create -f -

.PHONY: dev-deletedb
dev-deletedb:
	-kubectl delete -f resources/deletedb.yaml
	sed 's/{{STOLON_DELETE_DB}}/$(DB_NAME)/' resources/deletedb.yaml | kubectl create -f -


.PHONY: dev-backup
dev-backup:
	-kubectl delete -f resources/backup.yaml
	sed 's/{{STOLON_BACKUP_DB}}/$(DB_NAME)/' resources/backup.yaml | kubectl create -f -

BACKUP_FILE ?=
.PHONY: dev-restore
dev-restore:
	-kubectl delete -f resources/restore.yaml
	sed 's/{{STOLON_BACKUP_FILE}}/\backups\/$(BACKUP_FILE)/' resources/restore.yaml | kubectl create -f -

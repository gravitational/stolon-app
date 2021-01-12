GOLANGCI_LINT_VER ?= 1.30.0

.PHONY: all
all: images

BOOTSTRAP_TAG := stolon-bootstrap:$(VERSION)
UNINSTALL_TAG := stolon-uninstall:$(VERSION)
HOOK_TAG := stolon-hook:$(VERSION)
STOLON_TAG := stolon:$(VERSION)
TELEGRAF_TAG := stolon-telegraf:$(VERSION)
STOLONCTL_TAG := stolonctl:$(VERSION)
PGBOUNCER_TAG := stolon-pgbouncer:$(VERSION)
ETCD_TAG := stolon-etcd:$(VERSION)

.PHONY: images
images:	buildbox \
		bootstrap \
		uninstall \
		stolon \
		telegraf \
		hook \
		stolonctl \
		pgbouncer \
		etcd
	docker tag $(BOOTSTRAP_TAG) stolon-bootstrap:latest
	docker tag $(UNINSTALL_TAG) stolon-uninstall:latest
	docker tag $(STOLON_TAG) stolon:latest
	docker tag $(HOOK_TAG) stolon-hook:latest
	docker tag $(TELEGRAF_TAG) stolon-telegraf:latest
	docker tag $(STOLONCTL_TAG) stolonctl:latest
	docker tag $(PGBOUNCER_TAG) stolon-pgbouncer:latest
	docker tag $(ETCD_TAG) stolon-etcd:latest


.PHONY: buildbox
buildbox:
	docker build --pull --build-arg GOLANGCI_LINT_VER=$(GOLANGCI_LINT_VER) -t stolon-app-buildbox:latest $(PWD)/buildbox

STOLONBOOT_BINS := bootstrap/bin
.PHONY: bootstrap
bootstrap:
	cd ../ && make build-stolonboot
	mkdir -p $(STOLONBOOT_BINS)
	cp ../build/stolonboot $(STOLONBOOT_BINS)/stolonboot
	docker build --pull -t $(BOOTSTRAP_TAG) $(PWD)/bootstrap

STOLONCTL_BINS := stolonctl/bin
.PHONY: stolonctl
stolonctl:
	cd ../ && make build-stolonctl
	mkdir -p $(STOLONCTL_BINS)
	cp ../build/stolonctl $(STOLONCTL_BINS)/stolonctl
	docker build --pull --build-arg CHANGESET=$(CHANGESET) --build-arg VERSION=$(VERSION) -t $(STOLONCTL_TAG) $(PWD)/stolonctl

IMAGES_FOR_BUILD := uninstall telegraf pgbouncer etcd
TAG_VARIABLE_NAME = $(shell echo '$@' | tr '[:lower:]' '[:upper:]')_TAG
.PHONY: $(IMAGES_FOR_BUILD)
$(IMAGES_FOR_BUILD):
	docker build --pull -t $($(TAG_VARIABLE_NAME)) $@

.PHONY: hook
hook:
	$(eval CHANGESET = $(shell echo '$(VERSION)' | sed -e 's/[\.]//g'))
	if [ -z "$(CHANGESET)" ]; then \
		echo "CHANGESET is not set"; exit 1; \
	fi;
	docker build --pull --build-arg CHANGESET=stolon-$(CHANGESET) -t $(HOOK_TAG) $(PWD)/hook

.PHONY: stolon
stolon:
	docker build --pull -t $(STOLON_TAG) $(PWD)/stolon

.PHONY: publish-stolonctl
publish-stolonctl: stolonctl
	docker push $(STOLONCTL_TAG)

.PHONY: clean
clean:
	-rm -rf $(STOLONBOOT_BINS)
	-rm -rf $(STOLONCTL_BINS)

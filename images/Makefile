.PHONY: all
all: images

BOOTSTRAP_TAG := stolon-bootstrap:$(VERSION)
UNINSTALL_TAG := stolon-uninstall:$(VERSION)
HOOK_TAG := stolon-hook:$(VERSION)
STOLON_TAG := stolon:$(VERSION)
TELEGRAF_TAG := stolon-telegraf:$(VERSION)
STOLONTOOL_TAG := stolontool:$(VERSION)
PGBOUNCER_TAG := stolon-pgbouncer:$(VERSION)
ETCD_TAG := stolon-etcd:$(VERSION)

.PHONY: images
images: bootstrap \
		uninstall \
		stolon \
		telegraf \
		hook \
		stolontool \
		pgbouncer \
		etcd
	docker tag $(BOOTSTRAP_TAG) stolon-bootstrap:latest
	docker tag $(UNINSTALL_TAG) stolon-uninstall:latest
	docker tag $(STOLON_TAG) stolon:latest
	docker tag $(HOOK_TAG) stolon-hook:latest
	docker tag $(TELEGRAF_TAG) stolon-telegraf:latest
	docker tag $(STOLONTOOL_TAG) stolontool:latest
	docker tag $(PGBOUNCER_TAG) stolon-pgbouncer:latest
	docker tag $(ETCD_TAG) stolon-etcd:latest

STOLONBOOT_BINS := bootstrap/bin
.PHONY: bootstrap
bootstrap:
	cd ../ && make build-stolonboot
	mkdir -p $(STOLONBOOT_BINS)
	cp ../build/stolonboot $(STOLONBOOT_BINS)/stolonboot
	docker build --pull -t $(BOOTSTRAP_TAG) $(PWD)/bootstrap

STOLONTOOL_BINS := stolontool/bin
.PHONY: stolontool
stolontool:
	cd ../ && make build-stolontool
	mkdir -p $(STOLONTOOL_BINS)
	cp ../build/stolontool $(STOLONTOOL_BINS)/stolontool
	docker build --pull --build-arg CHANGESET=$(CHANGESET) --build-arg VERSION=$(VERSION) -t $(STOLONTOOL_TAG) $(PWD)/stolontool

IMAGES_FOR_BUILD := uninstall telegraf pgbouncer etcd
TAG_VARIABLE_NAME = $(shell echo '$@' | tr '[:lower:]' '[:upper:]')_TAG
.PHONY: $(IMAGES_FOR_BUILD)
$(IMAGES_FOR_BUILD):
	docker build --pull -t $($(TAG_VARIABLE_NAME)) $@

.PHONY: hook
hook:
	$(eval CHANGESET = $(shell echo $$VERSION | sed -e 's/[\.]//g'))
	if [ -z "$(CHANGESET)" ]; then \
		echo "CHANGESET is not set"; exit 1; \
	fi;
	docker build --pull --build-arg CHANGESET=stolon-$(CHANGESET) -t $(HOOK_TAG) $(PWD)/hook

.PHONY: stolon
stolon:
	docker build --pull -t $(STOLON_TAG) $(PWD)/stolon

.PHONY: publish-stolontool
publish-stolontool: stolontool
	docker push $(STOLONTOOL_TAG)

.PHONY: clean
clean:
	-rm -rf $(STOLONBOOT_BINS)
	-rm -rf $(STOLONTOOL_BINS)

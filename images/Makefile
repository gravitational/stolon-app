.PHONY: all
all: images

BOOTSTRAP_TAG := stolon-bootstrap:$(VERSION)
UNINSTALL_TAG := stolon-uninstall:$(VERSION)
HOOK_TAG := stolon-hook:$(VERSION)
JOBS_TAG := stolon-jobs:$(VERSION)
STOLON_TAG := stolon:$(VERSION)
TELEGRAF_TAG := stolon-telegraf:$(VERSION)
TELEGRAF_NODE_TAG := stolon-telegraf-node:$(VERSION)

.PHONY: images
images: bootstrap \
		uninstall \
		stolon \
		telegraf \
		telegraf-node \
		hook \
		jobs
	docker tag $(BOOTSTRAP_TAG) stolon-bootstrap:latest
	docker tag $(UNINSTALL_TAG) stolon-uninstall:latest
	docker tag $(STOLON_TAG) stolon:latest
	docker tag $(HOOK_TAG) stolon-hook:latest
	docker tag $(JOBS_TAG) stolon-jobs:latest
	docker tag $(TELEGRAF_TAG) stolon-telegraf:latest
	docker tag $(TELEGRAF_NODE_TAG) stolon-telegraf-node:latest

STOLONBOOT_BINS := bootstrap/bin
$(STOLONBOOT_BINS):
	cd ../ && make build-stolonboot
	mkdir $(STOLONBOOT_BINS)
	cp ../build/stolonboot $(STOLONBOOT_BINS)/stolonboot

.PHONY: bootstrap
bootstrap: $(STOLONBOOT_BINS)
	docker build --pull -t $(BOOTSTRAP_TAG) $(PWD)/bootstrap

.PHONY: uninstall
uninstall:
	docker build --pull -t $(UNINSTALL_TAG) $(PWD)/uninstall

.PHONY: hook
hook:
	$(eval CHANGESET = $(shell echo $$VERSION | sed -e 's/[\.]//g'))
	if [ -z "$(CHANGESET)" ]; then \
		echo "CHANGESET is not set"; exit 1; \
	fi;
	docker build --pull --build-arg CHANGESET=stolon-$(CHANGESET) -t $(HOOK_TAG) $(PWD)/hook

STOLON_BINS := stolon/bin
$(STOLON_BINS):
	-git submodule update --init
	-git submodule update --remote
	cd stolon/stolon && make
	mkdir -p $(STOLON_BINS)
	cp stolon/stolon/bin/* $(STOLON_BINS)/

.PHONY: stolon
stolon: $(STOLON_BINS)
	docker build --pull -t $(STOLON_TAG) $(PWD)/stolon
	docker tag $(STOLON_TAG) stolon:latest

.PHONY: telegraf
telegraf:
	docker build --pull -t $(TELEGRAF_TAG) $(PWD)/telegraf

.PHONY: telegraf-node
telegraf-node:
	docker build --pull -t $(TELEGRAF_NODE_TAG) $(PWD)/telegraf-node

.PHONY: jobs
jobs:
	docker build --pull -t $(JOBS_TAG) $(PWD)/jobs

.PHONY: clean
clean:
	-rm -rf $(STOLONBOOT_BINS)
	-rm -rf $(STOLON_BINS)

.PHONY: all
all: images

BOOTSTRAP_TAG := stolon-bootstrap:$(VERSION)
UNINSTALL_TAG := stolon-uninstall:$(VERSION)
HOOK_TAG := stolon-hook:$(VERSION)
JOBS_TAG := stolon-jobs:$(VERSION)
STOLON_TAG := stolon:$(VERSION)
TELEGRAF_TAG := stolon-telegraf:$(VERSION)
STOLONCTL_TAG := stolonctl:$(VERSION)

.PHONY: images
images: bootstrap \
		uninstall \
		stolon \
		telegraf \
		hook \
		jobs \
		stolonctl
	docker tag $(BOOTSTRAP_TAG) stolon-bootstrap:latest
	docker tag $(UNINSTALL_TAG) stolon-uninstall:latest
	docker tag $(STOLON_TAG) stolon:latest
	docker tag $(HOOK_TAG) stolon-hook:latest
	docker tag $(JOBS_TAG) stolon-jobs:latest
	docker tag $(TELEGRAF_TAG) stolon-telegraf:latest
	docker tag $(STOLONCTL_TAG) stolonctl:latest

STOLONBOOT_BINS := bootstrap/bin
$(STOLONBOOT_BINS):
	cd ../ && make build-stolonboot
	mkdir $(STOLONBOOT_BINS)
	cp ../build/stolonboot $(STOLONBOOT_BINS)/stolonboot

.PHONY: bootstrap
bootstrap: $(STOLONBOOT_BINS)
	docker build --pull -t $(BOOTSTRAP_TAG) $(PWD)/bootstrap

STOLONCTL_BINS := stolonctl/bin
$(STOLONCTL_BINS):
	cd ../ && make build-stolonctl
	mkdir $(STOLONCTL_BINS)
	cp ../build/stolonctl $(STOLONCTL_BINS)/stolonctl

.PHONY: stolonctl
stolonctl: $(STOLONCTL_BINS)
	docker build --pull --build-arg CHANGESET=$(CHANGESET) --build-arg VERSION=$(VERSION) -t $(STOLONCTL_TAG) $(PWD)/stolonctl

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

.PHONY: jobs
jobs:
	docker build --pull -t $(JOBS_TAG) $(PWD)/jobs

.PHONY: publish-stolonctl
publish-stolonctl: stolonctl
	docker push $(STOLONCTL_TAG)

.PHONY: clean
clean:
	-rm -rf $(STOLONBOOT_BINS)
	-rm -rf $(STOLONCTL_BINS)
	-rm -rf $(STOLON_BINS)

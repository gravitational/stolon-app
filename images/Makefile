.PHONY: all
all: images

BOOTSTRAP_TAG := stolon-bootstrap:$(VERSION)
UNINSTALL_TAG := stolon-uninstall:$(VERSION)
HATEST_TAG := stolon-hatest:$(VERSION)
HOOK_TAG := stolon-hook:$(VERSION)
STOLON_TAG := stolon:$(VERSION)
UTILS_TAG := stolon-utils:$(VERSION)
PGWEB_TAG := pgweb:$(VERSION)

.PHONY: images
images: bootstrap \
		uninstall \
		stolon \
		utils \
		hook \
		pgweb \
		hatest
	docker tag $(BOOTSTRAP_TAG) stolon-bootstrap:latest
	docker tag $(UNINSTALL_TAG) stolon-uninstall:latest
	docker tag $(STOLON_TAG) stolon:latest
	docker tag $(HOOK_TAG) stolon-hook:latest
	docker tag $(HATEST_TAG) stolon-hatest:latest
	docker tag $(UTILS_TAG) stolon-utils:latest
	docker tag $(PGWEB_TAG) pgweb:latest

STOLONBOOT_BINS := bootstrap/bin
$(STOLONBOOT_BINS):
	cd ../tool/stolonboot && make clean build
	mkdir $(STOLONBOOT_BINS)
	cp ../tool/stolonboot/bin/stolonboot $(STOLONBOOT_BINS)/stolonboot

.PHONY: bootstrap
bootstrap: $(STOLONBOOT_BINS)
	docker build -t $(BOOTSTRAP_TAG) $(PWD)/bootstrap

.PHONY: uninstall
uninstall:
	docker build -t $(UNINSTALL_TAG) $(PWD)/uninstall

.PHONY: hook
hook:
	$(eval CHANGESET = $(shell echo $$VERSION | sed -e 's/[\.]//g'))
	if [ -z "$(CHANGESET)" ]; then \
		echo "CHANGESET is not set"; exit 1; \
	fi;
	docker build --build-arg CHANGESET=stolon-$(CHANGESET) --pull -t $(HOOK_TAG) $(PWD)/hook

STOLON_BINS := stolon/bin
$(STOLON_BINS):
	-git submodule update --init
	-git submodule update --remote
	cd stolon/stolon && make
	mkdir -p $(STOLON_BINS)
	cp stolon/stolon/bin/* $(STOLON_BINS)/

.PHONY: stolon
stolon: $(STOLON_BINS)
	docker build -t $(STOLON_TAG) $(PWD)/stolon
	docker tag $(STOLON_TAG) stolon:latest

STOLONHATEST_BINS := hatest/bin
$(STOLONHATEST_BINS):
	cd ../tool/stolonhatest && make clean build
	mkdir $(STOLONHATEST_BINS)
	cp ../tool/stolonhatest/bin/stolonhatest $(STOLONHATEST_BINS)/stolonhatest

.PHONY: hatest
hatest: $(STOLON_BINS) $(STOLONHATEST_BINS)
	docker build -t $(HATEST_TAG) $(PWD)/hatest

.PHONY: utils
utils:
	docker build -t $(UTILS_TAG) $(PWD)/stolon-utils

.PHONY: pgweb
pgweb:
	docker build -t $(PGWEB_TAG) $(PWD)/pgweb

.PHONY: clean
clean:
	-rm -rf $(STOLONBOOT_BINS)
	-rm -rf $(STOLON_BINS)
	-rm -rf $(STOLONHATEST_BINS)

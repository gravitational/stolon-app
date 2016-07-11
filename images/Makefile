.PHONY: all
all: images

BOOTSTRAP_TAG := stolon-bootstrap:0.0.1
UNINSTALL_TAG := stolon-uninstall:0.0.1
BACKUP_TAG := stolon-backup:0.0.1
HATEST_TAG := stolon-hatest:0.0.1
STOLON_TAG_VERSION := 0.2.0
STOLON_TAG := stolon:$(STOLON_TAG_VERSION)
STOLONCTL_TAG := stolonctl:$(STOLON_TAG_VERSION)

.PHONY: images
images: bootstrap \
		uninstall \
		stolon \
		backup \
		hatest

STOLONBOOT_BINS := bootstrap/bin
$(STOLONBOOT_BINS):
	cd ../tool/stolonboot && make clean build
	mkdir $(STOLONBOOT_BINS)
	cp ../tool/stolonboot/bin/stolonboot $(STOLONBOOT_BINS)/stolonboot

.PHONY: bootstrap
bootstrap: $(STOLONBOOT_BINS)
	docker build -t $(BOOTSTRAP_TAG) $(PWD)/bootstrap
	docker tag $(BOOTSTRAP_TAG) apiserver:5000/$(BOOTSTRAP_TAG)

.PHONY: uninstall
uninstall:
	docker build -t $(UNINSTALL_TAG) $(PWD)/uninstall
	docker tag $(UNINSTALL_TAG) apiserver:5000/$(UNINSTALL_TAG)


.PHONY: backup
backup: stolon
	docker build -t $(BACKUP_TAG) $(PWD)/backup
	docker tag $(BACKUP_TAG) apiserver:5000/$(BACKUP_TAG)

STOLON_BINS := stolon/bin
STOLONCTL_BINS := stolonctl/bin
BUILDER_DOCKER_FLAGS := --rm=true \
						-u $$(id -u):$$(id -g) \
						-v $(PWD)/stolon/stolon:/go/src/github.com/sorintlab/stolon \
						-w /go/src/github.com/sorintlab/stolon \
						-e GO15VENDOREXPERIMENT=1
BUILDER_IMAGE := golang:1.5.4

.PHONY: stolon
stolon: $(STOLON_BINS)
	docker build -t $(STOLON_TAG) $(PWD)/stolon
	docker tag $(STOLON_TAG) apiserver:5000/$(STOLON_TAG)


stolon/bin:
	-git submodule update --init
	-git submodule update --remote
	docker run $(BUILDER_DOCKER_FLAGS) $(BUILDER_IMAGE) ./build
	mkdir -p $(STOLON_BINS)
	cp stolon/stolon/bin/* $(STOLON_BINS)/
	mkdir -p $(STOLONCTL_BINS)
	mv stolon/bin/stolonctl $(STOLONCTL_BINS)/

.PHONY: stolonctl
stolonctl: stolon/bin
	docker build -t $(STOLONCTL_TAG) $(PWD)/stolonctl
	docker tag $(STOLONCTL_TAG) apiserver:5000/$(STOLONCTL_TAG)

STOLONHATEST_BINS := hatest/bin
$(STOLONHATEST_BINS):
	$(MAKE) -C ../tool/stolonhatest
	mkdir $(STOLONHATEST_BINS)
	cp ../tool/stolonhatest/bin/stolonhatest $(STOLONHATEST_BINS)/

.PHONY: hatest
hatest: stolonctl $(STOLONHATEST_BINS)
	docker build -t $(HATEST_TAG) $(PWD)/hatest
	docker tag $(HATEST_TAG) apiserver:5000/$(HATEST_TAG)

.PHONY: clean
clean:
	-rm -rf $(STOLONBOOT_BINS)
	-rm -rf $(STOLON_BINS)
	-rm -rf $(STOLONHATEST_BINS)

.PHONY: all
all: images

BOOTSTRAP_TAG := stolon-bootstrap:$(VERSION)
UNINSTALL_TAG := stolon-uninstall:$(VERSION)
HATEST_TAG := stolon-hatest:$(VERSION)
STOLON_TAG := stolon:$(VERSION)

.PHONY: images
images: bootstrap \
		uninstall \
		stolon \
		hatest
	docker tag $(BOOTSTRAP_TAG) stolon-bootstrap:latest
	docker tag $(UNINSTALL_TAG) stolon-uninstall:latest
	docker tag $(STOLON_TAG) stolon:latest
	docker tag $(HATEST_TAG) stolon-hatest:latest

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

STOLON_BINS := stolon/bin
BUILDER_DOCKER_FLAGS := --rm=true \
						-u $$(id -u):$$(id -g) \
						-v $(PWD)/stolon/stolon:/go/src/github.com/sorintlab/stolon \
						-w /go/src/github.com/sorintlab/stolon \
						-e GO15VENDOREXPERIMENT=1
BUILDER_IMAGE := golang:1.5.4

.PHONY: stolon
stolon: $(STOLON_BINS)
	docker build -t $(STOLON_TAG) $(PWD)/stolon

$(STOLON_BINS):
	#TODO: uncomment
	#-git submodule update --init
	#-git submodule update --remote
	docker run $(BUILDER_DOCKER_FLAGS) $(BUILDER_IMAGE) ./build
	mkdir -p $(STOLON_BINS)
	cp stolon/stolon/bin/* $(STOLON_BINS)/

STOLONHATEST_BINS := hatest/bin
$(STOLONHATEST_BINS):
	$(MAKE) -C ../tool/stolonhatest clean
	$(MAKE) -C ../tool/stolonhatest
	mkdir $(STOLONHATEST_BINS)
	cp ../tool/stolonhatest/bin/stolonhatest $(STOLONHATEST_BINS)/

.PHONY: hatest
hatest: $(STOLON_BINS) $(STOLONHATEST_BINS)
	docker build -t $(HATEST_TAG) $(PWD)/hatest

.PHONY: clean
clean:
	-rm -rf $(STOLONBOOT_BINS)
	-rm -rf $(STOLON_BINS)
	-rm -rf $(STOLONHATEST_BINS)

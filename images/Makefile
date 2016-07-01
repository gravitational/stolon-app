.PHONY: all
all: images

BOOTSTRAP_TAG:=stolon-bootstrap:0.0.1
UNINSTALL_TAG:=stolon-uninstall:0.0.1
BACKUP_TAG:=stolon-backup:0.0.1
STOLON_TAG:=stolon:0.2.0

.PHONY: images
images: bootstrap uninstall stolon backup

STOLONBOOT_BINS:=bootstrap/stolonboot
$(STOLONBOOT_BINS):
	cd ../tool/stolonboot && make clean build
	cp ../tool/stolonboot/bin/stolonboot bootstrap/stolonboot

.PHONY: bootstrap
bootstrap: $(STOLONBOOT_BINS)
	docker build -t $(BOOTSTRAP_TAG) $(PWD)/bootstrap
	docker tag $(BOOTSTRAP_TAG) apiserver:5000/$(BOOTSTRAP_TAG)

.PHONY: uninstall
uninstall:
	docker build -t $(UNINSTALL_TAG) $(PWD)/uninstall
	docker tag $(UNINSTALL_TAG) apiserver:5000/$(UNINSTALL_TAG)

.PHONY: clean
clean:
	-rm bootstrap/stolonboot
	-rm -rf stolon/bin


.PHONY: stolon
stolon: stolon/bin
	docker build -t $(STOLON_TAG) $(PWD)/stolon
	docker tag $(STOLON_TAG) apiserver:5000/$(STOLON_TAG)

STOLON_BINS:=build/stolonctl build/stolon-keeper build/stolon-proxy build/stolon-sentinel
NOROOT=-u $$(id -u):$$(id -g)
SRCDIR=/go/src/github.com/sorintlab/stolon
DOCKERFLAGS=--rm=true $(NOROOT) -v $(PWD)/stolon/stolon:$(SRCDIR) -w $(SRCDIR) -e GO15VENDOREXPERIMENT=1
BUILDIMAGE=golang:1.5.4

stolon/bin:
	-git submodule update --init
	-git submodule update --remote
	docker run $(DOCKERFLAGS) $(BUILDIMAGE) ./build
	mkdir -p stolon/bin
	cp stolon/stolon/bin/* stolon/bin/

.PHONY: backup
backup: stolon
	docker build -t $(BACKUP_TAG) $(PWD)/backup
	docker tag $(BACKUP_TAG) apiserver:5000/$(BACKUP_TAG)


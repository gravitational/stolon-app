.PHONY: all
all: images

.PHONY: images
images: bootstrap uninstall stolon

STOLONBOOT_BINS:=bootstrap/stolonboot

$(STOLONBOOT_BINS):
	cd ../tool/stolonboot && make clean build
	cp ../tool/stolonboot/bin/stolonboot bootstrap/stolonboot

.PHONY: bootstrap
bootstrap: $(STOLONBOOT_BINS)
	docker build -t stolon-bootstrap:0.0.1 $(PWD)/bootstrap

.PHONY: uninstall
uninstall:
	docker build -t stolon-uninstall:0.0.1 $(PWD)/uninstall

.PHONY: clean
clean:
	-rm bootstrap/stolonboot
	-rm -rf stolon/bin

STOLON_BINS:=build/stolonctl build/stolon-keeper build/stolon-proxy build/stolon-sentinel
NOROOT=-u $$(id -u):$$(id -g)
SRCDIR=/go/src/github.com/sorintlab/stolon
DOCKERFLAGS=--rm=true $(NOROOT) -v $(PWD)/stolon/stolon:$(SRCDIR) -w $(SRCDIR) -e GO15VENDOREXPERIMENT=1
BUILDIMAGE=golang:1.5.4

.PHONY: stolon
stolon: stolon/bin
	cd stolon && docker build -t stolon:0.2.0 .

stolon/bin:
	-git submodule init
	-git submodule update
	docker run $(DOCKERFLAGS) $(BUILDIMAGE) ./build
	mkdir -p stolon/bin
	cp stolon/stolon/bin/* stolon/bin/


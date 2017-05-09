.PHONY: all
all: bin/stolonboot

NOROOT=-u $$(id -u):$$(id -g)
SRCDIR=/go/src/github.com/gravitational/stolon-app
DOCKERFLAGS=--rm=true $(NOROOT) -v $(PWD):$(SRCDIR) -w $(SRCDIR) -e GO15VENDOREXPERIMENT=1
BUILDIMAGE=golang:1.7.1

.PHONY: build
build:
	docker run $(DOCKERFLAGS) $(BUILDIMAGE) make clean all

bin/stolonboot: *.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o $@ .

.PHONY: clean
clean:
	rm -rf bin

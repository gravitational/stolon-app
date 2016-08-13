all: bin/stolonhatest

NOROOT=-u $$(id -u):$$(id -g)
SRCDIR=/go/src/github.com/gravitational/stolon-app
DOCKERFLAGS=--rm=true $(NOROOT) -v $(PWD):$(SRCDIR) -w $(SRCDIR) -e GO15VENDOREXPERIMENT=1
BUILDIMAGE=golang:1.5.4

.PHONY: build
build:
	docker run $(DOCKERFLAGS) $(BUILDIMAGE) make clean all

bin/stolonhatest: *.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o $@ .

.PHONY: start
start: bin/stolonhatest
	goreman start

.PHONY: test
test:
	go vet .

.PHONY: clean
clean:
	rm -rf bin

.PHONY: install-tools
install-tools:
	go get -u github.com/mattn/goreman

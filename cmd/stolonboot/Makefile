.PHONY: all
all: bin/stolonboot

NOROOT=-u $$(id -u):$$(id -g)
SRCDIR=/go/src/github.com/gravitational/stolon-app
DOCKERFLAGS=--rm=true $(NOROOT) -v $(PWD):$(SRCDIR) -w $(SRCDIR)
BUILDIMAGE=golang:1.9.4

.PHONY: build
build:
	docker run $(DOCKERFLAGS) $(BUILDIMAGE) make clean all

bin/stolonboot:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o $@ cmd/stolonboot

.PHONY: clean
clean:
	rm -rf bin

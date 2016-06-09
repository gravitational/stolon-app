.PHONY: all clean

all: bin/stolonboot

bin/stolonboot: *.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o $@ .

clean:
	rm -rf bin

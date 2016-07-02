APPLICATION ?= $$(basename $(CURDIR))
BUILD_DIR ?= bin

$(BUILD_DIR):
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o $(BUILD_DIR)/$(APPLICATION) .

.PHONY: start
start: $(BUILD_DIR)
	goreman start

.PHONY: test
test:
	go vet .

.PHONY: clean
clean:
	-rm -r $(BUILD_DIR)

.PHONY: install-tools
install-tools:
	go get -u github.com/mattn/goreman

.PHONY: all images bootstrap uninstall clean stolonboot stolon

all: images

images: bootstrap uninstall stolon

bootstrap: stolonboot
	docker build -t stolon-bootstrap:0.0.1 $(PWD)/bootstrap

stolonboot:
	cd ../tool/stolonboot && make all
	cp ../tool/stolonboot/bin/stolonboot bootstrap/stolonboot

stolon:
	git submodule init && git submodule update
	cd $(PWD)/stolon/stolon && ./build && cp bin/* ../bin/
	docker build -t stolon:0.2.0 $(PWD)/stolon

uninstall:
	docker build -t stolon-uninstall:0.0.1 $(PWD)/uninstall

clean:
	rm bootstrap/stolonboot

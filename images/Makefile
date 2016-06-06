.PHONY: all images bootstrap uninstall clean stolonboot

all: images

images: bootstrap uninstall

bootstrap: stolonboot
	docker build -t stolon-bootstrap:0.0.1 $(PWD)/bootstrap

stolonboot:
	cd ../tool/stolonboot && make all
	cp ../tool/stolonboot/bin/stolonboot bootstrap/stolonboot

uninstall:
	docker build -t stolon-uninstall:0.0.1 $(PWD)/uninstall

clean:
	rm bootstrap/stolonboot

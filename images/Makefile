.PHONY: all images bootstrap uninstall

all: images

images: bootstrap uninstall

bootstrap:
	docker build -t stolon-bootstrap:0.0.1 $(PWD)/bootstrap

uninstall:
	docker build -t stolon-uninstall:0.0.1 $(PWD)/uninstall

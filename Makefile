VER := 0.0.1
PACKAGE := gravitational.io/stolon-app:$(VER)

.PHONY: all images dev-push dev-destroy dev-deploy import

all: images

images:
	cd images && $(MAKE) -f Makefile

dev-push: images
	docker tag stolon-bootstrap:$(VER) apiserver:5000/stolon-bootstrap:$(VER)
	docker push apiserver:5000/stolon-bootstrap:$(VER)
	docker tag stolon-uninstall:$(VER) apiserver:5000/stolon-uninstall:$(VER)
	docker push apiserver:5000/stolon-uninstall:$(VER)
	docker tag quay.io/coreos/etcd:v2.3.6 apiserver:5000/quay.io/coreos/etcd:v2.3.6
	docker push apiserver:5000/quay.io/coreos/etcd:v2.3.6
	docker tag sorintlab/stolon:master apiserver:5000/sorintlab/stolon:master
	docker push apiserver:5000/sorintlab/stolon:master

dev-redeploy: dev-clean dev-deploy

dev-deploy: dev-push
	kubectl create -f dev/bootstrap.yml

dev-clean:
	-kubectl delete pod/stolon-init
	-kubectl delete \
		-f images/bootstrap/resources/keeper.yml \
		-f images/bootstrap/resources/proxy.yml \
		-f images/bootstrap/resources/secret.yml \
		-f images/bootstrap/resources/sentinel.yml \
		-f images/bootstrap/resources/etcd.yml

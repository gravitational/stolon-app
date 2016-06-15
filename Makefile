VER := 0.0.2
PACKAGE := gravitational.io/stolon-app:$(VER)

.PHONY: all images dev-push dev-destroy dev-redeploy dev-deploy import

all: images

images:
	cd images && $(MAKE) -f Makefile

dev-push: images
	docker tag stolon-bootstrap:$(VER) apiserver:5000/stolon-bootstrap:$(VER)
	docker push apiserver:5000/stolon-bootstrap:$(VER)
	docker tag stolon-uninstall:$(VER) apiserver:5000/stolon-uninstall:$(VER)
	docker push apiserver:5000/stolon-uninstall:$(VER)
	docker pull quay.io/coreos/etcd:v2.3.6
	docker tag quay.io/coreos/etcd:v2.3.6 apiserver:5000/quay.io/coreos/etcd:v2.3.6
	docker push apiserver:5000/quay.io/coreos/etcd:v2.3.6
	docker tag stolon:0.2.0 apiserver:5000/stolon:0.2.0
	docker push apiserver:5000/stolon:0.2.0

dev-redeploy: dev-clean dev-deploy

dev-deploy: dev-push
	-kubectl label nodes --all stolon-keeper=stolon-keeper
	kubectl create -f dev/bootstrap.yml

dev-clean:
	-kubectl label nodes -l stolon-keeper=stolon-keeper stolon-keeper-
	-kubectl delete pod/stolon-init secret/stolon
	-kubectl delete \
		-f images/bootstrap/resources/keeper.yml \
		-f images/bootstrap/resources/proxy.yml \
		-f images/bootstrap/resources/sentinel.yml \
		-f images/bootstrap/resources/etcd.yml

import:
	gravity app import --vendor --registry-url=apiserver:5000 --state-dir=/var/lib/gravity/opscenter . $(PACKAGE)

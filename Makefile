VER := 0.0.1
PACKAGE := gravitational.io/stolon-app:$(VER)


.PHONY: all images dev-push dev-destroy dev-deploy import

all: images

images:
	cd images && $(MAKE) -f Makefile

dev-push:
	docker tag stolon-bootstrap:$(VER) apiserver:5000/stolon-bootstrap:$(VER)
	docker push apiserver:5000/stolon-bootstrap:$(VER)
	docker tag stolon-uninstall:$(VER) apiserver:5000/stolon-uninstall:$(VER)
	docker push apiserver:5000/stolon-uninstall:$(VER)
	docker tag quay.io/coreos/etcd:v2.3.6 apiserver:5000/quay.io/coreos/etcd:v2.3.6
	docker push apiserver:5000/quay.io/coreos/etcd:v2.3.6
	docker tag sorintlab/stolon:master apiserver:5000/sorintlab/stolon:master
	docker push apiserver:5000/sorintlab/stolon:master


create:
	kubectl create -f etcd.yaml
	kubectl create -f sentinel.yaml
	kubectl create -f secret.yaml
	kubectl create -f proxy.yaml
	kubectl create -f keeper.yaml


clean:
	kubectl delete -f keeper.yaml
	kubectl delete -f proxy.yaml
	kubectl delete -f secret.yaml
	kubectl delete -f sentinel.yaml
	kubectl delete -f etcd.yaml

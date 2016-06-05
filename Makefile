.PHONY: all create clean

all: create

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

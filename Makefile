VER := 0.0.6
REPOSITORY := gravitational.io
NAME := stolon-app

OPS_URL ?= https://opscenter.localhost.localdomain:33009

CONTAINERS := stolon-bootstrap:$(VER) \
			  stolon-uninstall:$(VER) \
			  stolon:$(VER) \
			  stolon-hatest:$(VER)

IMPORT_IMAGE_DEP_FLAGS := --set-image=stolon-bootstrap:$(VER) \
	--set-image=stolon-uninstall:$(VER) \
	--set-image=stolon:$(VER) \
	--set-image=stolon-hatest:$(VER) \
	--set-dep=gravitational.io/k8s-onprem:$$(gravity app list --ops-url=$(OPS_URL) --insecure | grep -m 1 k8s-onprem | awk '{print $$3}' | cut -d: -f2 | cut -d, -f1)

IMPORT_OPTIONS := --vendor \
		--ops-url=$(OPS_URL) \
		--insecure \
		--repository=$(REPOSITORY) \
		--name=$(NAME) \
		--version=$(VER) \
		--glob=**/*.yaml \
		--ignore=dev \
		--ignore=images \
		--registry-url=apiserver:5000 \
		$(IMPORT_IMAGE_DEP_FLAGS)

.PHONY: all
all: clean images

.PHONY: images
images:
	cd images && $(MAKE) -f Makefile VERSION=$(VER)

.PHONY: import
import: images
	-gravity app delete --ops-url=$(OPS_URL) $(REPOSITORY)/$(NAME):$(VER) --force --insecure
	gravity app import $(IMPORT_OPTIONS) .

.PHONY: clean
clean:
	cd images && $(MAKE) clean

.PHONY: dev-push
dev-push: images
	for container in $(CONTAINERS); do \
		docker tag $$container apiserver:5000/$$container ;\
		docker push apiserver:5000/$$container ;\
	done

.PHONY: dev-redeploy
dev-redeploy: dev-clean dev-deploy

.PHONY: dev-deploy
dev-deploy: dev-push
	-kubectl label nodes -l role=node stolon-keeper=yes
	kubectl create -f dev/bootstrap.yaml

.PHONY: dev-clean
dev-clean:
	-kubectl label nodes -l stolon-keeper=yes stolon-keeper-
	-kubectl delete pod/stolon-init secret/stolon
	-kubectl delete \
		-f resources/keeper.yaml \
		-f resources/proxy.yaml \
		-f resources/sentinel.yaml

BACKUP_DB ?= postgres
.PHONY: dev-backup
dev-backup:
	-kubectl delete -f resources/backup.yaml
	sed 's/{{STOLON_BACKUP_DB}}/$(BACKUP_DB)/' resources/backup.yaml | kubectl create -f -

BACKUP_FILE ?=
.PHONY: dev-restore
dev-restore:
	-kubectl delete -f resources/restore.yaml
	sed 's/{{STOLON_BACKUP_FILE}}/\backups\/$(BACKUP_FILE)/' resources/restore.yaml | kubectl create -f -

.PHONY: dev-hatest
dev-hatest:
	-kubectl delete -f resources/hatest.yaml
	kubectl create -f resources/hatest.yaml

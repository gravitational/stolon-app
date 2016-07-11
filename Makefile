VER := 0.0.4
PACKAGE := gravitational.io/stolon-app:$(VER)
CONTAINERS := stolon-bootstrap:0.0.1 \
			  stolon-uninstall:0.0.1 \
			  stolon:0.2.0 \
			  stolon-backup:0.0.1

OUT := build/stolon-app.tar.gz

OPSCENTER_WORK_DIR := /var/lib/gravity/opscenter

APISERVER_HOST := apiserver

.PHONY: all
all: $(OUT)

.PHONY: images
images:
	cd images && $(MAKE) -f Makefile

$(OUT): $(shell find resources -type f)
	$(MAKE) images

.PHONY: delete
delete:
	-gravity app delete $(PACKAGE) --state-dir=$(OPSCENTER_WORK_DIR) --force

.PHONY: import
import:
	gravity app import  --state-dir=$(OPSCENTER_WORK_DIR) --registry-url=$(APISERVER_HOST):5000 --glob=**/*.yaml --ignore=examples --vendor .

.PHONY: reimport
reimport: delete import

.PHONY: clean
clean:
	rm -rf $(OUT)
	cd images && $(MAKE) clean

.PHONY: dev-push
dev-push: images
	docker push apiserver:5000/stolon-bootstrap:0.0.1
	docker push apiserver:5000/stolon-uninstall:0.0.1
	docker push apiserver:5000/stolon-backup:0.0.1
	docker push apiserver:5000/stolon-hatest:0.0.1
	docker push apiserver:5000/stolon:0.2.0

.PHONY: dev-redeploy
dev-redeploy: dev-clean dev-deploy

.PHONY: dev-deploy
dev-deploy: dev-push
	-kubectl label nodes -l role=node stolon-keeper=stolon-keeper
	kubectl create -f dev/bootstrap.yaml

.PHONY: dev-clean
dev-clean:
	-kubectl label nodes -l stolon-keeper=stolon-keeper stolon-keeper-
	-kubectl delete pod/stolon-init secret/stolon
	-kubectl delete \
		-f resources/keeper.yaml \
		-f resources/proxy.yaml \
		-f resources/sentinel.yaml

BACKUP_DB ?=
.PHONY: dev-backup
dev-backup:
	-kubectl delete -f resources/backup.yaml
	sed 's/{{STOLON_BACKUP_DB}}/$(BACKUP_DB)/' resources/backup.yaml | kubectl create -f -

BACKUP_FILE ?=
.PHONY: dev-restore
dev-restore:
	-kubectl delete -f resources/restore.yaml
	sed 's/{{STOLON_BACKUP_DB}}/$(BACKUP_DB)/g;s/{{STOLON_BACKUP_FILE}}/$(BACKUP_FILE)/g' resources/restore.yaml | kubectl create -f -

.PHONY: dev-hatest
dev-hatest:
	-kubectl delete -f resources/hatest.yaml
	kubectl create -f resources/hatest.yaml

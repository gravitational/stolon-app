VER:=0.0.2
PACKAGE:=gravitational.io/stolon-app:$(VER)
CONTAINERS:=stolon-bootstrap:0.0.1 stolon-uninstall:0.0.1
OUT:=build/stolon-app.tar.gz
LOCAL_WORK_DIR:=/var/lib/gravity/opscenter

REGISTRYIMAGE:=registry:2.1.1
REGISTRYPORT:=5056
RUNNINGREG:=$$(docker ps -q --filter=ancestor=$(REGISTRYIMAGE))
REGADDRESS=127.0.0.1:$(REGISTRYPORT)

.PHONY: all
all: $(OUT)

.PHONY: images
images:
	cd images && $(MAKE) -f Makefile

$(OUT): $(shell find resources -type f)
	$(MAKE) images
	$(MAKE) start-registry
	$(MAKE) push-layers-to-registry
	$(MAKE) export-layers-from-registry
	$(MAKE) make-tarball
	$(MAKE) stop-registry

#
# reimports (delete+import) the application into the locally running portal, for development
#
reimport: $(OUT)
	-gravity app --state-dir=$(LOCAL_WORK_DIR) delete $(PACKAGE) --force
	gravity app --state-dir=$(LOCAL_WORK_DIR) import $(OUT) $(PACKAGE)

#
# starts the temporary docker registry
#
.PHONY: start-registry
start-registry:
	$(MAKE) stop-registry
	@if [ -z "$(RUNNINGREG)" ]; then \
		docker run -d -p $(REGISTRYPORT):5000 $(REGISTRYIMAGE) ;\
		echo "Started temporary Docker registry on port $(REGISTRYPORT)\n" ;\
		sleep 2 ;\
	else \
		echo "Temporary Docker registry is already listening on port $(REGISTRYPORT)\n" ;\
	fi

#
# pushes images from local Docker to temporary registry. THIS TAKES A LOT OF TIME.
#
.PHONY: push-layers-to-registry
push-layers-to-registry:
	for container in $(CONTAINERS); do \
		echo "docker tag $$container $(REGADDRESS)/$$container" ;\
		docker tag $$container $(REGADDRESS)/$$container ;\
		docker push $(REGADDRESS)/$$container ;\
		docker rmi $(REGADDRESS)/$$container ;\
		echo "\n" ;\
	done

#
# exports /var/lib/registry from temporary registry container into 'registry' folder
#
.PHONY: export-layers-from-registry
export-layers-from-registry:
	@echo "Copying layers from temporary registry into registry folder..."
	@rm -rf registry
	@docker cp $(RUNNINGREG):/var/lib/registry registry

#
# builds app tarball
#
.PHONY: make-tarball
make-tarball:
	@echo "Making a tarball..."
	mkdir -p build
	tar -cvzf $(OUT) resources registry
	@echo "done ---> $(OUT)"

#
# stops the temporary docker registry
#
.PHONY: stop-registry
stop-registry:
	@if [ ! -z "$(RUNNINGREG)" ]; then \
		container=$(RUNNINGREG) ;\
		docker stop $$container >/dev/null && docker rm -v $$container >/dev/null ;\
	else \
		echo registry is not running ;\
	fi

.PHONY: clean
clean:
	rm -rf $(OUT)
	cd images && $(MAKE) clean


.PHONY: dev-push
dev-push: images
	docker tag stolon-bootstrap:0.0.1 apiserver:5000/stolon-bootstrap:0.0.1
	docker push apiserver:5000/stolon-bootstrap:0.0.1
	docker tag stolon-uninstall:0.0.1 apiserver:5000/stolon-uninstall:0.0.1
	docker push apiserver:5000/stolon-uninstall:0.0.1
	docker tag stolon:0.2.0 apiserver:5000/stolon:0.2.0
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
		-f images/bootstrap/resources/keeper.yaml \
		-f images/bootstrap/resources/proxy.yaml \
		-f images/bootstrap/resources/sentinel.yaml

.PHONY: vendor-import
vendor-import:
	-gravity app --state-dir=$(LOCAL_WORK_DIR) delete $(PACKAGE) --force
	gravity app import --debug --vendor --glob=**/*.yaml --ignore=examples --registry-url=apiserver:5000 --state-dir=$(LOCAL_WORK_DIR) . $(PACKAGE)


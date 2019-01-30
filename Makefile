SCRIPTS_PATH      := scripts
TEST_SCRIPTS_PATH := test/scripts

.PHONY: all
all: install test_all

.PHONY: depend
depend: 
	$(SCRIPTS_PATH)/dependencies.sh -f

.PHONY: depend-noforce
depend-noforce: 
	@$(SCRIPTS_PATH)/dependencies.sh

.PHONY: install
install: depend-noforce
	@GO111MODULE=on go install ./...

.PHONY: test_all
all: dovetail-tests hyperledger-fabric-tests corda-tests

.PHONY: dovetail-tests
dovetail-tests:
	@$(TEST_SCRIPTS_PATH)/dovetail.sh

.PHONY: hyperledger-fabric-tests
hyperledger-fabric-tests: fab-network-up iou-tests fab-network-down
	@$(TEST_SCRIPTS_PATH)/hyperledger-fabric.sh

.PHONY: corda-tests
corda-tests:
	@$(TEST_SCRIPTS_PATH)/corda.sh

.PHONY: fabadmin-tests
fabadmin-tests:
	@$(TEST_SCRIPTS_PATH)/fabadmin.sh

.PHONY: iou-tests
iou-tests: fabadmin-tests
	@$(TEST_SCRIPTS_PATH)/iou.sh

.PHONY: fab-network-up
fab-network-up:
	@$(TEST_SCRIPTS_PATH)/start-fab-network.sh

.PHONY: fab-network-down
fab-network-down:
	@$(TEST_SCRIPTS_PATH)/stop-fab-network.sh

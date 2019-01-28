SCRIPTS_PATH      := scripts
TEST_SCRIPTS_PATH := test/scripts

.PHONY: all
all: install iou-tests fab-network-down

.PHONY: depend
depend: 
	$(SCRIPTS_PATH)/dependencies.sh -f

.PHONY: depend-noforce
depend-noforce: 
	@$(SCRIPTS_PATH)/dependencies.sh

.PHONY: install
install: depend-noforce
	@GO111MODULE=on go install ./...

.PHONY: fab-network-up
fab-network-up:
	@$(TEST_SCRIPTS_PATH)/start-fab-network.sh

.PHONY: fab-network-down
fab-network-down:
	@$(TEST_SCRIPTS_PATH)/stop-fab-network.sh

.PHONY: fabadmin-tests
fabadmin-tests: fab-network-up
	@$(TEST_SCRIPTS_PATH)/fabadmin.sh

.PHONY: iou-tests
iou-tests: fabadmin-tests
	@$(TEST_SCRIPTS_PATH)/iou.sh

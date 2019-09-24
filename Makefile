SCRIPTS_PATH      := scripts
TEST_SCRIPTS_PATH := test/scripts

.PHONY: all
all: install dovetail-tests

.PHONY: depend
depend: 
	$(SCRIPTS_PATH)/dependencies.sh -f

.PHONY: depend-noforce
depend-noforce: 
	@$(SCRIPTS_PATH)/dependencies.sh

.PHONY: install
install: depend-noforce
	GO111MODULE=on go build -o $GOPATH/bin github.com/TIBCOSoftware/dovetail-cli/cmd/dovetail

.PHONY: buildtype
buildtype: install
	@$(SCRIPTS_PATH)/buildtype.sh

.PHONY: test_all
test_all: dovetail-tests hyperledger-fabric-tests corda-tests

.PHONY: dovetail-tests
dovetail-tests:
	@$(TEST_SCRIPTS_PATH)/dovetail.sh

.PHONY: hyperledger-fabric-tests
hyperledger-fabric-tests:
	@$(TEST_SCRIPTS_PATH)/hyperledger-fabric.sh

.PHONY: corda-tests
corda-tests:
	@$(TEST_SCRIPTS_PATH)/corda.sh

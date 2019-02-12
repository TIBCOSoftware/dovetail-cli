#!/bin/bash

GOPATH="${GOPATH:-${HOME}/go}"
DEPEND_FLOGO_TAG=v0.5.7
DEPEND_FABRIC_REL=release-1.4

function installDependencies {
    # update hyperledger fabric
    echo "download hyperledger fabric ..."
    if [ -d "${GOPATH}/src/github.com/hyperledger/fabric" ]; then
        # fix go get -u if the branch is not default branch
        cd ${GOPATH}/src/github.com/hyperledger/fabric
        git branch --set-upstream-to=origin/${DEPEND_FABRIC_REL} ${DEPEND_FABRIC_REL}
        git pull
    fi
    go get -u -d github.com/hyperledger/fabric
    cd ${GOPATH}/src/github.com/hyperledger/fabric
    git fetch $(git remote -v | grep '(fetch)' | grep hyperledger | awk '{print $1}') ${DEPEND_FABRIC_REL}
    echo "download hyperledger fabric-sdk-go ..."
    go get -u -d github.com/hyperledger/fabric-sdk-go
    cd ${GOPATH}/src/github.com/hyperledger/fabric-sdk-go
    echo "install dependencies of hyperledger fabric-sdk-go ..."
    make depend-noforce

    # update golang packages required by dovetail-cli runtime
    go get -u github.com/julienschmidt/httprouter
    go get -u github.com/jteeuwen/go-bindata/...
    go get -u github.com/kardianos/govendor

    # update golang packages required by fabric-sdk-go test
    go get -u github.com/golang/protobuf
    go get -u google.golang.org/grpc
    go get -u github.com/cloudflare/cfssl/...
    go get -u github.com/spf13/viper
    go get -u github.com/spf13/cast
    go get -u github.com/stretchr/testify
    go get -u golang.org/x/crypto
    go get -u github.com/pkg/errors
    go get -u github.com/golang/mock
    go install github.com/golang/mock/mockgen
    go get -u github.com/Knetic/govaluate

    # other required golang packages if GO111MODULE=auto
    go get -u go.uber.org/zap/zapcore
    go get -u github.com/spf13/cobra
    go get -u github.com/sirupsen/logrus

    # update flogo to required tag version
    go get -u -d github.com/TIBCOSoftware/flogo-lib
    cd ${GOPATH}/src/github.com/TIBCOSoftware/flogo-lib
    git fetch $(git remote -v | grep '(fetch)' | grep TIBCOSoftware | awk '{print $1}') tags/${DEPEND_FLOGO_TAG}
    if [ -d "${GOPATH}/src/github.com/TIBCOSoftware/flogo-contrib" ]; then
        # fix go get -u if the branch is not default branch
        cd ${GOPATH}/src/github.com/TIBCOSoftware/flogo-contrib
        git branch --set-upstream-to=origin/master master
        git pull
    fi
    go get -u -d github.com/TIBCOSoftware/flogo-contrib
    cd ${GOPATH}/src/github.com/TIBCOSoftware/flogo-contrib
    git fetch $(git remote -v | grep '(fetch)' | grep TIBCOSoftware | awk '{print $1}') tags/${DEPEND_FLOGO_TAG}

    # update dovetail packages
    go get -u -d github.com/TIBCOSoftware/dovetail-contrib
    return 0
}

# isDependenciesInstalled checks that Go tools are installed and help the user if they are missing
function isDependenciesInstalled {
    declare -a msgs=()

    # Check that required Go libs have been installed
    [ -d "${GOPATH}/src/github.com/hyperledger/fabric" ] || msgs+=("hyperledger fabric is not installed (go get -u -d github.com/hyperledger/fabric)")
    [ -d "${GOPATH}/src/github.com/hyperledger/fabric-sdk-go" ] || msgs+=("hyperledger fabric-sdk-go is not installed (go get -u -d github.com/hyperledger/fabric-sdk-go)")
    [ -d "${GOPATH}/src/github.com/TIBCOSoftware/flogo-lib" ] || msgs+=("TIBCO flogo-lib is not installed (go get -u -d github.com/TIBCOSoftware/flogo-lib)")
    [ -d "${GOPATH}/src/github.com/TIBCOSoftware/flogo-contrib" ] || msgs+=("TIBCO flogo-contrib is not installed (go get -u -d github.com/TIBCOSoftware/flogo-contrib)")
    [ -d "${GOPATH}/src/github.com/TIBCOSoftware/dovetail-contrib" ] || msgs+=("TIBCO dovetail-contrib is not installed (go get -u -d github.com/TIBCOSoftware/dovetail-contrib)")

    if [ ${#msgs[@]} -gt 0 ]; then
        echo ${msgs[@]} | tr ' ' '\n'
        return 1
    fi
}

function isForceMode {
    if [ "${BASH_ARGV[0]}" != "-f" ]; then
        return 1
    fi
}

# run from the specified GOPATH
if ! isDependenciesInstalled || isForceMode; then
    if [ ! -d "${GOPATH}" ] ; then
        echo "[Error] ${GOPATH} not found"
        exit 1
    fi
    echo "${GOPATH} found ..."
    cd ${GOPATH}
    installDependencies
else
    echo "No need to install dependencies"
fi

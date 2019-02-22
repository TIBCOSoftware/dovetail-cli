#!/bin/bash

echo "creating temp directory ..."
temp_dir=$(mktemp -d)
echo "created temp directory $temp_dir"
if [ -n "$GOPATH" ]; then
    echo "Removing /var/lib/build_server/buildtypes/ ..."
    rm -rf /var/lib/build_server/buildtypes/
    mkdir -p /var/lib/build_server/buildtypes/dovetail/
    echo "Copying GOPATH content '${GOPATH}' to tempdir ..."
    cp -r ${GOPATH}/ $temp_dir
    cp Dockerfile $temp_dir
    cd $temp_dir
    echo "Building tci-buildtype-dovetail:${BID} ..."
    docker build -t tci-buildtype-dovetail .
    docker tag tci-buildtype-dovetail reldocker.tibco.com/tibcolabs/tci-buildtype-dovetail:${BID}
    docker push reldocker.tibco.com/tibcolabs/tci-buildtype-dovetail:${BID}
    docker images reldocker.tibco.com/tibcolabs/tci-buildtype-dovetail:${BID} --format "{{.Repository}}        {{.Tag}}        {{.ID}}        {{.CreatedSince}}        {{.Size}}" > ${WORKDIR}/images.txt

    echo "cleaning up..."
    rm -Rf ${temp_dir}
else
    echo "GOPATH is empty, please set it up to your go path directory"
fi
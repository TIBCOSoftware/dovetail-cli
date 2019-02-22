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
    echo "Building dovetail-buildtype:${BID} ..."
    docker build -t dovetail-buildtype .
    docker tag dovetail-buildtype reldocker.tibco.com/tibcolabs/dovetail-buildtype:${BID}
    docker push reldocker.tibco.com/tibcolabs/dovetail-buildtype:${BID}
    docker images reldocker.tibco.com/tibcolabs/dovetail-buildtype:${BID} --format "{{.Repository}}        {{.Tag}}        {{.ID}}        {{.CreatedSince}}        {{.Size}}" > ${WORKDIR}/images.txt

    echo "cleaning up..."
    rm -Rf ${temp_dir}
else
    echo "GOPATH is empty, please set it up to your go path directory"
fi
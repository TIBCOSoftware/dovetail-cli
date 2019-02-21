#!/bin/bash

echo "creating temp directory ..."
temp_dir=$(mktemp -d)
echo "created temp directory $temp_dir"
if [ -n "$GOPATH" ]; then
    cp -r ${GOPATH}/ $temp_dir
    cp Dockerfile $temp_dir
    cd $temp_dir
    docker build -t dovetail-buildtype .
    docker tag dovetail-buildtype reldocker.tibco.com/tibcolabs/dovetail-buildtype:$env.BID
    docker push reldocker.tibco.com/tibcolabs/dovetail-buildtype:$env.BID
    docker images reldocker.tibco.com/tibcolabs/dovetail-buildtype:$env.BID --format "{{.Repository}}        {{.Tag}}        {{.ID}}        {{.CreatedSince}}        {{.Size}}" > $env.WORKDIR/artifacts.txt

    echo "cleaning up..."
    rm -Rf ${temp_dir}
else
    echo "GOPATH is empty, please set it up to your go path directory"
fi
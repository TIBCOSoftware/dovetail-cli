#!/bin/bash
# cleanup docker containers of fabric-sdk network 

echo "shutdown fabric test network ..."
docker kill $(docker ps | egrep 'fabsdkgo|fabric-tools|fabric-ccenv' | awk '{print $1}')
docker rm $(docker ps -a | egrep 'fabsdkgo|fabric-tools|fabric-ccenv' | awk '{print $1}')
docker rmi $(docker images | grep fabsdkgo | awk '{print $3}')

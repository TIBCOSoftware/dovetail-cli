#!/bin/bash
# cleanup docker containers of fabric-sdk network 

echo "shutdown fabric test network ..."
docker kill $(docker ps | egrep 'fabsdkgo|fabric-tools|fabric-ccenv' | awk '{print $1}')
docker rm $(docker ps -a | egrep 'fabsdkgo|fabric-tools|fabric-ccenv' | awk '{print $1}')
docker images -a | grep "fabsdkgo" | awk '{print $3}' | xargs docker rmi

#!/bin/bash
# cleanup docker containers of fabric-sdk network 

echo "shutdown fabric test network ..."
if [ -n  "$(docker ps | egrep 'fabsdkgo|fabric-tools|fabric-ccenv' | awk '{print $1}')" ]; then
echo "killing containers ..."
    docker kill $(docker ps | egrep 'fabsdkgo|fabric-tools|fabric-ccenv' | awk '{print $1}')
fi
if [ -n  "$(docker ps -a | egrep 'fabsdkgo|fabric-tools|fabric-ccenv' | awk '{print $1}')" ]; then
echo "removing containers ..."
    docker rm $(docker ps -a | egrep 'fabsdkgo|fabric-tools|fabric-ccenv' | awk '{print $1}')
fi
if [ -n  "$(docker images | grep 'fabric' | awk '{print $3}')" ]; then
echo "removing images ..."
    docker rmi $(docker images | grep 'fabric' | awk '{print $3}')
fi

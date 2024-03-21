#!/bin/bash

IMAGE_NAME="alpine:edge"
CONTAINER_NAME="server"
CONTAINER_NETWORK="tp0_testing_net"
SERVER_IP=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $CONTAINER_NAME)
SERVER_PORT=12345
HEALTHCHECK_SIGNAL="healthcheck"
NO_COLOR='\033[0m'
RED_COLOR='\033[0;31m'
GREEN_COLOR='\033[0;32m'

echo "Sending signal to container named $CONTAINER_NAME"
SIGNAL_RECEIVED=$(docker run --network=$CONTAINER_NETWORK $IMAGE_NAME echo -n $HEALTHCHECK_SIGNAL | nc $SERVER_IP $SERVER_PORT)

if [ $SIGNAL_RECEIVED == $HEALTHCHECK_SIGNAL ]; then
    echo -e "${GREEN_COLOR}Container is OK${NO_COLOR}"
else
    echo -e "${RED_COLOR}Container is KO${NO_COLOR}"
fi

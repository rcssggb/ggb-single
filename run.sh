#!/usr/bin/env bash

HOST_OS=$(uname -s)

echo "Detected OS: ${HOST_OS}"

echo "Exposing X11"

if [ $HOST_OS == "Linux" ]; then
    xhost +local:root
elif [ $HOST_OS == "Darwin" ]; then
    xhost + 127.0.0.1
else
    echo "unknown host os"
    exit
fi

if [ $HOST_OS == "Linux" ]; then
    export DISPLAY=${DISPLAY}
elif [ $HOST_OS == "Darwin" ]; then
    export DISPLAY="host.docker.internal:0"
else
    echo "unknown host os"
    exit
fi

docker-compose up
docker-compose down

echo "Disabling X11 exposure"

if [ $HOST_OS == "Linux" ]; then
    xhost -local:root
elif [ $HOST_OS == "Darwin" ]; then
    xhost - 127.0.0.1
fi

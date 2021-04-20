#!/usr/bin/env bash

docker network create ggb-network
docker-compose up -d rcssserver single-agent
# docker-compose down


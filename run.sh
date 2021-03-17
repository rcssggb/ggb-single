#!/usr/bin/env bash

docker network create ggb-network
docker-compose up rcssserver single-agent
docker-compose down


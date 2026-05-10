#!/usr/bin/env bash

source ./.env

podman run --name go_basic_template_db -d \
    -p 5432:5432 \
    -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
    -e POSTGRES_USER=${POSTGRES_USER} \
    -e POSTGRES_DB=${POSTGRES_DB} \
    --volume ./.docker/sql:/docker-entrypoint-initdb.d/ \
    docker.io/library/postgres

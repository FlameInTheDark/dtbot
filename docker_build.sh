#!/usr/bin/env bash
go build
docker kill $(docker ps -a -q --filter="name=dtbot")
docker rmi $(docker images --format '{{.Repository}}:{{.Tag}}' | grep 'imagename')
docker build -t dtbot .
docker run -d --rm --restart ALWAYS -e BOT_TOKEN=$BOT_TOKEN -e MONGO_CONN=$MONGO_CONN --name dtbot dtbot:latest

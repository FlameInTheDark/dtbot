#!/usr/bin/env bash
git pull origin master
go build
docker kill $(docker ps -a -q --filter="name=dtbot")
docker rm $(docker ps -a -q --filter="name=dtbot")
docker rmi $(docker images --format '{{.Repository}}:{{.Tag}}' | grep 'imagename')
docker build -t dtbot .
docker run -d --restart always -e BOT_TOKEN=$BOT_TOKEN -e MONGO_CONN=$MONGO_CONN --name dtbot dtbot:latest

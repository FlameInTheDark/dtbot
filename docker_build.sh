go build
docker kill $(docker ps -a -q --filter="name=dtbot")
docker rmi $(docker images --format '{{.Repository}}:{{.Tag}}' | grep 'imagename')
docker build -t dtbot .
docker run -d --rm -e BOT_TOKEN=$BOT_TOKEN --name dtbot dtbot:latest

go build
docker build -t dtbot .
docker run -d --rm -e BOT_TOKEN=$BOT_TOKEN --name dtbot dtbot:latest
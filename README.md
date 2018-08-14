# Discord Tools Bot

## Build for docker

Easy way to build docker image for Alpine:
Compile app with command:

`GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dtalp`

Make (or use my) Dockerfile:

```
FROM alpine
COPY . .
RUN apk add ca-certificates
ENTRYPOINT ["./dtalp"]
```

Go inside app directory and build docker image

`docker build -t dtbot .`

Add environment variable `BOT_TOKEN` with token of discord bot.
And run container:

`docker run -d --rm -e BOT_TOKEN=$BOT_TOKEN --name dtbot dtbot:latest`
FROM alpine
COPY . .
RUN apk add ca-certificates
ENTRYPOINT ["./dtalp"]
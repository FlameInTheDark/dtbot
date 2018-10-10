FROM ubuntu:18.04
RUN apt-get update
RUN apt-get install -y ca-certificates ffmpeg
COPY . .
ENTRYPOINT ["./dtbot"]
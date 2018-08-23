FROM ubuntu:18.04
COPY . .
RUN apt-get update
RUN apt-get install -y ca-certificates ffmpeg 
ENTRYPOINT ["./dtbot"]
FROM ubuntu:18.04
RUN apt-get update
RUN apt-get install -y ca-certificates ffmpeg
RUN wget https://yt-dl.org/downloads/latest/youtube-dl
RUN chmod a+rx /usr/local/bin/youtube-dl
COPY . .
ENTRYPOINT ["./dtbot"]
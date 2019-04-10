FROM ubuntu:18.04
RUN apt-get update
RUN apt-get install -y wget ca-certificates ffmpeg python
RUN wget https://yt-dl.org/downloads/latest/youtube-dl
RUN chmod a+rx youtube-dl
ADD /usr/local/go/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip
COPY . .
ENTRYPOINT ["./dtbot"]

![Discord Tools Bot](/logo.png?raw=true "Discord Tools Bot")
# Discord Tools Bot

[![Go report](http://goreportcard.com/badge/FlameInTheDark/dtbot)](http://goreportcard.com/report/FlameInTheDark/dtbot)
[![Build Status](https://travis-ci.org/FlameInTheDark/dtbot.svg?branch=master)](https://travis-ci.org/FlameInTheDark/dtbot)
[![Scrutinizer Quality Score](https://img.shields.io/scrutinizer/g/FlameInTheDark/dtbot/master.svg)](https://scrutinizer-ci.com/g/FlameInTheDark/dtbot/)

## Used APIs

* OpenWeatherMap
* Yandex Translate
* Newsapi.org
* Geonames.org
* cbr-xml-daily.ru
* youtube-dl
* owfont

## Bot's features

* Shows weather
* Translate words
* Shows news
* Shows currency
* Plays music from Youtube and Soundcloud
* Plays music from online radiostations
* Makes polls

## How to use

Bot commands

To use the `!b` command you need to add a guild role named `bot.admin` and add it to you!

Command | Description
------- | -----------
`!v join` | Add bot into you voice channel
`!v leave` | Remove bot from voice channel
`!b clear [from_num]` | Remove bot's messages `!b clear` or `!b clear 3` removes all messages from 3rd message
`!b setconf [parameter] [value]` | Sets the bot configuration for your channel.
`!help` | Shows help
`!help [command]` | Detail help `!help !y`
`!play [youtube_url]` | Adds track (or playlist) in queue and start playing
`!y add [song]` | Adds song from youtube or soundcloud
`!y clear` | Removes all songs from queue
`!y play` | Starts playing queue
`!y stop` | Stops playing queue
`!y skip` | Skipping one song
`!y list` | List of songs in queue
`!r play [radio_station]` | Plays specified network radio station `!r play http://air2.radiorecord.ru:9003/rr_320`
`!r stop` | Stops radio
`!w [place]` | Shows the weather in a specified location `!w New York`
`!n [category]` | Displays news in the specified category `!n technology`
`!t [target_lang] [text]` | Translator `!t ru Hello world`
`!c` | Shows currencies (default from config)
`!c list` | Shows list of available currencies
`!c [currency]` | Shows specified currency `!c USD EUR`
`!p new [fields]` | Creates new poll `!p new field one \| field two \| field three`
`!p vote [field_num]` | Votes in poll
`!p end` | Ends poll and shows results

## Set config parameters

Parameter | Description
--------- | -----------
`general.language [string]` | Sets bot language
`general.timezone [num]` | Sets bot timezone
`general.nick [string]` | Sets bot nickname
`embed.color [hex color like #007700]` | Sets bot embed color
`news.country [string]` | Sets bot news country
`weather.city [string]` | Sets default city for weather

## Build for docker

Easy way to build docker image for Ubuntu:

Install MongoDB and set environment variable with mongo connection string  

`export MONGO_CONN=mongodb://user:password@some-host.com/dtbot`

Clone repository and move inside app directory.  
Compile app with command:

`go build`

Create `config.toml` file from sample `sample.config.toml`

```toml
[weather]
WeatherToken = "OpenWeatherMap API Token"
City = "Moscow" # Default forecast city

[news]
ApiKey = "Api key from Newsapi.org"
Country = "us" # News country
Articles = 5 # Count of articles per request

[translate]
ApiKey = "Yandex Translate API key"

[general]
GeonamesUsername = "Username from Geonames.org"
Language = "en" # Bot language
Timezone = 5 # UTC 0 + Timezone
Game = "Half-Life 3" # Shows "Play in Half-Life 3" status
EmbedColor = 4039680 # Default embed color (Hex color converted to int)
ServiceURL = "https://youtube.com" # Dont touch this value
MessagePool = 10 # Count of indexed messages for delete
DatabaseName = "dtbot" # Bot's database name in MongoDB 

[currency]
Default = ["USD", "EUR"] # Array of default currencies
```

Make (or use my) Dockerfile:

```Dockerfile
FROM ubuntu:18.04
RUN apt-get update
RUN apt-get install -y wget ca-certificates ffmpeg python
RUN wget https://yt-dl.org/downloads/latest/youtube-dl
RUN chmod a+rx youtube-dl
COPY . .
ENTRYPOINT ["./dtbot"]
```

Build docker image

`docker build -t dtbot .`

Add environment variable `BOT_TOKEN` with token of discord bot.  
And run container:

`docker run -d --restart always -e BOT_TOKEN=$BOT_TOKEN -e MONGO_CONN=$MONGO_CONN --name dtbot dtbot:latest`
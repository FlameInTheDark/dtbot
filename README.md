![Discord Tools Bot](https://github.com/FlameIntheDark/dtbot/blob/master/logo.png?raw=true "Discord Tools Bot")
# Discord Tools Bot

[![Go report](http://goreportcard.com/badge/FlameInTheDark/dtbot)](http://goreportcard.com/report/FlameInTheDark/dtbot)
[![Build Status](https://travis-ci.org/FlameInTheDark/dtbot.svg?branch=master)](https://travis-ci.org/FlameInTheDark/dtbot)
[![Scrutinizer Quality Score](https://img.shields.io/scrutinizer/g/FlameInTheDark/dtbot/master.svg)](https://scrutinizer-ci.com/g/FlameInTheDark/dtbot/)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FFlameInTheDark%2Fdtbot.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FFlameInTheDark%2Fdtbot?ref=badge_shield)
[![Discord Bots](https://discordbots.org/api/widget/status/424221765321883658.svg)](https://discordbots.org/bot/424221765321883658)

[![Musor stats on Discord Bot List](https://discordbotlist.com/bots/424221765321883658/widget)](https://discordbotlist.com/bots/424221765321883658)

## [Documentation](https://dtbot.realpha.ru)

## Used APIs and external software

* [Dark Sky](https://darksky.net/poweredby/)
* [Yandex Translate](https://tech.yandex.ru/translate/)
* [News API](https://newsapi.org)
* [Geonames](https://www.geonames.org)
* [cbr-xml-daily.ru](https://www.cbr-xml-daily.ru)
* [youtube-dl](https://ytdl-org.github.io/youtube-dl/index.html)
* [FFmpeg](https://ffmpeg.org)
* [Weather Icons](https://erikflowers.github.io/weather-icons/)
* [Sypex Geo](https://sypexgeo.net/)
* [Twitch API](https://dev.twitch.tv)

## Bot's features

* Shows weather
* Translate words
* Shows news
* Shows currency
* Makes polls
* Plays music from Youtube and Soundcloud
* Plays music from online radio stations
* Announcing if Twitch stream is started
* Greetings new users

## How to use

Bot commands

To use the `!b` or `!cron` commands you need to create a guild role named `bot.admin` and add it to you!

Command | Description
------- | -----------
`!v join` | Add bot into you voice channel
`!v leave` | Remove bot from voice channel
`!b clear [from_num]` | Remove bot's messages `!b clear` or `!b clear 3` removes all messages from 3rd message
`!b setconf [parameter] [value]` | Sets the bot configuration for your channel.
`!help` | Shows help
`!help [command]` | Detail help `!help y`
`!help bot.admin` | Shows help how get `bot.admin` role
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
`!m [map/sat] [location]` | Sends location image from yandex map `!m map New-York` or `!m sat New-York`
`!cron add [cron_time] [command]` | Creates cron job for command `!cron add 0 0 12 * * * !w Chelyabinsk` - everyday in 12:00 UTC 0 use command `!w`
`!cron list` | Shows cron jobs
`!cron remove [id]` | Removes cron job by ID `!cron remove 1`
`!geoip [some_ip_address]` | Shows geographic information about IP
`!twitch add [twitch_login]` | Adds streamer in announcer
`!twitch remove [twitch_login] [custom_announce_message]` | Removes streamer from announcer (custom message is optional)
`!greetings add [text]` | Adds new greetings message with specified text
`!greetings remove` | Removes greetings message
`!greetings test` | Sends you a greetings message

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
# Old weather API token (deprecated and unused now)
WeatherToken = "OpenWeatherMap API Token"
# Default forecast city
City = "Moscow"

[news]
ApiKey = "Api key from Newsapi.org"
Country = "us"
Articles = 5

[translate]
ApiKey = "Yandex Translate API key"

[general]
GeonamesUsername = "Username from Geonames.org"
Language = "en"
# UTC 0 + Timezone
Timezone = 5
Game = "Half-Life 3"
# Default embed color (Hex color converted to int)
EmbedColor = 4039680
ServiceURL = "https://youtube.com"
MessagePool = 10
DatabaseName = "dtbot"
GeocodingApiKey = "yandex_geocode_api_key"

[currency]
Default = ["USD", "EUR"]

[metrics]
# InfluxDB connection address
Address = "http://some_server.com:8086"
Database = "dtbot"
User = "user"
Password = "password"
# Discord Bot List
[dbl]
Token = "discordbots.org_bot_token"
# Twitch announcer
[twitch]
ClientID = "twitch_application_client_id"
# Weather API
[darksky]
Token = "darksky_api_token"
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
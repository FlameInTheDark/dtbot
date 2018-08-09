FROM golang:latest

COPY . .
RUN go get github.com/bwmarrin/discordgo
RUN go get github.com/BurntSushi/toml

CMD go build

CMD ./dtbot
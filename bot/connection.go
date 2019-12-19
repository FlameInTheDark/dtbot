package bot

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Connection : Voice connection struct
type Connection struct {
	voiceConnection *discordgo.VoiceConnection
	playing         bool
	quitChan chan struct{}{}
}

// NewConnection creates and returns new voice connection
func NewConnection(voiceConnection *discordgo.VoiceConnection) *Connection {
	connection := new(Connection)
	connection.voiceConnection = voiceConnection
	connection.playing = false
	connection.send = make(chan []int16, 2)
	quitChan = new(chan struct{}, 1)
	return connection
}

// Disconnect remove from voice channel and connection
func (c *Connection) Disconnect() {
	_ = c.voiceConnection.Disconnect()
}

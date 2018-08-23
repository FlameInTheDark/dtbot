package bot

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Connection : Voice connection struct
type Connection struct {
	voiceConnection *discordgo.VoiceConnection
	send            chan []int16
	lock            sync.Mutex
	sendpcm         bool
	stopRunning     bool
	playing         bool
}

// NewConnection creates and returns new voice connection
func NewConnection(voiceConnection *discordgo.VoiceConnection) *Connection {
	connection := new(Connection)
	connection.voiceConnection = voiceConnection
	connection.playing = false
	return connection
}

// Disconnect remove from voice channel and connection
func (c Connection) Disconnect() {
	c.voiceConnection.Disconnect()
}

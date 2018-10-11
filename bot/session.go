package bot

import (
	"github.com/bwmarrin/discordgo"
)

type (
	// Session : Session with radio player and voice connection struct
	Session struct {
		Player             RadioPlayer
		guildID, ChannelID string
		connection         *Connection
	}

	// SessionManager : Contains all sessions
	SessionManager struct {
		sessions map[string]*Session
	}

	// JoinProperties : Voice connection propperties struct
	JoinProperties struct {
		Muted    bool
		Deafened bool
	}
)

// Creates and returns new session
func newSession(newGuildID, newChannelID string, conn *Connection) *Session {
	session := &Session{
		guildID:    newGuildID,
		ChannelID:  newChannelID,
		connection: conn,
	}
	return session
}

// GetConnection returns vice connection struct
func (sess *Session) GetConnection() *Connection {
	return sess.connection
}

// Play starts to play radio
func (sess Session) Play(source string) error {
	return sess.connection.Play(source)
}

// Stop stops radio
func (sess *Session) Stop() {
	sess.connection.Stop()
}

// NewSessionManager creates and returns new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{make(map[string]*Session)}
}

// GetByGuild returns session by guild ID
func (manager SessionManager) GetByGuild(guildID string) *Session {
	for _, sess := range manager.sessions {
		if sess.guildID == guildID {
			return sess
		}
	}
	return nil
}

// GetByChannel returns session by channel ID
func (manager SessionManager) GetByChannel(channelID string) (*Session, bool) {
	sess, found := manager.sessions[channelID]
	return sess, found
}

// Join add bot to voice channel
func (manager *SessionManager) Join(discord *discordgo.Session, guildID, channelID string,
	properties JoinProperties) (*Session, error) {
	vc, err := discord.ChannelVoiceJoin(guildID, channelID, properties.Muted, properties.Deafened)
	if err != nil {
		return nil, err
	}
	sess := newSession(guildID, channelID, NewConnection(vc))
	manager.sessions[channelID] = sess
	return sess, nil
}

// Leave remove bot from voice channel
func (manager *SessionManager) Leave(discord *discordgo.Session, session Session) {
	session.connection.Stop()
	session.connection.Disconnect()
	delete(manager.sessions, session.ChannelID)
}

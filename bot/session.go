package bot

import (
	"github.com/bwmarrin/discordgo"
)

type (
	Session struct {
		Player             RadioPlayer
		guildId, ChannelId string
		connection         *Connection
	}

	SessionManager struct {
		sessions map[string]*Session
	}

	JoinProperties struct {
		Muted    bool
		Deafened bool
	}
)

func newSession(guildID, channelID string, conn *Connection) *Session {
	session := &Session{
		guildId:    guildID,
		ChannelId:  channelID,
		connection: conn,
	}
	return session
}

func (sess *Session) GetConnection() *Connection {
	return sess.connection
}

func (sess Session) Play(source string) error {
	return sess.connection.Play(source)
}

func (sess *Session) Stop() {
	sess.connection.Stop()
}

func NewSessionManager() *SessionManager {
	return &SessionManager{make(map[string]*Session)}
}

func (manager SessionManager) GetByGuild(guildId string) *Session {
	for _, sess := range manager.sessions {
		if sess.guildId == guildId {
			return sess
		}
	}
	return nil
}

func (manager SessionManager) GetByChannel(channelId string) (*Session, bool) {
	sess, found := manager.sessions[channelId]
	return sess, found
}

func (manager *SessionManager) Join(discord *discordgo.Session, guildId, channelId string,
	properties JoinProperties) (*Session, error) {
	vc, err := discord.ChannelVoiceJoin(guildId, channelId, properties.Muted, properties.Deafened)
	if err != nil {
		return nil, err
	}
	sess := newSession(guildId, channelId, NewConnection(vc))
	manager.sessions[channelId] = sess
	return sess, nil
}

func (manager *SessionManager) Leave(discord *discordgo.Session, session Session) {
	session.connection.Stop()
	session.connection.Disconnect()
	delete(manager.sessions, session.ChannelId)
}

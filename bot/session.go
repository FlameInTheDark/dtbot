package bot

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type (
	// Session structure with radio player and voice connection
	Session struct {
		Queue              *SongQueue
		Player             RadioPlayer
		guildID, ChannelID string
		connection         *Connection
		Volume             float32
	}

	// SessionManager contains all sessions
	SessionManager struct {
		sessions map[string]*Session
	}

	// JoinProperties voice connection properties struct
	JoinProperties struct {
		Muted    bool
		Deafened bool
	}
)

// Creates and returns new session
func newSession(newGuildID, newChannelID string, conn *Connection, volume float32) *Session {
	session := &Session{
		Queue:      newSongQueue(),
		guildID:    newGuildID,
		ChannelID:  newChannelID,
		connection: conn,
		Volume:     volume,
	}
	return session
}

// GetConnection returns vice connection struct
func (sess *Session) GetConnection() *Connection {
	return sess.connection
}

// Play starts to play the thing in source
func (sess *Session) Play(source string, volume float32) error {
	ffmpeg := exec.Command("ffmpeg", "-i", source, "-f", "s16le", "-reconnect", "1", "-reconnect_at_eof", "1", "-reconnect_streamed", "1", "-reconnect_delay_max", "2", "-filter:a", fmt.Sprintf("volume=%.3f", volume), "-ar", strconv.Itoa(FRAME_RATE), "-ac", strconv.Itoa(CHANNELS), "pipe:1")
	reader, err := ffmpeg.StdoutPipe()
	if err != nil {
		return err
	}
	err = ffmpeg.Start()
	if err != nil {
		return err
	}
	defer func() {
		reader.Close()
		fferr := ffmpeg.Process.Kill()
		if fferr != nil {
			fmt.Println("FFMPEG close err: ", fferr)
		}
	}()
	return sess.connection.EncodeOpusAndSend(reader)
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
func (manager *SessionManager) GetByGuild(guildID string) *Session {
	for _, sess := range manager.sessions {
		if sess.guildID == guildID {
			return sess
		}
	}
	return nil
}

// GetByChannel returns session by channel ID
func (manager *SessionManager) GetByChannel(channelID string) (*Session, bool) {
	sess, found := manager.sessions[channelID]
	return sess, found
}

func (manager *SessionManager) GetAll() map[string]*Session {
	return manager.sessions
}

func (s *Session) IsOk() bool {
	if s.connection.voiceConnection != nil {
		if s.connection.voiceConnection.Ready == true {
			return true
		}
	}
	return false
}

// Join add bot to voice channel
func (manager *SessionManager) Join(discord *discordgo.Session, guildID, channelID string,
	properties JoinProperties, volume float32) (*Session, error) {
	vc, err := discord.ChannelVoiceJoin(guildID, channelID, properties.Muted, properties.Deafened)
	if err != nil {
		log.Println("Voice connection error: ", err)
		if vc != nil {
			err := vc.Disconnect()
			if err != nil {
				log.Println(err)
			}
			vc.Close()
		}
		return nil, err
	}
	log.Println("Voice joined")
	if vc == nil {
		return nil, errors.New("no voice connection")
	}
	if vc.Ready != true {
		return nil, errors.New("voice connection not ready")
	}
	log.Println("Voice connection OK and Ready")
	sess := newSession(guildID, channelID, NewConnection(vc), volume)
	manager.sessions[channelID] = sess
	log.Println("Voice session created and saved")
	return sess, nil
}

// Leave remove bot from voice channel
func (manager *SessionManager) Leave(discord *discordgo.Session, session Session) {
	session.connection.Stop()
	session.connection.Disconnect()
	delete(manager.sessions, session.ChannelID)
}

// Count returns count of voice sessions
func (manager *SessionManager) Count() int {
	return len(manager.sessions)
}

func (manager *SessionManager) GetChannels() []string {
	var ids []string
	for _, s := range manager.sessions {
		ids = append(ids, s.ChannelID)
	}
	return ids
}

func (manager *SessionManager) GetGuilds() []string {
	var ids []string
	for _, s := range manager.sessions {
		ids = append(ids, s.guildID)
	}
	return ids
}

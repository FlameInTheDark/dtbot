package bot

import (
	"fmt"
	"os/exec"
	"strconv"
)

// Song contains information about song
type Song struct {
	Media    string
	Title    string
	Duration *string
	Id       string
}

// Ffmpeg returns ffmpeg executable command
func (song Song) Ffmpeg(volume float32) *exec.Cmd {
	return exec.Command("ffmpeg", "-i", song.Media, "-f", "s16le", "-reconnect", "1", "-reconnect_at_eof", "1", "-reconnect_streamed", "1", "-reconnect_delay_max", "2", "-filter:a", fmt.Sprintf("volume=%.3f", volume), "-ar", strconv.Itoa(FRAME_RATE), "-ac",
		strconv.Itoa(CHANNELS), "pipe:1")
}

// NewSong creates and returns new song
func NewSong(media, title, id string) *Song {
	song := new(Song)
	song.Media = media
	song.Title = title
	song.Id = id
	return song
}

package bot

import (
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

// Ffmpeg returns ffmpeg executable commans
func (song Song) Ffmpeg() *exec.Cmd {
	return exec.Command("ffmpeg", "-i", song.Media, "-f", "s16le", "-ar", strconv.Itoa(FRAME_RATE), "-ac",
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

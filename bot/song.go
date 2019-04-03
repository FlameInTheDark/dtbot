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

// Ffmpeg returns ffmpeg executable command
func (song Song) Ffmpeg(volume float32) *exec.Cmd {
	return exec.Command("ffmpeg", "-i", song.Media, "-f", "s16le", /*fmt.Sprintf("-filter:a volume=%.2f", volume),*/ "-ar", strconv.Itoa(FRAME_RATE), "-ac",
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

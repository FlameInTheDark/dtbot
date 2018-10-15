package bot

import (
	"github.com/bwmarrin/discordgo"
)

type SongQueue struct {
	list    []Song
	current *Song
	Running bool
}

func (queue SongQueue) Get() []Song {
	return queue.list
}

func (queue *SongQueue) Set(list []Song) {
	queue.list = list
}

func (queue *SongQueue) Add(song *Song) {
	queue.list = append(queue.list, *song)
}

func (queue *SongQueue) HasNext() bool {
	return len(queue.list) > 0
}

func (queue *SongQueue) Next() Song {
	song := queue.list[0]
	queue.list = queue.list[1:]
	queue.current = &song
	return song
}

func (queue *SongQueue) Clear() {
	queue.list = make([]Song, 0)
	queue.Running = false
	queue.current = nil
}

func (queue *SongQueue) Start(sess *Session, msg *discordgo.Message, callback func(string, *discordgo.Message)) {
	queue.Running = true
	for queue.HasNext() && queue.Running {
		song := queue.Next()
		callback(song.Title, msg)
		sess.PlayYoutube(song)
	}
	if !queue.Running {
		callback("stop", msg)
	} else {
		callback("finish", msg)
	}
}

func (queue *SongQueue) Current() *Song {
	return queue.current
}

func (queue *SongQueue) Pause() {
	queue.Running = false
}

func newSongQueue() *SongQueue {
	queue := new(SongQueue)
	queue.list = make([]Song, 0)
	return queue
}

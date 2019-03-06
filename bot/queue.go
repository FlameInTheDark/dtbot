package bot

// SongQueue struct contains songs array
type SongQueue struct {
	list    []Song
	current *Song
	Running bool
}

// Get returns songs array
func (queue *SongQueue) Get() []Song {
	return queue.list
}

// Set sets songs array
func (queue *SongQueue) Set(list []Song) {
	queue.list = list
}

// Add adds one song in songs array
func (queue *SongQueue) Add(song *Song) {
	queue.list = append(queue.list, *song)
}

// HasNext check if exist newx song in queue
func (queue *SongQueue) HasNext() bool {
	return len(queue.list) > 0
}

// Next returns next song from queue
func (queue *SongQueue) Next() Song {
	song := queue.list[0]
	queue.list = queue.list[1:]
	queue.current = &song
	return song
}

// Clear removes all songs from queue
func (queue *SongQueue) Clear() {
	queue.list = make([]Song, 0)
	queue.Running = false
	queue.current = nil
}

// Start starts queue playing
func (queue *SongQueue) Start(sess *Session, callback func(string)) {
	queue.Running = true
	for queue.HasNext() && queue.Running {
		song := queue.Next()
		callback(song.Title)
		_ = sess.PlayYoutube(song)
	}
	if !queue.Running {
		callback("stop")
	} else {
		callback("finish")
	}
}

// Current returns current song
func (queue *SongQueue) Current() *Song {
	return queue.current
}

// Pause pauses queue playing
func (queue *SongQueue) Pause() {
	queue.Running = false
}

func newSongQueue() *SongQueue {
	queue := new(SongQueue)
	queue.list = make([]Song, 0)
	return queue
}

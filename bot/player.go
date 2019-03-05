package bot

// RadioPlayer radio player struct
type RadioPlayer struct {
	Running bool
}

// Start starts radio playback
func (player *RadioPlayer) Start(sess *Session, source string, callback func(string)) {
	player.Running = true
	for player.Running {
		callback("Now playing `" + source + "`.")
		_=sess.Play(source)
		player.Running = false
	}
	if !player.Running {
		callback("Stopped playing.")
	} else {
		callback("Player closed.")
	}
}

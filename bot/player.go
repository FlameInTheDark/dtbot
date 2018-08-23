package bot

type RadioPlayer struct {
	Running bool
}

func (player *RadioPlayer) Start(sess *Session, source string, callback func(string)) {
	player.Running = true
	for player.Running {
		callback("Now playing `" + source + "`.")
		sess.Play(source)
		player.Running = false
	}
	if !player.Running {
		callback("Stopped playing.")
	} else {
		callback("Player closed.")
	}
}

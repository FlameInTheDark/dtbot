package bot

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strconv" // https://github.com/layeh/gopus

	"github.com/FlameInTheDark/gopus"
)

const (
	// CHANNELS count of audio channels
	CHANNELS int = 2
	// FRAME_RATE ...
	FRAME_RATE int = 48000
	// FRAME_SIZE ...
	FRAME_SIZE int = 960
	// MAX_BYTES max bytes per sample
	MAX_BYTES int = (FRAME_SIZE * 2) * 2
)

// Play start playback
func (connection *Connection) Play(source string, volume float32) error {
	if connection.playing {
		return errors.New("song already playing")
	}
	ffmpeg := exec.Command("ffmpeg", "-i", source, "-f", "s16le", "-reconnect", "1", "-reconnect_at_eof", "1", "-reconnect_streamed", "1", "-reconnect_delay_max", "2", "-filter:a", fmt.Sprintf("volume=%.3f", volume), "-ar", strconv.Itoa(FRAME_RATE), "-ac", strconv.Itoa(CHANNELS), "pipe:1")
	out, err := ffmpeg.StdoutPipe()
	if err != nil {
		return err
	}
	buffer := bufio.NewReaderSize(out, 16384)
	err = ffmpeg.Start()
	if err != nil {
		return err
	}
	defer func() {
		fferr := ffmpeg.Process.Kill()
		if fferr != nil {
			fmt.Println("FFMPEG close err: ", fferr)
		}
	}()
	connection.playing = true
	_ = connection.voiceConnection.Speaking(true)
	defer func() { _ = connection.voiceConnection.Speaking(false) }()

	encoder, err := gopus.NewEncoder(FRAME_RATE, CHANNELS, gopus.Audio)
	if err != nil {
		fmt.Println("Can's create a gopus encoder", err)
		return
	}

loop:
	for {
		select {
		case _ = <-connection.quitChan:
			break loop
		default:
			opus, err := encoder.Encode(buffer, FRAME_SIZE, MAX_BYTES)
			if err != nil {
				fmt.Println("Gopus encoding error,", err)
				return
			}
			if !voice.Ready || voice.OpusSend == nil {
				fmt.Printf("Discordgo not ready for opus packets. %+v : %+v", voice.Ready, voice.OpusSend)
				return
			}
			voice.OpusSend <- opus
		}
	}

	return nil
}

// PlayYoutube starts playing song from youtube
func (connection *Connection) PlayYoutube(ffmpeg *exec.Cmd) error {
	if connection.playing {
		return errors.New("song already playing")
	}
	out, err := ffmpeg.StdoutPipe()
	if err != nil {
		return err
	}
	buffer := bufio.NewReaderSize(out, 16384)
	err = ffmpeg.Start()
	if err != nil {
		return err
	}
	defer func() {
		fferr := ffmpeg.Process.Kill()
		if fferr != nil {
			fmt.Println("FFMPEG close err: ", fferr)
		}
	}()
	connection.playing = true
	defer func() {
		connection.playing = false
	}()
	_ = connection.voiceConnection.Speaking(true)
	defer func() { _ = connection.voiceConnection.Speaking(false) }()

	//audioBuffer := make([]int16, FRAME_SIZE*CHANNELS)

	encoder, err := gopus.NewEncoder(FRAME_RATE, CHANNELS, gopus.Audio)
	if err != nil {
		fmt.Println("Can's create a gopus encoder", err)
		return
	}
loop:
	for {
		select {
		case _ = <-connection.quitChan:
			break loop
		default:
			/*
				err = binary.Read(buffer, binary.LittleEndian, &audioBuffer)
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					return nil
				} else if err != nil {
					return err
				}
			*/
			opus, err := encoder.Encode(buffer, FRAME_SIZE, MAX_BYTES)
			if err != nil {
				fmt.Println("Gopus encoding error,", err)
				return
			}
			if !voice.Ready || voice.OpusSend == nil {
				fmt.Printf("Discordgo not ready for opus packets. %+v : %+v", voice.Ready, voice.OpusSend)
				return
			}
			voice.OpusSend <- audioBuffer
		}
		return nil
	}
}

// Stop stops playback
func (connection *Connection) Stop() {
	connection.playing = false
	connection.quitChan <- struct{}{}
}

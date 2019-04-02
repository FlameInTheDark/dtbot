package bot

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv" // https://github.com/layeh/gopus

	"github.com/FlameInTheDark/gopus"
	"github.com/bwmarrin/discordgo"
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
	ffmpeg := exec.Command("ffmpeg", "-i", source, "-f", "s16le", fmt.Sprintf("-filter:a volume=%.2f", volume), "-ar", strconv.Itoa(FRAME_RATE), "-ac", strconv.Itoa(CHANNELS), "pipe:1")
	connection.stopRunning = false
	out, err := ffmpeg.StdoutPipe()
	if err != nil {
		return err
	}
	buffer := bufio.NewReaderSize(out, 16384)
	err = ffmpeg.Start()
	if err != nil {
		return err
	}
	connection.playing = true
	defer func() {
		connection.playing = false
	}()
	_ = connection.voiceConnection.Speaking(true)
	defer func() { _ = connection.voiceConnection.Speaking(false) }()
	if connection.send == nil {
		connection.send = make(chan []int16, 2)
	}
	go connection.sendPCM(connection.voiceConnection, connection.send)
	for {
		if connection.stopRunning {
			_ = ffmpeg.Process.Kill()
			break
		}
		audioBuffer := make([]int16, FRAME_SIZE*CHANNELS)
		err = binary.Read(buffer, binary.LittleEndian, &audioBuffer)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		if err != nil {
			return err
		}
		connection.send <- audioBuffer
	}
	return nil
}

// sendPCM sends pulse code modulation to discord voice channel
func (connection *Connection) sendPCM(voice *discordgo.VoiceConnection, pcm <-chan []int16) {
	connection.lock.Lock()
	if connection.sendpcm || pcm == nil {
		connection.lock.Unlock()
		return
	}
	connection.sendpcm = true
	connection.lock.Unlock()
	defer func() {
		connection.sendpcm = false
	}()
	encoder, err := gopus.NewEncoder(FRAME_RATE, CHANNELS, gopus.Audio)
	if err != nil {
		fmt.Println("NewEncoder error,", err)
		return
	}
	for {
		receive, ok := <-pcm
		if !ok {
			fmt.Println("PCM channel closed")
			return
		}
		opus, err := encoder.Encode(receive, FRAME_SIZE, MAX_BYTES)
		if err != nil {
			fmt.Println("Encoding error,", err)
			return
		}
		if !voice.Ready || voice.OpusSend == nil {
			fmt.Printf("Discordgo not ready for opus packets. %+v : %+v", voice.Ready, voice.OpusSend)
			return
		}
		voice.OpusSend <- opus
	}
}

// PlayYoutube starts playing song from youtube
func (connection *Connection) PlayYoutube(ffmpeg *exec.Cmd) error {
	if connection.playing {
		return errors.New("song already playing")
	}
	connection.stopRunning = false
	out, err := ffmpeg.StdoutPipe()
	if err != nil {
		return err
	}
	buffer := bufio.NewReaderSize(out, 16384)
	err = ffmpeg.Start()
	if err != nil {
		return err
	}
	connection.playing = true
	defer func() {
		connection.playing = false
	}()
	_ = connection.voiceConnection.Speaking(true)
	defer func() { _ = connection.voiceConnection.Speaking(false) }()
	if connection.send == nil {
		connection.send = make(chan []int16, 2)
	}
	go connection.sendPCM(connection.voiceConnection, connection.send)
	for {
		if connection.stopRunning {
			_ = ffmpeg.Process.Kill()
			break
		}
		audioBuffer := make([]int16, FRAME_SIZE*CHANNELS)
		err = binary.Read(buffer, binary.LittleEndian, &audioBuffer)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		if err != nil {
			return err
		}
		connection.send <- audioBuffer
	}
	return nil
}

// Stop stops playback
func (connection *Connection) Stop() {
	connection.stopRunning = true
	connection.playing = false
}

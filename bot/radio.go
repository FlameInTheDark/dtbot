package bot

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

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

func (connection *Connection) EncodeOpusAndSend(reader io.Reader) error {
	if connection.playing {
		return errors.New("song already playing")
	}

	connection.playing = true
	_ = connection.voiceConnection.Speaking(true)
	defer func() { _ = connection.voiceConnection.Speaking(false) }()

	breader := bufio.NewReaderSize(reader, 16384)
	var buffer [FRAME_SIZE*CHANNELS]int16
	encoder, err := gopus.NewEncoder(FRAME_RATE, CHANNELS, gopus.Audio)
	if err != nil {
		fmt.Println("Can's create a gopus encoder", err)
		return err
	}

	voice := connection.voiceConnection
loop:
	for {
		select {
		case _ = <-connection.quitChan:
			break loop
		default:
			err = binary.Read(breader, binary.LittleEndian, &buffer)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return nil
			} else if err != nil {
				return err
			}
			opus, err := encoder.Encode(buffer[:], FRAME_SIZE, MAX_BYTES)
			if err != nil {
				fmt.Println("Gopus encoding error,", err)
				return err
			}
			if !voice.Ready || voice.OpusSend == nil {
				fmt.Printf("Discordgo not ready for opus packets. %+v : %+v", voice.Ready, voice.OpusSend)
				return err
			}
			voice.OpusSend <- opus
		}
	}

	return nil
}

// Stop stops playback
func (connection *Connection) Stop() {
	connection.playing = false
	connection.quitChan <- struct{}{}
}

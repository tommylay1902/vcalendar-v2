package audio

import (
	"context"
	"fmt"
	"log"
	"time"

	"changeme/model"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/gordonklaus/portaudio"
	"github.com/olebedev/when"
)

type VoskCommunication struct {
	ctx         context.Context
	ws          *websocket.Conn
	stream      *portaudio.Stream
	audioBuffer []int16
}

func InitVoskCommunication(ctx context.Context, ws *websocket.Conn, stream *portaudio.Stream, audioBuffer []int16) *VoskCommunication {
	return &VoskCommunication{
		ctx:         ctx,
		ws:          ws,
		stream:      stream,
		audioBuffer: audioBuffer,
	}
}

func (vc *VoskCommunication) SendAudioToWs(messageChan chan any, errorChan chan error, doneChan chan struct{}) {
	defer close(messageChan)
	defer close(errorChan)
	defer close(doneChan)

	for {
		var msg any
		err := wsjson.Read(vc.ctx, vc.ws, &msg)
		if err != nil {
			errorChan <- err
			return
		}
		select {
		case messageChan <- msg:
		case <-doneChan:
			return
		}
	}
}

func (vc *VoskCommunication) RecordAudio(wc *when.Parser, gc model.GcClient, qc model.QdrantClient,
	messageChan chan any, errorChan chan error, stopChan chan struct{},
) {
	recording := true
	for recording {

		// Read audio from microphone
		err := vc.stream.Read()
		if err != nil {
			log.Printf("Error reading audio: %v", err)
			panic(err)
		}

		// Send audio to Vosk when we have enough samples
		if len(vc.audioBuffer) >= 160 { // ~10ms of 16kHz audio
			audioBytes := make([]byte, len(vc.audioBuffer)*2)
			for i, sample := range vc.audioBuffer {
				audioBytes[i*2] = byte(sample)
				audioBytes[i*2+1] = byte(sample >> 8)
			}

			// Send raw audio to Vosk
			err = vc.ws.Write(vc.ctx, websocket.MessageBinary, audioBytes)
			if err != nil {
				log.Printf("Error sending audio: %v", err)
				break
			}
		}

		// Check for messages or stop signal
		select {
		case msg := <-messageChan:

			finalText := handleVoskMessage(msg)
			var date *time.Time
			if finalText != nil {
				r, err := wc.Parse(*finalText, time.Now())
				if err != nil {
					fmt.Println(error.Error)
					panic(err)
				}
				if r == nil {
					fmt.Println("no matches found")
				} else {
					date = &r.Time
					fmt.Println(date)
				}
			}

			operation := qc.GetOperation(finalText)
			switch operation {
			case "List":
				gc.GetEventsForTheDay(date)
			case "Add":
				fmt.Println("Creating event...")
			case "Delete":
				fmt.Println("Deleting event...")

			}
		case err := <-errorChan:
			if err != nil {
				log.Printf("WebSocket error: %v", err)
			}
			recording = false
		case <-stopChan:
			recording = false
		default:
			// Continue recording
		}
	}
}

func handleVoskMessage(msg any) *string {
	// currPartial := []string{}
	// Try to parse as JSON object
	if m, ok := msg.(map[string]any); ok {
		if text, ok := m["text"].(string); ok && text != "" {
			fmt.Print("\r\033[2K") // \033[2K clears entire line

			fmt.Printf("Final: %s\n", text)
			return &text

		} else if partial, ok := m["partial"].(string); ok && partial != "" {
			fmt.Printf("\rListening: %s", partial)
		}
	} else if str, ok := msg.(string); ok {
		fmt.Printf("Message: %s\n", str)
	}
	return nil
}

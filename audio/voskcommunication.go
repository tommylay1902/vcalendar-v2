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
	config      map[string]any
}

func InitVoskCommunication(ctx context.Context, ws *websocket.Conn, stream *portaudio.Stream, audioBuffer []int16, config map[string]any) *VoskCommunication {
	return &VoskCommunication{
		ctx:         ctx,
		ws:          ws,
		stream:      stream,
		audioBuffer: audioBuffer,
		config:      config,
	}
}

func (vc *VoskCommunication) StartVoskCommunication() {
	writeCtx, wsCancel := context.WithTimeout(vc.ctx, 100*time.Millisecond)
	defer wsCancel()
	err := wsjson.Write(writeCtx, vc.ws, vc.config)
	if err != nil {
		fmt.Println("error writing data to websocket ")
	}
}

func (vc *VoskCommunication) FormatWebsocketToJson(messageChan chan any, errorChan chan error, stopChan chan struct{}) {
	defer close(messageChan)
	defer close(errorChan)
	defer func() {
		if r := recover(); r != nil {
			// Ignore "send on closed channel" panic
			fmt.Println("Recovered from panic in FormatWebsocketToJson:", r)
		}
	}()
	for {
		select {
		case <-stopChan:
			return
		case <-vc.ctx.Done():
			return
		default:
			var msg any
			err := wsjson.Read(vc.ctx, vc.ws, &msg)
			fmt.Println(msg)
			if err != nil {
				fmt.Println("err reading from websocket")
				errorChan <- err
				return
			}
			messageChan <- msg
		}
	}
}

func (vc *VoskCommunication) HandleMessage(messageChan chan any, errorChan chan error, stopChan chan struct{}) {
	for {
		select {
		case msg := <-messageChan:
			fmt.Println("Received message from Vosk:")
			fmt.Printf("  %+v\n", msg)
		case err := <-errorChan:
			fmt.Printf("WebSocket error: %v\n", err)
			return
		case <-stopChan:
			fmt.Println("Stopping message handler")
			return
		}
	}
}

func (vc *VoskCommunication) RecordAudioTest(messageChan chan any, errorChan chan error, stopChan chan struct{}) {
	defer func() {
		// Close channels safely
		select {
		case <-stopChan:
			// Already closed, do nothing
		default:
			close(stopChan)
		}
	}()
	for {
		select {
		case <-vc.ctx.Done():
			return
		case <-stopChan:
			return
		default:
			err := vc.stream.Read()
			if err != nil {
				select {
				case errorChan <- err:
					// Error sent
				case <-vc.ctx.Done():
					// Context cancelled while trying to send error
				}
				return
			}
			if len(vc.audioBuffer) >= 160 {
				audioBytes := make([]byte, len(vc.audioBuffer)*2)
				for i, sample := range vc.audioBuffer {
					audioBytes[i*2] = byte(sample)
					audioBytes[i*2+1] = byte(sample >> 8)
				}
				err := vc.ws.Write(vc.ctx, websocket.MessageBinary, audioBytes)
				if err != nil {
					break
				}
			}

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

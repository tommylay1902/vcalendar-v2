package audio

import (
	"context"
	"fmt"
	"log"
	"time"
	"vcalendar-v2/model"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/gordonklaus/portaudio"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type VoskCommunication struct {
	ctx                 context.Context
	ws                  *websocket.Conn
	stream              *portaudio.Stream
	audioBuffer         []int16
	config              map[string]any
	gc                  *model.GcClient
	qc                  *model.QdrantClient
	wc                  *when.Parser
	finalTranscriptChan chan string
}

func InitVoskCommunication(ctx context.Context, ws *websocket.Conn, stream *portaudio.Stream, audioBuffer []int16, config map[string]any, gc *model.GcClient) *VoskCommunication {
	qc, err := model.InitializeQdrantClient()
	if err != nil {
		fmt.Println("error initalizing qc client")
		panic(err)
	}
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	return &VoskCommunication{
		ctx:                 ctx,
		ws:                  ws,
		stream:              stream,
		audioBuffer:         audioBuffer,
		config:              config,
		gc:                  gc,
		qc:                  qc,
		wc:                  w,
		finalTranscriptChan: make(chan string, 128), // Add buffer size
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
	defer func() {
		// Close channels safely
		closeSafe(messageChan)
		closeSafe(errorChan)
		if r := recover(); r != nil {
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
			if err != nil {
				fmt.Println("err reading from websocket")
				if !isContextError(err) {
					select {
					case errorChan <- err:
					case <-vc.ctx.Done():
					case <-stopChan:
					}
				}
			}
			messageChan <- msg
		}
	}
}

func closeSafe[T any](ch chan T) {
	defer func() { recover() }()
	select {
	case _, ok := <-ch:
		if !ok {
			return // Already closed
		}
	default:
		close(ch)
	}
}

func isContextError(err error) bool {
	return err == context.Canceled ||
		err == context.DeadlineExceeded ||
		err.Error() == "context canceled" ||
		err.Error() == "context deadline exceeded"
}

func (vc *VoskCommunication) HandleMessage(messageChan chan any, errorChan chan error, stopChan chan struct{}) {
	for {
		select {
		case msg := <-messageChan:
			// First, type assert msg to map[string]any
			if m, ok := msg.(map[string]any); ok {
				if text, ok := m["text"].(string); ok && text != "" {
					// fmt.Printf("\nFinal: %s\n", text)
					application.Get().Event.Emit("vcalendar-v2:send-transcription", model.Transcription{
						Message: text,
						IsFinal: true,
					})
					select {
					case vc.finalTranscriptChan <- text:
					default:
						// Skip if channel is full
					}
					// vc.findOperation(messageChan, errorChan, stopChan)
				} else if partial, ok := m["partial"].(string); ok && partial != "" {
					// fmt.Printf("Listening: %s", partial)
					application.Get().Event.Emit("vcalendar-v2:send-transcription", model.Transcription{
						Message: partial,
						IsFinal: false,
					})
				}
			} else if str, ok := msg.(string); ok {
				fmt.Printf("Message: %s\n", str)
			}
		case err := <-errorChan:
			fmt.Printf("WebSocket error: %v\n", err)
			return
		case <-stopChan:
			fmt.Println("Stopping message handler")
			return
		}
	}
}

func (vc *VoskCommunication) ProcessTranscripts(stopChan chan struct{}) {
	for {
		select {
		case text := <-vc.finalTranscriptChan:
			vc.processFinalTranscript(text)
		case <-stopChan:
			return
		}
	}
}

func (vc *VoskCommunication) processFinalTranscript(text string) {
	fmt.Println("Processing:", text)

	// Parse the date from the text
	r, err := vc.wc.Parse(text, time.Now())
	if err != nil {
		fmt.Println("Error parsing date:", err)
		// Don't panic, just return
		return
	}

	var date *time.Time
	if r == nil {
		fmt.Println("no matches found")
	} else {
		date = &r.Time
		fmt.Println("Parsed date:", date)
	}
	// Get the operation from Qdrant
	operation := vc.qc.GetOperation(&text)
	// Execute the operation
	switch operation {
	case "List":

		events := vc.gc.GetEventsForTheDay(date)
		if events != nil {
			application.Get().Event.Emit("vcalendar-v2:send-events", model.CalendarEvents{
				Summary:     events.Summary,
				Description: events.Description,
				Events:      events.Items,
			})
		}
	case "Add":
		fmt.Println("Creating event...")
	case "Delete":
		fmt.Println("Deleting event...")
	default:
		fmt.Println("Unknown operation:", operation)
	}
}

func (vc *VoskCommunication) RecordAudioTest(messageChan chan any, errorChan chan error, stopChan chan struct{}) {
	defer func() {
		// Don't close stopChan here - let the service handle it
		if r := recover(); r != nil {
			fmt.Println("Recovered in RecordAudioTest:", r)
		}
	}()

	for {
		select {
		case <-vc.ctx.Done():
			fmt.Println("RecordAudioTest: context done")
			return
		case <-stopChan:
			fmt.Println("RecordAudioTest: stop channel closed")
			return
		default:
			err := vc.stream.Read()
			if err != nil {
				fmt.Printf("RecordAudioTest: stream read error: %v\n", err)
				// Don't send to errorChan - just exit
				return
			}

			if len(vc.audioBuffer) >= 160 {
				audioBytes := make([]byte, len(vc.audioBuffer)*2)
				for i, sample := range vc.audioBuffer {
					audioBytes[i*2] = byte(sample)
					audioBytes[i*2+1] = byte(sample >> 8)
				}

				// Check context before writing
				select {
				case <-vc.ctx.Done():
					return
				case <-stopChan:
					return
				default:
					err := vc.ws.Write(vc.ctx, websocket.MessageBinary, audioBytes)
					if err != nil {
						fmt.Printf("RecordAudioTest: write error: %v\n", err)
						return
					}
				}
			}
		}
	}
}

func (vc *VoskCommunication) findOperation(messageChan chan any, errorChan chan error, stopChan chan struct{}) {
	fmt.Println("find operation!")
	select {
	case msg := <-messageChan:
		fmt.Println("inside ")
		finalText := handleVoskMessage(msg)
		var date *time.Time
		if finalText != nil {
			r, err := vc.wc.Parse(*finalText, time.Now())
			if err != nil {
				fmt.Println(err.Error)
				panic(err)
			}
			if r == nil {
				fmt.Println("no matches found")
			} else {
				date = &r.Time
				fmt.Println(date)
			}
		}

		operation := vc.qc.GetOperation(finalText)
		switch operation {
		case "List":
			vc.gc.GetEventsForTheDay(date)
		case "Add":
			fmt.Println("Creating event...")
		case "Delete":
			fmt.Println("Deleting event...")
		}
	case err := <-errorChan:
		if err != nil {
			log.Printf("WebSocket error: %v", err)
		}
	case <-stopChan:
		return
	default:
		// Continue recording
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

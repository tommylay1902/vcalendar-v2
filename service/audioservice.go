package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"changeme/audio"

	"github.com/coder/websocket"
	"github.com/gordonklaus/portaudio"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type AudioService struct {
	stream      *portaudio.Stream
	audioSample []int16
	ws          *websocket.Conn
	vc          *audio.VoskCommunication
	stop        chan struct{}
	ctx         context.Context
	cancel      context.CancelFunc
	App         *application.App
}

func (as *AudioService) StartRecord() {
	as.ctx, as.cancel = context.WithCancel(context.Background())
	err := portaudio.Initialize()
	if err != nil {
		fmt.Println("error initializing portaudio")
		panic(err)
	}
	as.stop = make(chan struct{})
	messageChan := make(chan any, 100) // Buffered channel
	errorChan := make(chan error, 10)
	as.audioSample = make([]int16, 256)
	stream, err := portaudio.OpenDefaultStream(
		1, 0,
		16000, len(as.audioSample),
		as.audioSample,
	)
	if err != nil {
		fmt.Println("error opening stream")
		panic(err)

	}
	as.stream = stream

	err = stream.Start()
	if err != nil {
		fmt.Println("error starting stream")
		panic(err)
	}
	wsCtx, cancel := context.WithCancel(as.ctx)
	defer cancel()
	// WebSocket connection to Vosk
	ws, _, err := websocket.Dial(wsCtx, "ws://localhost:2700", nil)
	as.ws = ws
	if err != nil {
		fmt.Println("error trying to connnect to websocket")
		panic(err)
	}

	// Send configuration to Vosk
	config := map[string]any{
		"config": map[string]any{
			"sample_rate": 16000.0, // Vosk expects 16kHz
		},
	}
	as.vc = audio.InitVoskCommunication(as.ctx, ws, as.stream, as.audioSample, config)
	as.vc.StartVoskCommunication()

	go as.vc.RecordAudioTest(messageChan, errorChan, as.stop)
	go as.vc.FormatWebsocketToJson(messageChan, errorChan, as.stop)
	go as.vc.HandleMessage(messageChan, errorChan, as.stop)
}

func (as *AudioService) StopRecord() {
	// Signal all goroutines to stop
	close(as.stop)

	// Cancel context to propagate cancellation
	if as.cancel != nil {
		as.cancel()
	}

	// Wait a bit for goroutines to clean up
	time.Sleep(100 * time.Millisecond)

	// Cleanup resources
	as.cleanup()
}

func (as *AudioService) cleanup() {
	// Stop and close stream if exists
	if as.stream != nil {
		as.stream.Stop()
		as.stream.Close()
		as.stream = nil
	}

	// Close WebSocket if exists
	if as.ws != nil {
		as.ws.Close(websocket.StatusNormalClosure, "cleanup")
		as.ws = nil
	}

	// Clear other references
	as.audioSample = nil
	as.vc = nil

	// Terminate PortAudio - use defer to handle panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic during PortAudio termination: %v", r)
		}
		portaudio.Terminate()
	}()
}

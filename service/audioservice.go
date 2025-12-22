package service

import (
	"context"
	"fmt"
	"time"
	"vcalendar-v2/audio"
	"vcalendar-v2/model"

	"github.com/coder/websocket"
	"github.com/gordonklaus/portaudio"
)

type AudioService struct {
	stream      *portaudio.Stream
	audioSample []int16
	ws          *websocket.Conn
	vc          *audio.VoskCommunication
	stop        chan struct{}
	ctx         context.Context
	cancel      context.CancelFunc
	Gc          *model.GcClient
}

func (as *AudioService) StartRecord() {
	as.ctx, as.cancel = context.WithCancel(context.Background())

	err := portaudio.Initialize()
	if err != nil {
		fmt.Println("error initializing portaudio")
		panic(err)
	}

	// Create stop channel
	as.stop = make(chan struct{})
	messageChan := make(chan any, 100)
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

	wsCtx, cancel := context.WithTimeout(as.ctx, 5*time.Second)
	defer cancel()

	ws, _, err := websocket.Dial(wsCtx, "ws://localhost:2700", nil)
	as.ws = ws
	if err != nil {
		fmt.Println("error trying to connect to websocket")
		panic(err)
	}

	config := map[string]any{
		"config": map[string]any{
			"sample_rate": 16000.0,
		},
	}

	as.vc = audio.InitVoskCommunication(as.ctx, ws, stream, as.audioSample, config, as.Gc)
	as.vc.StartVoskCommunication()

	// Start all goroutines with the same stop channel
	go as.vc.RecordAudioTest(messageChan, errorChan, as.stop)
	go as.vc.FormatWebsocketToJson(messageChan, errorChan, as.stop)
	go as.vc.HandleMessage(messageChan, errorChan, as.stop)
	go as.vc.ProcessTranscripts(as.stop)

	fmt.Println("Audio service started successfully")
}

func (as *AudioService) StopRecord() {
	fmt.Println("Stopping audio service...")

	// 1. Close the stop channel FIRST to signal all goroutines
	if as.stop != nil {
		close(as.stop)
		// Wait a moment for goroutines to receive the signal
		time.Sleep(50 * time.Millisecond)
	}

	// 2. Cancel the context
	if as.cancel != nil {
		as.cancel()
	}

	// 3. Close WebSocket (this will unblock wsjson.Read)
	if as.ws != nil {
		// Send EOF in a goroutine with timeout
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			as.ws.Write(ctx, websocket.MessageText, []byte(`{"eof":1}`))
		}()

		as.ws.Close(websocket.StatusNormalClosure, "stopping")
		as.ws = nil
	}

	// 4. Stop the audio stream
	if as.stream != nil {
		as.stream.Stop()
		as.stream.Close()
		as.stream = nil
	}

	// 5. Terminate portaudio
	portaudio.Terminate()

	fmt.Println("Audio service stopped")
}

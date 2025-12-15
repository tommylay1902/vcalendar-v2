package service

import (
	"fmt"

	"changeme/audio"

	"github.com/gordonklaus/portaudio"
)

type AudioService struct {
	stream      *portaudio.Stream
	audioSample []int16
	vc          *audio.VoskCommunication
}

func (as *AudioService) StartRecord() {
	err := portaudio.Initialize()
	if err != nil {
		panic(err)
	}

	as.audioSample = make([]int16, 128)
	stream, err := portaudio.OpenDefaultStream(
		1, 0,
		16000, len(as.audioSample),
		func(in []int16) {
			fmt.Println(in)
		},
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
}

func (as *AudioService) StopRecord() {
	if as.stream == nil {
		return // Nothing to do
	}

	// Just stop and close the stream
	streamStopErr := as.stream.Stop()
	streamCloseErr := as.stream.Close()
	if streamStopErr != nil {
		panic(streamStopErr)
	}
	if streamCloseErr != nil {
		panic(streamCloseErr)
	}
	// Terminate portaudio (assumes Initialize was called)
	paErr := portaudio.Terminate()
	if paErr != nil {
		panic(paErr)
	}
	// Clear references
	as.stream = nil
	as.audioSample = nil
}

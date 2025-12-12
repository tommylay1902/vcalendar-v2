package service

import "github.com/gordonklaus/portaudio"

type AudioService struct {
	stream      *portaudio.Stream
	audioSample []int16
}

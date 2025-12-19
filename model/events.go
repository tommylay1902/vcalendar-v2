package model

type GoogleAuth struct {
	TokenNeeded bool
}

type AuthCodeToken struct {
	Token string
}

type Transcription struct {
	Message string
	IsFinal bool
}

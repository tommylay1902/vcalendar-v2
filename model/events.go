package model

import "google.golang.org/api/calendar/v3"

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

type CalendarEventRequest struct {
	Data calendar.Events `json:"data"`
}

type CalendarEvents struct {
	Summary     string
	Description string
	Events      []*calendar.Event `json:"events"`
}

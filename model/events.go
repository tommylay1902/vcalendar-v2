package model

import (
	"time"

	"google.golang.org/api/calendar/v3"
)

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
	Date        *time.Time
	Summary     string
	Description string
	Events      []*calendar.Event `json:"events"`
}

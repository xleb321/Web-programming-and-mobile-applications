package models

import (
	"time"
)

type EventType string

const (
	EventUserRegistered EventType = "user.registered"
)

type Event struct {
	EventID   string          `json:"eventId"`
	EventType EventType       `json:"eventType"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   interface{}     `json:"payload"`
	Metadata  EventMetadata   `json:"metadata"`
}

type EventMetadata struct {
	Attempt       int    `json:"attempt"`
	SourceService string `json:"sourceService"`
}

type UserRegisteredPayload struct {
	UserID      string `json:"userId"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName,omitempty"`
}
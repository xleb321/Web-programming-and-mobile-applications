package handlers

import (
	"encoding/json"
	"log"

	"rugram-api/internal/models"
	"rugram-api/internal/service"
)

type EventHandler struct {
    emailService *service.EmailService
}

func NewEventHandler(emailService *service.EmailService) *EventHandler {
    return &EventHandler{
        emailService: emailService,
    }
}

func (h *EventHandler) HandleUserRegistered(event models.Event) error {
    // Convert payload to bytes then to map
    payloadBytes, err := json.Marshal(event.Payload)
    if err != nil {
        log.Printf("Failed to marshal payload for event: %s", event.EventID)
        return err
    }

    var payload models.UserRegisteredPayload
    if err := json.Unmarshal(payloadBytes, &payload); err != nil {
        log.Printf("Failed to parse payload for event: %s", event.EventID)
        return err
    }

    displayName := payload.DisplayName
    if displayName == "" {
        displayName = payload.Email
    }

    log.Printf("Processing welcome email for user: %s (%s)", payload.UserID, payload.Email)

    if h.emailService == nil || !h.emailService.IsConfigured() {
        log.Printf("SMTP not configured, skipping email for user: %s", payload.UserID)
        return nil
    }

    return h.emailService.SendWelcomeEmail(payload.Email, displayName, payload.UserID)
}
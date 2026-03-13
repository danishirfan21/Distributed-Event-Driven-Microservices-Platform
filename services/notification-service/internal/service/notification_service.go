package service

import (
	"context"
	"log"
)

type NotificationService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

type notificationService struct{}

func NewNotificationService() NotificationService {
	return &notificationService{}
}

func (s *notificationService) SendEmail(ctx context.Context, to, subject, body string) error {
	log.Printf("[NOTIFICATION] Sending Email to: %s, Subject: %s, Body: %s", to, subject, body)
	return nil
}

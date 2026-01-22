package services

import (
	"encoding/json"
	"errors"

	"github.com/SherClockHolmes/webpush-go"
)

type PushService struct {
	PublicKey  string
	PrivateKey string
	Subject    string
}

func NewPushService(pub, priv, subject string) *PushService {
	return &PushService{
		PublicKey:  pub,
		PrivateKey: priv,
		Subject:    subject,
	}
}

type PushSub struct {
	Endpoint string
	P256dh   string
	Auth     string
}

func (s *PushService) Send(sub PushSub, title, message string) error {
	if sub.Endpoint == "" {
		return errors.New("invalid subscription")
	}

	payload := map[string]string{
		"title":   title,
		"message": message,
	}
	data, _ := json.Marshal(payload)

	webSub := &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			// ⚠️ pakai yang cocok sama versi library kamu:
			// P256DH: sub.P256dh,
			P256dh: sub.P256dh,
			Auth:   sub.Auth,
		},
	}

	opts := &webpush.Options{
		Subscriber:      s.Subject,
		VAPIDPublicKey:  s.PublicKey,
		VAPIDPrivateKey: s.PrivateKey,
		TTL:             60,
	}

	_, err := webpush.SendNotification(data, webSub, opts)
	return err
}

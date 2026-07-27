package adapters

import (
	"fmt"

	expo "github.com/wagon-official/expo-notifications-sdk-golang"
)

type PushClient struct {
	client *expo.PushClient
}

func NewPushClient() *PushClient {
	return &PushClient{
		client: expo.NewPushClient(nil),
	}
}

// Send delivers a push notification to a single Expo push token.
// Returns an error instead of panicking, so a failed push never crashes the server.
func (p *PushClient) SendPushNotification(token, title, body string, data map[string]interface{}) error {
	pushToken, err := expo.NewExpoPushToken(token)
	if err != nil {
		return fmt.Errorf("invalid push token: %w", err)
	}

	responses, err := p.client.Publish(&expo.PushMessage{
		To:       []expo.ExpoPushToken{pushToken},
		Title:    title,
		Body:     body,
		Data:     data,
		Sound:    "default",
		Priority: expo.DefaultPriority,
	})
	if err != nil {
		return fmt.Errorf("failed to publish push: %w", err)
	}

	for _, res := range responses {
		if res.ValidateResponse() != nil {
			fmt.Println("push rejected for token:", res.PushMessage.To)
		}
	}
	return nil
}
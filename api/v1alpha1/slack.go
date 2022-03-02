package v1alpha1

import (
	"fmt"
)

// SlackNotifier defines the spec for integrating with Slack
// +k8s:openapi-gen=true
type SlackNotifier struct {
	WebhookUrl string `json:"webhook_url"`
	Channel    string `json:"channel,omitempty"`
	Username   string `json:"username,omitempty"`
	IconEmoji  string `json:"icon_emoji,omitempty"`
}

func (n *SlackNotifier) Send(message string) error {
	channel := ``
	if n.Channel != `` {
		channel = fmt.Sprintf(`, "channel":"%s"`, n.Channel)
	}
	username := ``
	if n.Username != `` {
		username = fmt.Sprintf(`, "username":"%s"`, n.Username)
	}
	icon_emoji := ``
	if n.IconEmoji != `` {
		icon_emoji = fmt.Sprintf(`, "icon_emoji":"%s"`, n.IconEmoji)
	}
	payload := fmt.Sprintf(`{"text":"%s"%s%s%s}`, escapeString(message), channel, username, icon_emoji)
	return sendJSON(n.WebhookUrl, payload)
}

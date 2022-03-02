package v1alpha1

import (
	"fmt"
)

// FeishuNotifier defines the spec for integrating with Slack
// +k8s:openapi-gen=true
type FeishuNotifier struct {
	WebhookUrl string `json:"webhook_url"`
}

func (n *FeishuNotifier) Send(message string) error {
	payload := fmt.Sprintf(`{"msg_type":"text","content":{"text":"%s"}}`, escapeString(message))
	return sendJSON(n.WebhookUrl, payload)
}

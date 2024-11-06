package models

type DTO_EmqxWebhookData struct {
	ID                string `json:"id,omitempty"`
	Topic             string `json:"topic,omitempty"`
	Payload           string `json:"payload,omitempty"`
	ClientID          string `json:"clientid,omitempty"`
	Event             string `json:"event,omitempty"`
	Username          string `json:"username,omitempty"`
	Timestamp         int64  `json:"timestamp,omitempty"`
	Qos               int    `json:"qos,omitempty"`
	PublishReceivedAt int64  `json:"publish_received_at,omitempty"`
	Peerhost          string `json:"peerhost,omitempty"`
	Node              string `json:"node,omitempty"`
	PubProps          struct {
		UserProperty struct {
		} `json:"User-Property,omitempty"`
	} `json:"pub_props,omitempty"`

	Metadata struct {
		RuleID string `json:"rule_id,omitempty"`
	} `json:"metadata,omitempty"`
	Flags struct {
		Retain bool `json:"retain,omitempty"`
		Dup    bool `json:"dup,omitempty"`
	} `json:"flags,omitempty"`
}

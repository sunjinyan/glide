package subscription

import "github.com/glide-im/glide/pkg/messages"

// Message is a message that can be publishing to channel.
type Message interface {
	GetFrom() SubscriberID

	// GetChatMessage convert message body to *messages.ChatMessage
	GetChatMessage() (*messages.ChatMessage, error)
}

const (
	NotifyTypeOffline = 2001
	NotifyTypeOnline  = 2002
	NotifyTypeJoin    = 2003
	NotifyTypeLeave   = 2004

	NotifyOnlineMembers = 2005
)

type NotifyMessage struct {
	From string      `json:"from"`
	Type int         `json:"type"`
	Body interface{} `json:"body"`
}

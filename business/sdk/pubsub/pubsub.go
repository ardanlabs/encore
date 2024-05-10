// Package pubsub provides topics that can be access by the business layer.
package pubsub

import (
	"encore.dev/pubsub"
	"github.com/ardanlabs/encore/business/sdk/delegate"
)

// Delegate represents a topic for handling delegate calls.
var Delegate = pubsub.NewTopic[delegate.Data]("delegate", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})

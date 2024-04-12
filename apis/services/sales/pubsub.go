package sales

import (
	"context"

	"encore.dev/pubsub"
	"github.com/ardanlabs/encore/business/api/delegate"
	bpubsub "github.com/ardanlabs/encore/business/pubsub"
)

// We need a single subscription which will route a message to the
// delegate system.
var _ = pubsub.NewSubscription(bpubsub.Delegate, "handle-delegate-call",
	pubsub.SubscriptionConfig[delegate.Data]{
		Handler: pubsub.MethodHandler((*Service).DelegateHandler),
	},
)

// DelegateHandler receives a message from the pubsub system and passes it
// into the delegate system.
func (s *Service) DelegateHandler(ctx context.Context, data delegate.Data) error {
	s.log.Info("DelegateHandler", "data", data)
	return s.delegate.Call(ctx, data)
}

package bootstrap

import (
	"errors"

	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/glide/pkg/messaging"
	"github.com/glide-im/glide/pkg/subscription"
)

type Options struct {
	Messaging    messaging.Messaging
	Gate         gate.Gateway
	Subscription subscription.Interface
}

func Bootstrap(opts *Options) error {

	err := setupDependence(opts)
	if err != nil {
		return err
	}

	_, ok := opts.Gate.(gate.Server)
	if ok {
		return bootGatewayServer(opts)
	}
	_, ok = opts.Messaging.(messaging.Messaging)
	if ok {
		return bootMessagingServer(opts)
	}
	_, ok = opts.Subscription.(subscription.Server)
	if ok {
		return bootSubscriptionServer(opts)
	}

	return errors.New("no server found")
}

func setupDependence(opts *Options) error {
	m, ok := opts.Messaging.(messaging.Messaging)
	if ok {
		g, ok := opts.Gate.(gate.Gateway)
		if ok {
			m.SetGate(g)
		} else {
			return errors.New("gateway not found")
		}
		m.SetSubscription(opts.Subscription)
	}

	sb, ok := opts.Subscription.(subscription.Subscribe)
	if ok {
		sb.SetGateInterface(opts.Gate.(gate.DefaultGateway))
	}
	return nil
}

func bootSubscriptionServer(opts *Options) error {
	server, ok := opts.Subscription.(subscription.Server)
	if !ok {
		return errors.New("subscription server not implemented")
	}
	server.SetGateInterface(opts.Gate.(gate.DefaultGateway))
	return server.Run()
}

func bootMessagingServer(opts *Options) error {
	server, ok := opts.Messaging.(messaging.Messaging)
	if !ok {
		return errors.New("messaging does not implement Messaging.impl")
	}

	manager, ok := opts.Gate.(gate.Gateway)
	if ok {
		server.SetGate(manager)
	}
	server.SetSubscription(opts.Subscription)
	return nil
}

func bootGatewayServer(opts *Options) error {

	gateway, ok := opts.Gate.(gate.Server)
	if !ok {
		return errors.New("Gate is not a gateway server")
	}

	if opts.Messaging == nil {
		return errors.New("can't boot a gateway server without a Messaging interface")
	}

	gateway.SetMessageHandler(func(cliInfo *gate.Info, message *messages.GlideMessage) {
		err := opts.Messaging.Handle(cliInfo, message)
		if err != nil {
			// TODO: Log error
		}
	})

	return gateway.Run()
}

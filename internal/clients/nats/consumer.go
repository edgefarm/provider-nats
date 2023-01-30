package nats

import (
	"errors"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
)

// ConsumerList returns a list of consumer names for a given domain
func ConsumerList(c *Client, domain string, stream string) ([]string, error) {
	jsopts := []jsm.Option{}
	if domain != "" {
		jsopts = append(jsopts, jsm.WithDomain(domain))
	}

	c.conn.NewRespInbox()
	mgr, err := jsm.New(c.conn, jsopts...)
	if err != nil {
		return nil, err
	}

	names, err := mgr.ConsumerNames(stream)
	if err != nil {
		return nil, err
	}
	return names, nil
}

// ConsumerInfo returns the consumer info for a given consumer name for a given domain and stream
func ConsumerInfo(c *Client, domain string, consumer string, stream string) (*nats.ConsumerInfo, error) {
	jsctx, err := c.conn.JetStream(nats.Domain(domain))
	if err != nil {
		return nil, err
	}

	info, err := jsctx.ConsumerInfo(stream, consumer)
	if err != nil {
		if errors.Is(err, nats.ErrConsumerNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return info, nil
}

// CreateConsumer creates a new jetstream consumer with a given configuration for a given domain and stream
func (c *Client) CreateConsumer(domain string, stream string, config *nats.ConsumerConfig) error {
	jsOpts := []nats.JSOpt{}
	if domain != "" {
		jsOpts = append(jsOpts, nats.Domain(domain))
	}

	jsctx, err := c.conn.JetStream(jsOpts...)
	if err != nil {
		return err
	}

	_, err = jsctx.AddConsumer(stream, config)
	if err != nil {
		return err
	}

	return nil
}

// DeleteConsumer deletes a jetstream consumer with a given name for a given domain and stream
func (c *Client) DeleteConsumer(domain string, stream string, consumer string) error {
	jsOpts := []nats.JSOpt{}
	if domain != "" {
		jsOpts = append(jsOpts, nats.Domain(domain))
	}

	jsctx, err := c.conn.JetStream(jsOpts...)
	if err != nil {
		return err
	}

	err = jsctx.DeleteConsumer(stream, consumer)
	if err != nil {
		return err
	}

	return nil
}

// UpdateConsumer updates a jetstream consumer with a given configuration for a given domain and stream
func (c *Client) UpdateConsumer(domain string, stream string, config *nats.ConsumerConfig) error {
	jsOpts := []nats.JSOpt{}
	if domain != "" {
		jsOpts = append(jsOpts, nats.Domain(domain))
	}

	jsctx, err := c.conn.JetStream(jsOpts...)
	if err != nil {
		return err
	}

	_, err = jsctx.UpdateConsumer(stream, config)
	if err != nil {
		return err
	}

	return nil
}

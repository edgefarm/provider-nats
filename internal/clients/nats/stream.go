package nats

import (
	"errors"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
)

// StreamList returns a list of stream names for a given domain
func StreamList(c *Client, domain string) ([]string, error) {
	jsopts := []jsm.Option{}
	if domain != "" {
		jsopts = append(jsopts, jsm.WithDomain(domain))
	}

	c.conn.NewRespInbox()
	mgr, err := jsm.New(c.conn, jsopts...)
	if err != nil {
		return nil, err
	}

	names := []string{}

	err = mgr.EachStream(&jsm.StreamNamesFilter{}, func(stream *jsm.Stream) {
		names = append(names, stream.Name())
	})
	if err != nil {
		return nil, err
	}
	return names, nil
}

// StreamInfo returns the stream info for a given stream name for a given domain
func StreamInfo(c *Client, domain string, stream string) (*nats.StreamInfo, error) {
	jsctx, err := c.conn.JetStream(nats.Domain(domain))
	if err != nil {
		return nil, err
	}

	info, err := jsctx.StreamInfo(stream)
	if err != nil {
		if errors.Is(err, nats.ErrStreamNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return info, nil
}

// CreateStream creates a new jetstream stream with a given configuration for a given domain
func (c *Client) CreateStream(domain string, config *nats.StreamConfig) error {
	jsOpts := []nats.JSOpt{}
	if domain != "" {
		jsOpts = append(jsOpts, nats.Domain(domain))
	}

	jsctx, err := c.conn.JetStream(jsOpts...)
	if err != nil {
		return err
	}

	_, err = jsctx.AddStream(config)
	if err != nil {
		return err
	}

	return nil
}

// DeleteStream deletes a jetstream stream with a given name for a given domain
func (c *Client) DeleteStream(domain string, name string) error {
	jsOpts := []nats.JSOpt{}
	if domain != "" {
		jsOpts = append(jsOpts, nats.Domain(domain))
	}

	jsctx, err := c.conn.JetStream(jsOpts...)
	if err != nil {
		return err
	}

	err = jsctx.DeleteStream(name)
	if err != nil {
		return err
	}

	return nil
}

// UpdateStream updates a jetstream stream with a given configuration for a given domain
func (c *Client) UpdateStream(domain string, config *nats.StreamConfig) error {
	jsOpts := []nats.JSOpt{}
	if domain != "" {
		jsOpts = append(jsOpts, nats.Domain(domain))
	}

	jsctx, err := c.conn.JetStream(jsOpts...)
	if err != nil {
		return err
	}

	_, err = jsctx.UpdateStream(config)
	if err != nil {
		return err
	}

	return nil
}

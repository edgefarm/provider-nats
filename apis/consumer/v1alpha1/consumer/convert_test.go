package consumer

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestConvertToNatsMinimal(t *testing.T) {
	assert := assert.New(t)

	customConfig := &ConsumerConfig{
		Description:       "my consumer",
		DeliverPolicy:     "New",
		AckPolicy:         "Explicit",
		AckWait:           "2m",
		MaxDeliver:        -1,
		FilterSubject:     "foo",
		ReplayPolicy:      "Instant",
		MaxAckPending:     5000,
		InactiveThreshold: "1m",
		Replicas:          0,
		MemoryStorage:     false,
	}

	natsConfig, err := ConfigV1Alpha1ToNats("mystream", customConfig)
	assert.Nil(err)
	assert.Equal(natsConfig.Description, "my consumer")
	assert.Equal(natsConfig.DeliverPolicy, nats.DeliverNewPolicy)
	assert.Equal(natsConfig.AckPolicy, nats.AckExplicitPolicy)
	assert.Equal(natsConfig.AckWait, func() time.Duration {
		t, _ := time.ParseDuration("2m")
		return t
	}())
	assert.Equal(natsConfig.MaxDeliver, int(-1))
	assert.Equal(natsConfig.FilterSubject, "foo")
	assert.Equal(natsConfig.ReplayPolicy, nats.ReplayInstantPolicy)
	assert.Equal(natsConfig.MaxAckPending, int(5000))
	assert.Equal(natsConfig.InactiveThreshold, func() time.Duration {
		t, _ := time.ParseDuration("1m")
		return t
	}())
	assert.Equal(natsConfig.Replicas, 0)
	assert.Equal(natsConfig.MemoryStorage, false)
	assert.Equal(natsConfig.MaxWaiting, int(DefaultMaxWait))
}

func TestConvertToNatsPull(t *testing.T) {
	assert := assert.New(t)

	customConfig := &ConsumerConfig{
		Description:       "my consumer",
		DeliverPolicy:     "New",
		AckPolicy:         "Explicit",
		AckWait:           "2m",
		MaxDeliver:        -1,
		FilterSubject:     "foo",
		ReplayPolicy:      "Instant",
		MaxAckPending:     5000,
		InactiveThreshold: "1m",
		Replicas:          0,
		MemoryStorage:     false,
		PullConsumer: &PullConsumerSpec{
			MaxWaiting:         func() *int { i := 100; return &i }(),
			MaxRequestExpires:  "1m",
			MaxRequestBatch:    100,
			MaxRequestMaxBytes: 1024,
		},
	}

	natsConfig, err := ConfigV1Alpha1ToNats("mystream", customConfig)
	assert.Nil(err)
	assert.Equal(natsConfig.Description, "my consumer")
	assert.Equal(natsConfig.DeliverPolicy, nats.DeliverNewPolicy)
	assert.Equal(natsConfig.AckPolicy, nats.AckExplicitPolicy)
	assert.Equal(natsConfig.AckWait, func() time.Duration {
		t, _ := time.ParseDuration("2m")
		return t
	}())
	assert.Equal(natsConfig.MaxDeliver, int(-1))
	assert.Equal(natsConfig.FilterSubject, "foo")
	assert.Equal(natsConfig.ReplayPolicy, nats.ReplayInstantPolicy)
	assert.Equal(natsConfig.MaxAckPending, int(5000))
	assert.Equal(natsConfig.InactiveThreshold, func() time.Duration {
		t, _ := time.ParseDuration("1m")
		return t
	}())
	assert.Equal(natsConfig.Replicas, 0)
	assert.Equal(natsConfig.MemoryStorage, false)
	assert.Equal(natsConfig.MaxWaiting, int(100))
	assert.Equal(natsConfig.MaxRequestExpires, func() time.Duration {
		t, _ := time.ParseDuration("1m")
		return t
	}())
	assert.Equal(natsConfig.MaxRequestBatch, int(100))
	assert.Equal(natsConfig.MaxRequestMaxBytes, int(1024))
}

func TestConvertToNatsPush(t *testing.T) {
	assert := assert.New(t)

	customConfig := &ConsumerConfig{
		Description:       "my consumer",
		DeliverPolicy:     "New",
		AckPolicy:         "Explicit",
		AckWait:           "2m",
		MaxDeliver:        -1,
		FilterSubject:     "foo",
		ReplayPolicy:      "Instant",
		MaxAckPending:     5000,
		InactiveThreshold: "1m",
		Replicas:          0,
		MemoryStorage:     false,
		PushConsumer: &PushConsumerSpec{
			RateLimit:      1000,
			HeadersOnly:    true,
			DeliverSubject: "mysubject",
			DeliverGroup:   "mydelivergroup",
			FlowControl:    true,
			IdleHeartbeat:  "30s",
		},
	}

	natsConfig, err := ConfigV1Alpha1ToNats("mystream", customConfig)
	assert.Nil(err)
	assert.Equal(natsConfig.Description, "my consumer")
	assert.Equal(natsConfig.DeliverPolicy, nats.DeliverNewPolicy)
	assert.Equal(natsConfig.AckPolicy, nats.AckExplicitPolicy)
	assert.Equal(natsConfig.AckWait, func() time.Duration {
		t, _ := time.ParseDuration("2m")
		return t
	}())
	assert.Equal(natsConfig.MaxDeliver, int(-1))
	assert.Equal(natsConfig.FilterSubject, "foo")
	assert.Equal(natsConfig.ReplayPolicy, nats.ReplayInstantPolicy)
	assert.Equal(natsConfig.MaxAckPending, int(5000))
	assert.Equal(natsConfig.InactiveThreshold, func() time.Duration {
		t, _ := time.ParseDuration("1m")
		return t
	}())
	assert.Equal(natsConfig.Replicas, 0)
	assert.Equal(natsConfig.MemoryStorage, false)
	assert.Equal(natsConfig.RateLimit, uint64(1000))
	assert.Equal(natsConfig.HeadersOnly, true)
	assert.Equal(natsConfig.DeliverSubject, "mysubject")
	assert.Equal(natsConfig.DeliverGroup, "mydelivergroup")
	assert.Equal(natsConfig.FlowControl, true)
	assert.Equal(natsConfig.Heartbeat, func() time.Duration {
		t, _ := time.ParseDuration("30s")
		return t
	}())
}

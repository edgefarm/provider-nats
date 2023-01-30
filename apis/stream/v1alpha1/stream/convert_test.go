package stream

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestConvertToNats(t *testing.T) {
	assert := assert.New(t)
	optStartTime := "2023-01-09T14:48:32Z"
	maxAge := "1h2m3s"
	duplicates := "2m"

	customConfig := &StreamConfig{
		Description: "this is a test stream",
		Subjects: []string{
			"foo",
			"bar",
			"baz.>",
		},
		Retention:            "Limits",
		MaxConsumers:         2,
		MaxMsgs:              100,
		MaxBytes:             1024,
		Discard:              "New",
		DiscardNewPerSubject: false,
		MaxAge:               maxAge,
		MaxMsgsPerSubject:    -1,
		MaxMsgSize:           -1,
		Duplicates:           duplicates,
		Storage:              "File",
		Replicas:             1,
		NoAck:                false,
		Placement:            &Placement{},
		Mirror: &StreamSource{
			Name:      "mymirror",
			StartSeq:  0,
			StartTime: optStartTime,
			External: &ExternalStream{
				APIPrefix:     "$JS.mydomain.API,",
				DeliverPrefix: "",
			},
			Domain: "mydomain",
		},
		Sources:      []*StreamSource{},
		Sealed:       false,
		DenyDelete:   false,
		DenyPurge:    false,
		AllowRollup:  false,
		RePublish:    &RePublish{},
		AllowDirect:  false,
		MirrorDirect: false,
	}

	natsConfig, err := ConfigV1Alpha1ToNats("mystream", customConfig)
	assert.Nil(err)
	assert.Equal(natsConfig.Name, "mystream")
	assert.Equal(natsConfig.Description, "this is a test stream")
	assert.Equal(natsConfig.Subjects, []string{
		"foo",
		"bar",
		"baz.>",
	})
	assert.Equal(natsConfig.Retention, nats.LimitsPolicy)
	assert.Equal(natsConfig.MaxConsumers, 2)
	assert.Equal(natsConfig.MaxMsgs, int64(100))
	assert.Equal(natsConfig.MaxBytes, int64(1024))
	assert.Equal(natsConfig.Discard, nats.DiscardNew)
	assert.Equal(natsConfig.DiscardNewPerSubject, false)
	assert.Equal(natsConfig.MaxAge, func() time.Duration {
		t, _ := time.ParseDuration(maxAge)
		return t
	}())
	assert.Equal(natsConfig.MaxMsgsPerSubject, int64(-1))
	assert.Equal(natsConfig.MaxMsgSize, int32(-1))
	assert.Equal(natsConfig.Duplicates, func() time.Duration {
		t, _ := time.ParseDuration(duplicates)
		return t
	}())
	assert.Equal(natsConfig.Storage, nats.FileStorage)
	assert.Equal(natsConfig.Replicas, 1)
	assert.Equal(natsConfig.NoAck, false)
	assert.Equal(natsConfig.Placement, &nats.Placement{})
	assert.Equal(natsConfig.Mirror, &nats.StreamSource{
		Name:        "mymirror",
		OptStartSeq: 0,
		OptStartTime: func() *time.Time {
			t, _ := time.Parse(time.RFC3339, optStartTime)
			return &t
		}(),
		External: &nats.ExternalStream{
			APIPrefix:     "$JS.mydomain.API,",
			DeliverPrefix: "",
		},
		Domain: "mydomain",
	})
	assert.Equal(natsConfig.Sources, []*nats.StreamSource{})
	assert.Equal(natsConfig.Sealed, false)
	assert.Equal(natsConfig.DenyDelete, false)
	assert.Equal(natsConfig.DenyPurge, false)
	assert.Equal(natsConfig.AllowRollup, false)
	assert.Equal(natsConfig.RePublish, &nats.RePublish{})
	assert.Equal(natsConfig.AllowDirect, false)
	assert.Equal(natsConfig.MirrorDirect, false)
}

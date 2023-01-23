package nats

import (
	"errors"

	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

var (
	ErrNatsConfig = errors.New("secret does not contain a valid nats configuration")
)

type Config struct {
	// JWT is the NATS users JWT token to use for authentication.
	JWT string `json:"jwt"`
	// SeedKey is the NATS users seed key to use for authentication.
	SeedKey string `json:"seed_key"`
	// Address is the NATS address to use for authentication.
	Address string `json:"address"`
}

type Client struct {
	conn          *natsgo.Conn
	Address       string
	UserPublicKey string
}

func GetUserPublicKey(seed string) (string, error) {
	user, err := nkeys.FromSeed([]byte(seed))
	if err != nil {
		return "", err
	}

	pub, err := user.PublicKey()
	if err != nil {
		return "", err
	}
	return pub, nil
}

func NewClient(creds []byte) (*Client, error) {
	var config Config
	if err := jsonutil.DecodeJSON(creds, &config); err != nil {
		return nil, err
	}

	if config.JWT == "" || config.SeedKey == "" || config.Address == "" {
		return nil, ErrNatsConfig
	}

	c, err := natsgo.Connect(config.Address, natsgo.UserJWTAndSeed(config.JWT, config.SeedKey))
	if err != nil {
		return nil, err
	}

	pub, err := GetUserPublicKey(config.SeedKey)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:          c,
		Address:       config.Address,
		UserPublicKey: pub,
	}, nil
}

func (c *Client) Disconnect() error {
	c.conn.Close()
	return nil
}

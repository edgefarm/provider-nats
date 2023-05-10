package nats

import (
	"errors"

	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	natsjwt "github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats.go"
	natsgo "github.com/nats-io/nats.go"
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
	conn             *natsgo.Conn
	Address          string
	UserPublicKey    string
	AccountPublicKey string
}

func GetPublicKeys(jwt string) (string, string, error) {
	c, err := natsjwt.DecodeUserClaims(jwt)
	if err != nil {
		return "", "", err
	}
	return c.Issuer, c.Subject, nil
}

func NewClient(creds []byte) (*Client, error) {
	var config Config
	if err := jsonutil.DecodeJSON(creds, &config); err != nil {
		return nil, err
	}

	if config.JWT == "" || config.SeedKey == "" || config.Address == "" {
		return nil, ErrNatsConfig
	}

	opts := []nats.Option{}
	opts = append(opts, nats.UserJWTAndSeed(config.JWT, config.SeedKey))
	c, err := natsgo.Connect(config.Address, opts...)
	if err != nil {
		return nil, err
	}

	accountPub, userPub, err := GetPublicKeys(config.JWT)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:             c,
		Address:          config.Address,
		UserPublicKey:    userPub,
		AccountPublicKey: accountPub,
	}, nil
}

func (c *Client) Disconnect() error {
	c.conn.Close()
	return nil
}

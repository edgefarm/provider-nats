package nats

import (
	"errors"
	"log"
	"sync"
	"time"

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

var lock = &sync.Mutex{}
var clientInstance *Client

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

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectHandler(func(nc *nats.Conn) {
		log.Printf("Disconnected: will attempt reconnects for %.0fm", totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Printf("Closing connection")
	}))
	return opts
}

func NewClient(creds []byte) (*Client, error) {
	if clientInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if clientInstance == nil {
			var config Config
			if err := jsonutil.DecodeJSON(creds, &config); err != nil {
				return nil, err
			}

			if config.JWT == "" || config.SeedKey == "" || config.Address == "" {
				return nil, ErrNatsConfig
			}

			opts := []nats.Option{}
			opts = setupConnOptions(opts)
			opts = append(opts, nats.UserJWTAndSeed(config.JWT, config.SeedKey))
			c, err := natsgo.Connect(config.Address, opts...)
			if err != nil {
				return nil, err
			}

			accountPub, userPub, err := GetPublicKeys(config.JWT)
			if err != nil {
				return nil, err
			}
			clientInstance = &Client{
				conn:             c,
				Address:          config.Address,
				UserPublicKey:    userPub,
				AccountPublicKey: accountPub,
			}
		}
	}
	return clientInstance, nil
}

func (c *Client) Disconnect() error {
	c.conn.Close()
	return nil
}

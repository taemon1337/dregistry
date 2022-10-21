package server

import (
  "log"
  "time"

	"github.com/hashicorp/serf/client"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	config *client.Config         `json:"config" yaml:"config"`
	c      *client.RPCClient      `json:"-" yaml:"-"`
}

func (c *Client) Connect() error {
  var conn *client.RPCClient
  var err error

	if conn, err = client.ClientFromConfig(c.config); err != nil {
		return err
	}

  log.Println("Connected to %s", c.config.Addr)
	c.c = conn // connected
  return nil
}

func (c *Client) Listen(ch chan map[string]interface{}) {
	for {
		data := <-ch
		log.Printf("PONG: %v\n", data)
	}
}

func (c *Client) Start() error {
	if c.c == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	var g errgroup.Group

  log.Println("Starting serf client")
	g.Go(func() error {
    log.Println("Starting listener...")
		ch := make(chan map[string]interface{})
		c.c.Stream("user", ch)
		c.Listen(ch)
    return nil
	})

  g.Go(func() error {
    for {
      log.Println("PINGING")
      if err := c.c.UserEvent("test", []byte("PING"), false); err != nil {
        return err
      }
      time.Sleep(5 * time.Second)
    }
  })

	// wait for all go routines to complete
	// and return first non-nil error
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func NewClientConfig(addr string) *client.Config {
	return &client.Config{
		Addr: addr,
	}
}

func NewClient(cfg *client.Config) *Client {
	return &Client{
		config: cfg,
		c:      nil,
	}
}

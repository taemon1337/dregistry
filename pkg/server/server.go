package server

import (
  "log"
  "time"

	"github.com/hashicorp/serf/serf"
)

type Server struct {
  srv     *serf.Serf
  cfg     *serf.Config
}

func NewConfig(name string) *serf.Config {
  cfg := serf.DefaultConfig()

  if name != "" {
    cfg.NodeName = name
  }

  return cfg
}

func New(cfg *serf.Config) (*Server, error) {
  var s *serf.Serf
  var err error

  if s, err = serf.Create(cfg); err != nil {
    return nil, err
  }

  return &Server {
    srv: s,
    cfg: cfg,
  }, nil
}

func (s *Server) Start() error {
  log.Println("Starting Serf...")
  for {
    time.Sleep(5 * time.Second)
    resp, err := s.srv.Query("hello", []byte("is anybody out there?")); err != nil {
      log.Printf("error sending hello query - %s", err)
    }

    log.Printf("HELLO QUERY RESPONSE: %s", resp)
  }
}

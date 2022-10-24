package server

import (
  "log"
  "time"
  "sync"
  "context"

	"github.com/hashicorp/serf/serf"
)

type Server struct {
  srv     *serf.Serf
  cfg     *serf.Config
  events  chan serf.Event
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

  eventCh := make(chan serf.Event, 4) // s.srv.EventCh
  cfg.EventCh = eventCh

  if s, err = serf.Create(cfg); err != nil {
    return nil, err
  }

  return &Server {
    srv: s,
    cfg: cfg,
    events: eventCh,
  }, nil
}

func (s *Server) Listen() {
  var wg sync.WaitGroup
  defer wg.Wait()

  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  wg.Add(1)
  go func() {
    defer wg.Done()
    for {
      select {
        case <- ctx.Done():
          return
        case r := <-s.events:
          log.Printf("EVENT: %s", r)
      }
    }
  }()
}

func (s *Server) Start() error {
  log.Println("Start listener...")
  s.Listen()

  log.Println("Starting events...")
  for {
    time.Sleep(3*time.Second)
    if err := s.srv.UserEvent("hello", []byte("is anybody out there?"), true); err != nil {
      log.Printf("[ERROR] error sending user event hello - %s", err)
    }

    log.Printf("SENT 'hello'")
  }
}

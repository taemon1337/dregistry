package server

import (
  "log"
  "fmt"
  "time"
  "sync"
  "errors"
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

  eventCh := make(chan serf.Event, 4)
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

func (s *Server) Start(joins []string, emit bool) error {
  log.Println("Start listener...")
  s.Listen()

  if _, err := s.srv.Join(joins, false); err != nil {
    return errors.New(fmt.Sprintf("could not join - %s", err))
  }

  log.Println("Starting events...")
  for {
    time.Sleep(3*time.Second)
    if emit {
      resp, err := s.srv.Query("ping", []byte("return-me"), &serf.QueryParam{}); 
      if err != nil {
        log.Printf("[ERROR] error sending ping query - %s", err)
      }
  
      log.Printf("[RESP] %s", resp)

      if err := s.srv.UserEvent("hello", []byte("is anybody out there?"), true); err != nil {
        log.Printf("[ERROR] error sending user event hello - %s", err)
      }

      log.Printf("SENT 'hello'")

      for _, member := range s.srv.Members() {
        log.Printf("MEMBER: %s@%s:%d", member.Name, member.Addr, member.Port)
      }
    }
  }
  return nil
}

package agent

import (
  "fmt"
  "log"
  "time"
  "errors"

	"github.com/hashicorp/serf/cmd/serf/command/agent"
	"github.com/hashicorp/serf/serf"
  runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

type Agent struct {
  agent         *agent.Agent
  image         *ImageService
}

func NewAgent(nodename, criaddr string, timeout time.Duration) (*Agent, error) {
  agentConf := agent.DefaultConfig()
  conf := serf.DefaultConfig()

  agentConf.NodeName = nodename
  conf.NodeName = nodename

  a, err := agent.Create(agentConf, conf, nil)
  if err != nil {
    return nil, err
  }

  is, err := NewImageService(criaddr, timeout)
  if err != nil {
    log.Printf("could not get image service at %s - %s", criaddr, err)
    return nil, err
  }

  resp, err := is.ListImages(&runtimeapi.ImageFilter{})
  if err != nil {
    log.Printf("could not list images - %s", err)
  }

  if err == nil {
    for _, img := range resp {
      log.Printf("IMAGE: %s", img)
    }
  }

  return &Agent{
    agent:        a,
    image:        is,
  }, nil
}

func (a *Agent) Start(joins []string, emit bool) error {
  if err := a.agent.Start(); err != nil {
    return errors.New(fmt.Sprintf("could not start agent - %s", err))
  }

  if _, err := a.agent.Join(joins, true); err != nil {
    return errors.New(fmt.Sprintf("could not join cluster - %s", err))
  }

  a.agent.RegisterEventHandler(a)

  log.Printf("%s", a.agent.Stats())

  for {
    time.Sleep(10*time.Second)
    ml := a.agent.Serf().Memberlist()
    if ml.NumMembers() < 1 {
      return nil
    }

    if emit {
      if err := a.agent.UserEvent("ping", []byte("reply with: pong"), false); err != nil {
        return errors.New(fmt.Sprintf("could not send ping - %s", err))
      }
    }
  }

  return nil
}

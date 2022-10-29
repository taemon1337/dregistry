package agent

import (
  "log"

	"github.com/hashicorp/serf/serf"
)

func (a *Agent) HandleEvent(e serf.Event) {
  if e.EventType() == serf.EventUser {
    ue := e.(serf.UserEvent)
    log.Printf("EVENT: %s", ue.Name)
  }
}

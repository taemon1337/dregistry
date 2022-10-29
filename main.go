package main

import (
	"flag"
  "time"
	"github.com/taemon1337/dregistry/pkg/agent"
	"github.com/taemon1337/dregistry/pkg/server"
	"log"
)

func main() {
	var nodename string
  var criaddr string
	var joins server.ArrayFlags
	var emit bool
  var timeout int64

	flag.StringVar(&nodename, "name", "", "The unique node identifier")
	flag.StringVar(&criaddr, "sock", "unix:///var/run/docker.sock", "The container runtime endpoint")
  flag.Int64Var(&timeout, "timeout", 5, "The number of seconds to wait before ending connectiong (timeout)")
	flag.BoolVar(&emit, "emit", false, "Emit events (for testing)")
	flag.Var(&joins, "join", "Join the following node(s)")
	flag.Parse()

	s, err := agent.NewAgent(nodename, criaddr, time.Duration(timeout) * time.Second)
	if err != nil {
		log.Fatal(err)
	}

	s.Start(joins, emit)
}

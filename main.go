package main

import (
  "log"
  "flag"
  "github.com/taemon1337/dregistry/pkg/server"
)

func main() {
  var nodename string

  flag.StringVar(&nodename, "name", "", "The unique node identifier")
  flag.Parse()

  cfg := server.NewConfig(nodename)

  s, err := server.New(cfg);
  if err != nil {
    log.Fatal(err)
  }

  s.Start()
}

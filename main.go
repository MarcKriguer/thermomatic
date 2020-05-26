package main

import (
  "github.com/MarcKriguer/thermomatic/internal/common"
  "github.com/MarcKriguer/thermomatic/internal/server"
)

func main() {
  common.LogOutput("Starting thermomatic service.")
  server.StartServer(common.DefaultTheromaticPort)
}

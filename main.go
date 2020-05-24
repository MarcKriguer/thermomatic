package main

import (
  "github.com/MarcKriguer/thermomatic/internal/common"
  "github.com/MarcKriguer/thermomatic/internal/server"
)

const PORT = 1337

func main() {
  common.LogOutput("Starting thermomatic service.")
  server.StartServer(PORT)
}

// TODO: document the package.
package server

import (
  "net"
  "strconv"
  "github.com/MarcKriguer/thermomatic/internal/common"
)

// This is the handler that is called when a client connects to the server.
func handleConnection(conn net.Conn) {
  common.LogOutput("Connection accepted.")
}

// This function is called to start the socket server on the given port.
func StartServer(port int) {
  common.LogOutput("Starting server on port " + strconv.Itoa(port))

  // Set up the socket and start listening on it.
  link, err := net.Listen("tcp", ":" + strconv.Itoa(port))
  if err != nil {
    common.LogError(err)
    return
  }

  // Make sure the socket is eventually closed.
  defer link.Close()

  // Wait for client connections here. Handle new connections in their own goroutine.
  for {
    conn, err := link.Accept()
    if err != nil {
      common.LogError(err)
      return
    }
    go handleConnection(conn)
  }
}

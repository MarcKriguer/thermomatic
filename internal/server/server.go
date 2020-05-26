// Server package. Handles the Connection, and reading in of logins (IMEI) and Reading packets.
package server

import (
  "errors"
  "fmt"
  "github.com/MarcKriguer/thermomatic/internal/client"
  "github.com/MarcKriguer/thermomatic/internal/common"
  "github.com/MarcKriguer/thermomatic/internal/imei"
  "io"
  "net"
  "strconv"
  "time"
)

var (
  ErrImeiTimeout    = errors.New("server: imei login timeout")
  ErrReadingTimeout = errors.New("server: data reading timeout")
)

// This is the handler that is called when a client connects to the server.
func handleConnection(conn net.Conn) {
  common.LogOutput("Connection accepted.")

  // In case of a panic, recover by closing the connection
  defer func() {
    recover()
    conn.Close()
    return
  }()

  // largest valid message is a Reading
  buffer := make([]byte, client.READING_LENGTH)

  // client has only 1 second to login (send IMEU)
  conn.SetReadDeadline(time.Now().Add(time.Second))

  // read it in
  bytesRead, err := conn.Read(buffer)
  if err != nil {
    common.LogError(err)
    conn.Close()
    return
  }

  // If nothing read in, we've timed out, so display the error and close the connection.
  if bytesRead == 0 {
    common.LogError(ErrImeiTimeout)
    conn.Close()
    return
  }

  // validate the login attempt
  imei, err := imei.Decode(buffer[0:imei.IMEI_LENGTH])
  if err != nil {
    common.LogError(err)
    conn.Close()
    return
  }

  // repeatedly read in next Reading (with a 1-second timeout) and just output it.
  var reading client.Reading

  for {
    // Set the timeout for the data reading to 1 second
    conn.SetReadDeadline(time.Now().Add(1 * time.Second))

    // read in next Reading
    bytesRead, err := conn.Read(buffer)

    // Check the length, if it's 0 then we've timed out and need to close the connection.
    if bytesRead == 0 {
      common.LogError(ErrReadingTimeout)
      conn.Close()
      return
    }
    if err != nil {
      if err == io.EOF {
        common.LogOutput("Connection closed by client!")
        return
      }
    }

    // Decode the Reading
    reading.Decode(buffer[0:client.READING_LENGTH])

    // Log the Reading's data
    common.LogOutput("valid Reading read.\n" + fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v\n", imei,
        time.Now().UnixNano(), reading.Temperature, reading.Altitude, reading.Latitude,
        reading.Longitude, reading.BatteryLevel))
  }
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

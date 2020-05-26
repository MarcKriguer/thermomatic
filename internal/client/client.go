// Client package. Simulates a thermomatic device logging in and sending any number of random
// Readings as desired.
package client

import (
  "github.com/MarcKriguer/thermomatic/internal/common"
  "net"
  "strconv"
  "time"
)

var ValidImei = []byte { 4, 9, 0, 1, 5, 4, 2, 0, 3, 2, 3, 7, 5, 1, 8}

// function to connect to the server, send a number of messages, and close the connection.
// return value is a diagnostic message.
func Connect(imei []byte, imei_timeout_in_millis uint64, reading_timeout_in_millis uint64,
    readings_to_send int) string {
  url := "localhost:" + strconv.Itoa(common.DefaultTheromaticPort)
  conn, err := net.Dial("tcp", url)
  if err != nil {
    common.LogError(err)
    return "Unable to connect: " + err.Error()
  }
  // Set the default timeout for all operations to 5 seconds (login and read time outs are smaller)
  conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

  // Actually do a Sleep (for the specified number of milliseconds) before logging in, in order to
  // test the imei_timeout functionality.
  time.Sleep(time.Duration(imei_timeout_in_millis) * time.Millisecond)

  // login (send the IMEI byte array)
  _, err = conn.Write(imei)
  if err != nil {
    common.LogError(err)
    conn.Close()
    return "Unable to login: " + err.Error()
  }

  var reading Reading
  // send "readings_to_send" readings to the server
  for i := 0; i < readings_to_send; i++ {
    _, err = conn.Write(reading.GenerateRandomReading())
    if err != nil {
      common.LogError(err)
      conn.Close()
      return "Unable to send reading #" + strconv.Itoa(i) + " of " +
          strconv.Itoa(readings_to_send) + ": " + err.Error()
    }

    // Actually do a Sleep (for the specified number of milliseconds) before logging in, in order to
    // test the reading_timeout functionality.
    time.Sleep(time.Duration(reading_timeout_in_millis) * time.Millisecond)
  }

  conn.Close()
  return "OK" // all readings successfully sent.
}

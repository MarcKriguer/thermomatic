package main

import (
  "github.com/MarcKriguer/thermomatic/internal/client"
  "testing"
  "strings"
)

// Start the server (e.g. run "go run server.go") in a different window before running these tests.
func TestClientHappyPath(t *testing.T) {
  // It will send 10 randomly generated Readings. (The timeout delays specified before are all
  // within server limits, so all Readings should get sent successfully.
  result := client.Connect(client.ValidImei, 200, 200, 10)

  if result != "OK" {
    t.Error("Server error: " + result)
  }
}

// This tests specifying an invalid IMEI.
func TestConnectionInvalidImei(t *testing.T) {
  invalidImei := []byte { 4, 9, 0, 1, 5, 4, 2, 0, 3, 2, 3, 7, 5, 1, 9}
  result := client.Connect(invalidImei, 200, 200, 10)

  if !strings.HasPrefix(result, "Unable to send reading #1 ") {
    t.Error("Server error: " + result)
  }
}

// This tests if the IMEI isn't sent within the timeout.
func TestConnectionImeiTimeout(t *testing.T) {

  // Set the IMEI timeout to 5 seconds to force the error
  result := client.Connect(client.ValidImei, 5000, 200, 10)

  if !strings.HasPrefix(result, "Unable to login:") {
    t.Error("Server error: " + result)
  }
}

// This tests if a Reading isn't sent within the timeout.
func TestConnectionReadingTimeout(t *testing.T) {

  // Set the Reading timeout to 5 seconds to force the error
  result := client.Connect(client.ValidImei, 200, 5000, 10)

  if !strings.HasPrefix(result, "Unable to send reading #") {
    t.Error("Server error: " + result)
  }
}

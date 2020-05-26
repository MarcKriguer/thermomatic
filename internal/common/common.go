// Package common implements utilities & constants commonly consumed by the rest of the packages.
package common

import "fmt"
import "os"

// Tcp port to use for the server
var DefaultTheromaticPort = 1337

// Outputs an error message to StdErr
func LogError(input error) {
  os.Stderr.WriteString("ERROR: ")
  fmt.Fprintln(os.Stderr, input)
  os.Stderr.WriteString("\n")
}

// Outputs a server message to StdOut
func LogOutput(input string) {
  os.Stdout.WriteString("Thermomatic: " + input + "\n")
}

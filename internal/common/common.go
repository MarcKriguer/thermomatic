// Package common implements utilities & functionality commonly consumed by the
// rest of the packages.
package common

import "errors"
import "fmt"
import "os"

// ErrNotImplemented is raised throughout the codebase of the challenge to
// denote implementations to be done by the candidate.
var ErrNotImplemented = errors.New("not implemented")

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

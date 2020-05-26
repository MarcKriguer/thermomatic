package client

import (
  "bytes"
  "encoding/binary"
  "errors"
  "fmt"
  "github.com/MarcKriguer/thermomatic/internal/common"
  "math"
  "math/rand"
)

var (
  ErrReadingLength = errors.New("client: Reading byte array is less than 40 characters long.")
)

const READING_LENGTH = 40

// Reading is the set of device readings.
type Reading struct {
	// Temperature denotes the temperature reading of the message.  Valid range [-300, 300].
	Temperature float64

	// Altitude denotes the altitude reading of the message.  Valid range [-20000, 20000].
	Altitude float64

	// Latitude denotes the latitude reading of the message.  Valid range [-90, 90].
	Latitude float64

	// Longitude denotes the longitude reading of the message.  Valid range [-180, 180].
	Longitude float64

	// BatteryLevel denotes the battery level reading of the message.  Valid range (0, 100].
	BatteryLevel float64
}

// Decode the reading message payload in the given byte array into a Reading.
//
// Returns true if all fields are within min/max ranges, false if any are outside their range.
//
// Decode does NOT allocate under any condition.
// Additionally, it panics if b isn't at least 40 bytes long.
func (r *Reading) Decode(b []byte) (ok bool) {
  // panic if byte array is too small
  if len(b) < READING_LENGTH {
    common.LogError(ErrReadingLength)
    panic(ErrReadingLength)
  }

  // extract each field
  r.Temperature  = math.Float64frombits(binary.BigEndian.Uint64(b[0:8]))
  r.Altitude     = math.Float64frombits(binary.BigEndian.Uint64(b[8:16]))
  r.Latitude     = math.Float64frombits(binary.BigEndian.Uint64(b[16:24]))
  r.Longitude    = math.Float64frombits(binary.BigEndian.Uint64(b[24:32]))
  r.BatteryLevel = math.Float64frombits(binary.BigEndian.Uint64(b[32:40]))

  // validate each field, returning false if any of them are outside their range
  if r.Temperature > 300 || r.Temperature < -300 {
    common.LogOutput("Temperature out of range: " + fmt.Sprintf("%f", r.Temperature))
    return false
  }
  if r.Altitude > 20000 || r.Altitude < -20000 {
    common.LogOutput("Altitude out of range: " + fmt.Sprintf("%f", r.Altitude))
    return false
  }
  if r.Latitude > 90 || r.Latitude < -90 {
    common.LogOutput("Latitude out of range: " + fmt.Sprintf("%f", r.Latitude))
    return false
  }
  if r.Longitude > 180 || r.Longitude < -180 {
    common.LogOutput("Longitude out of range: " + fmt.Sprintf("%f", r.Longitude))
    return false
  }
  if r.BatteryLevel > 100 || r.BatteryLevel <= 0 {
    common.LogOutput("BatteryLevel out of range: " + fmt.Sprintf("%f", r.BatteryLevel))
    return false
  }

  return true
}

// Encode encodes the reading message payload in the given r into a byte array.
func (r *Reading) Encode() (buf []byte) {
  var buffer []byte
  buffer = append(buffer, Float64ToByteArray(r.Temperature)...)
  buffer = append(buffer, Float64ToByteArray(r.Altitude)...)
  buffer = append(buffer, Float64ToByteArray(r.Latitude)...)
  buffer = append(buffer, Float64ToByteArray(r.Longitude)...)
  buffer = append(buffer, Float64ToByteArray(r.BatteryLevel)...)
  return buffer
}

// Generate an 8-byte array containing a 64-bit float
func Float64ToByteArray(field float64) (buf []byte) {
  buffer := new (bytes.Buffer)
  err := binary.Write(buffer, binary.BigEndian, field)
  if err != nil { 
    common.LogError(err)
  }
  return buffer.Bytes()
}

// Generate a byte array of a valid Reading [this method only used by test methods]
func (r *Reading) GenerateRandomReading() (b []byte) {
  // Generate a Reading and populate with random data
  var reading Reading
  reading.Temperature = (rand.Float64() * 600) - 300
  reading.Altitude = (rand.Float64() * 40000) - 20000
  reading.Latitude = (rand.Float64() * 180) - 90
  reading.Longitude = (rand.Float64() * 360) - 180
  reading.BatteryLevel = (rand.Float64() * 100)
  // reroll BatteryLevel while it's 0
  for reading.BatteryLevel == 0 {
    reading.BatteryLevel = (rand.Float64() * 100)
  }

  return reading.Encode()
}

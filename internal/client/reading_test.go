package client

import (
  "bytes"
  "testing"
)

// Test the Decode function, when exactly 40 bytes are passed in.
func TestReadingDecodeExact40(t *testing.T) {
  var reading Reading

  // Declare byte array to Decode
  byteArray := []uint8{
    0x40, 0x50, 0xf1, 0x47, 0xae, 0x14, 0x7a, 0xe1, // temperature
    0x40, 0x05, 0x15, 0x9b, 0x3d, 0x07, 0xc8, 0x4b, // altitude
    0x40, 0x40, 0xb4, 0x7a, 0xe1, 0x47, 0xae, 0x14, // latitude
    0x40, 0x46, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // longitude
    0x3f, 0xd0, 0x6d, 0x1e, 0x10, 0x8c, 0x3f, 0x3e, // battery level
  }

  result := reading.Decode(byteArray)
  if !result {
    t.Errorf("Failed to decode byte array")
  }

  if reading.Temperature != 67.77 {
    t.Errorf("Error decoding Temperature (found %x)", reading.Temperature)
  }
  if reading.Altitude != 2.63555 {
    t.Errorf("Error decoding Altitude (found %x)", reading.Altitude)
  }
  if reading.Latitude != 33.41 {
    t.Errorf("Error decoding Latitude (found %x)", reading.Latitude)
  }
  if reading.Longitude != 44.4 {
    t.Errorf("Error decoding Longitude (found %x)", reading.Longitude)
  }
  if reading.BatteryLevel != 0.25666 {
    t.Errorf("Error decoding BatteryLevel (found %x)", reading.BatteryLevel)
  }
}

// Test the Decode function, for a randomly-generated Reading.
func TestReadingDecodeRandom(t *testing.T) {
  var reading Reading
  result := reading.Decode(reading.GenerateRandomReading())
  if !result {
    t.Errorf("Failed to decode random byte array")
  }
}

// Test the Decode function, when more than 40 bytes are passed in.
func TestReadingDecodeOver40(t *testing.T) {
  var reading Reading

  // Declare byte array to Decode
  byteArray := []uint8{
    0x40, 0x50, 0xf1, 0x47, 0xae, 0x14, 0x7a, 0xe1, // temperature
    0x40, 0x05, 0x15, 0x9b, 0x3d, 0x07, 0xc8, 0x4b, // altitude
    0x40, 0x40, 0xb4, 0x7a, 0xe1, 0x47, 0xae, 0x14, // latitude
    0x40, 0x46, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // longitude
    0x3f, 0xd0, 0x6d, 0x1e, 0x10, 0x8c, 0x3f, 0x3e, // battery level
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // ignored
  }

  result := reading.Decode(byteArray)
  if !result {
    t.Errorf("Failed to decode larger byte array")
  }

  if reading.Temperature != 67.77 {
    t.Errorf("Error decoding Temperature (found %x)", reading.Temperature)
  }
  if reading.Altitude != 2.63555 {
    t.Errorf("Error decoding Altitude (found %x)", reading.Altitude)
  }
  if reading.Latitude != 33.41 {
    t.Errorf("Error decoding Latitude (found %x)", reading.Latitude)
  }
  if reading.Longitude != 44.4 {
    t.Errorf("Error decoding Longitude (found %x)", reading.Longitude)
  }
  if reading.BatteryLevel != 0.25666 {
    t.Errorf("Error decoding BatteryLevel (found %x)", reading.BatteryLevel)
  }
}

// Test that the Decode function panics when less than 40 bytes are passed in.
func TestReadingDecodeUnder40(t *testing.T) {
  var reading Reading

  // Declare byte array to Decode
  byteArray := []uint8{
    0x40, 0x50, 0xf1, 0x47, 0xae, 0x14, 0x7a, 0xe1, // temperature
    0x40, 0x05, 0x15, 0x9b, 0x3d, 0x07, 0xc8, 0x4b, // altitude
    0x40, 0x40, 0xb4, 0x7a, 0xe1, 0x47, 0xae, 0x14, // latitude
    0x40, 0x46, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // longitude
    0x3f, 0xd0, 0x6d, 0x1e, 0x10, 0x8c, 0x3f,       // battery level (missing final byte)
  }

  defer func() {
   if r := recover(); r == nil {
      t.Errorf("Reading.Decode did not panic when expected")
    }
  }()

  result := reading.Decode(byteArray)
  if !result {
    t.Errorf("Error (but not expected panic) when calling Reading.Decode with too few bytes.")
  }
}

// Test that the Decode function returns false (but doesn't panic) when a field is out of range
func TestReadingDecodeWithInvalidBatteryLevel(t *testing.T) {
  var reading Reading

  // Declare byte array to Decode
  byteArray := []uint8{
    0x40, 0x50, 0xf1, 0x47, 0xae, 0x14, 0x7a, 0xe1, // temperature
    0x40, 0x05, 0x15, 0x9b, 0x3d, 0x07, 0xc8, 0x4b, // altitude
    0x40, 0x40, 0xb4, 0x7a, 0xe1, 0x47, 0xae, 0x14, // latitude
    0x40, 0x46, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // longitude
    0xc0, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // invalid BatteryLevel of -5
  }

  defer func() {
   if r := recover(); r != nil {
      t.Errorf("Reading.Decode panicked when it just should have returned false")
    }
  }()

  result := reading.Decode(byteArray)
  if result {
    t.Errorf("Did not get a false (error) result when calling Reading.Decode with invalid data.")
  }
}

// This tests the Encode function (used by the client to test the server).
func TestReadingEncode(t *testing.T) {
  // Declare a Reading struc and stuff in the values from the provided assignment's format example..
  var reading Reading
  reading.Temperature = 67.77
  reading.Altitude = 2.63555
  reading.Latitude = 33.41
  reading.Longitude = 44.4
  reading.BatteryLevel = 0.25666

  // Encode to get the byte array
  encoded := reading.Encode()

  // Declare byte arrays for comparison
  temperature := []uint8{0x40, 0x50, 0xf1, 0x47, 0xae, 0x14, 0x7a, 0xe1}
  altitude    := []uint8{0x40, 0x05, 0x15, 0x9b, 0x3d, 0x07, 0xc8, 0x4b}
  latitude    := []uint8{0x40, 0x40, 0xb4, 0x7a, 0xe1, 0x47, 0xae, 0x14}
  longitude   := []uint8{0x40, 0x46, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33}
  battery     := []uint8{0x3f, 0xd0, 0x6d, 0x1e, 0x10, 0x8c, 0x3f, 0x3e}

  // Compare each element of the Reading with the expected Byte arrays
  if !bytes.Equal([]byte(temperature), []byte(encoded[0:8])) {
    t.Errorf("Error in encoding Temperature (was %x)", encoded[0:8])
  }
  if !bytes.Equal([]byte(altitude), []byte(encoded[8:16])) {
    t.Errorf("Error in encoding Altitude (was %x)", encoded[8:16])
  }
  if !bytes.Equal([]byte(latitude), []byte(encoded[16:24])) {
    t.Errorf("Error in encoding Latitude (was %x)", encoded[16:24])
  }
  if !bytes.Equal([]byte(longitude), []byte(encoded[24:32])) {
    t.Errorf("Error in encoding Longitude (was %x)", encoded[24:32])
  }
  if !bytes.Equal([]byte(battery), []byte(encoded[32:40])) {
    t.Errorf("Error in encoding BatteryLevel (was %x)", encoded[32:40])
  }
}

func BenchmarkReadingDecode(b *testing.B) {
  b.ReportAllocs()
  var reading Reading

  // Declare byte array to Decode
  byteArray := []uint8{
    0x40, 0x50, 0xf1, 0x47, 0xae, 0x14, 0x7a, 0xe1, // temperature
    0x40, 0x05, 0x15, 0x9b, 0x3d, 0x07, 0xc8, 0x4b, // altitude
    0x40, 0x40, 0xb4, 0x7a, 0xe1, 0x47, 0xae, 0x14, // latitude
    0x40, 0x46, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // longitude
    0x3f, 0xd0, 0x6d, 0x1e, 0x10, 0x8c, 0x3f, 0x3e, // battery level
  }

  for i := 0; i < b.N; i++ {
    result := reading.Decode(byteArray)
    if !result {
      b.Errorf("Failed to decode byte array in Benchmark.")
    }
  
    if reading.Temperature != 67.77 {
      b.Errorf("Error decoding Temperature (found %x)", reading.Temperature)
    }
    if reading.Altitude != 2.63555 {
      b.Errorf("Error decoding Altitude (found %x)", reading.Altitude)
    }
    if reading.Latitude != 33.41 {
      b.Errorf("Error decoding Latitude (found %x)", reading.Latitude)
    }
    if reading.Longitude != 44.4 {
      b.Errorf("Error decoding Longitude (found %x)", reading.Longitude)
    }
    if reading.BatteryLevel != 0.25666 {
      b.Errorf("Error decoding BatteryLevel (found %x)", reading.BatteryLevel)
    }
  }
}

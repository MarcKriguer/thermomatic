package imei

import (
	"github.com/MarcKriguer/thermomatic/internal/imei"
	"testing"
)

func TestDecode(t *testing.T) {
  // Happy path: make sure valid IMEI is decoded properly
  var validimei = [...]byte { 4, 9, 0, 1, 5, 4, 2, 0, 3, 2, 3, 7, 5, 1, 8}
  result, err := imei.Decode(validimei[:])
  if err != nil {
    t.Errorf("Calling Decode received unexpected error.")
  }

  // make sure it was decoded properly
  if result != 490154203237518 {
    t.Errorf("Received unexpected result")
  }
}

func TestDecodePanics(t *testing.T) {
  // unhappy path 1: invalid IMEI length
  var invalidimei = [...]byte { 4, 9, 0, 1, 5, 4, 2, 0, 3, 2, 3, 7, 5, 1}

  defer func() {
    if r := recover(); r == nil {
      t.Errorf("Decode did not panic when expected")
    }
  }()

  result, err := imei.Decode(invalidimei[:])
  if (err != nil) {
    t.Errorf("Decode returned an error instead of panicing.")
  }
  if (result != 0) {
    t.Errorf("Decode returned a value instead of panicing.")
  }
}

func TestDecodeBadChecksum(t *testing.T) {
  // unhappy path 2: invalid IMEI checksum
  var invalidimei = [...]byte { 4, 9, 0, 1, 5, 4, 2, 0, 3, 2, 3, 7, 5, 1, 9}

  defer func() {
    if r := recover(); r != nil {
      t.Errorf("Decode unexpectedly paniced")
    }
  }()

  result, err := imei.Decode(invalidimei[:])
  if err != imei.ErrChecksum {
      t.Errorf("Decode received a different error than expected.");
  }
  if (result != 0) {
    t.Errorf("Decode returned a value instead of an error.")
  }
}

func TestDecodeBadCharacter(t *testing.T) {
  // unhappy path 3: invalid IMEI character
  var invalidimei = [...]byte { 52, 9, 0, 1, 5, 4, 2, 0, 3, 2, 3, 7, 5, 1, 8}

  defer func() {
    if r := recover(); r != nil {
      t.Errorf("Decode unexpectedly paniced")
    }
  }()

  result, err := imei.Decode(invalidimei[:])
  if err != imei.ErrInvalid {
      t.Errorf("Decode received a different error than expected.");
  }
  if (result != 0) {
    t.Errorf("Decode returned a value instead of an error.")
  }
}

func BenchmarkDecode(b *testing.B) {
  b.ReportAllocs()
  var validimei = [...]byte { 4, 9, 0, 1, 5, 4, 2, 0, 3, 2, 3, 7, 5, 1, 8}
  for i := 0; i < b.N; i++ {
    result, err := imei.Decode(validimei[:])
    if err != nil {
      b.Error("Calling Decode received unexpected error in Benchmark.")
    }

    // make sure it was decoded properly
    if result != 490154203237518 {
      b.Error("Received unexpected result in Benchmark.")
    }
  }
}

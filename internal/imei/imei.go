// Package imei implements an IMEI decoder.
package imei

// NOTE: for more information about IMEI codes and their structure you may
// consult with:
//
// https://en.wikipedia.org/wiki/International_Mobile_Equipment_Identity.

import (
  "errors"
  "github.com/MarcKriguer/thermomatic/internal/common"
)

var (
  ErrChecksum = errors.New("imei: invalid IMEI checksum")
  ErrImeiSize = errors.New("imei: invalid IMEI size (must be 15 characters)")
  ErrInvalid  = errors.New("imei: invalid IMEI character(s)")
)

const IMEI_LENGTH = 15

// Decode returns the IMEI code contained in b.
//
// If b isn't exactly 15 bytes long, the returned error will be ErrImeiSize.
//
// In case b isn't strictly composed of digits, the returned error will be ErrInvalid.
//
// In case b's checksum is wrong, the returned error will be ErrChecksum.
//
// Decode does NOT allocate under any condition. Additionally, it panics if b
// isn't exactly 15 bytes long.
func Decode(b []byte) (code uint64, err error) {
  // make sure length is correct
  if (len(b) != IMEI_LENGTH) {
    common.LogError(ErrImeiSize)
    panic(ErrImeiSize)
  }

  // go through each byte of b -- build up checksum and code, byte by byte (return an error, but
  // not panic, if either the checksum is wrong or an invalid byte is encountered).
  var checksum uint8  = 0
  var results  uint64 = 0
  var digit    uint8

  for i := 0; i < IMEI_LENGTH; i++ {
    // make sure the byte represents a digit
    digit = uint8(b[i])
    if digit > 9 {
      common.LogError(ErrInvalid)
      return 0, ErrInvalid
    }

    results *= uint64(10)
    results += uint64(digit)
    // checksum adds the value of even-positioned digits,
    // and double the value of odd-positioned digits (plus 1 if a 2-digit result)
    if (i % 2 == 0) {
      checksum += digit
    } else {
      checksum += (digit * 2)
      if (digit > 4) {
        checksum++
      }
    }
  }

  // checksum needs to end with a zero (after all bytes have been counted) to be valid
  if checksum % 10 != 0 {
    common.LogError(ErrChecksum)
    return 0, ErrChecksum
  }
  return results, nil
}

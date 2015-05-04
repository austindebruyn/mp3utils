package pickbits

import (
  "fmt"
  "errors"
)

// Helper function
func PickBits(bytes []byte, offset, count int) (uint, error) {
  if offset+count > len(bytes)*8 {
    format := "offset %d and count %d larger than %d bit slice"
    message := fmt.Sprintf(format, offset, count, len(bytes))
    return 0, errors.New(message)
  }

  // Determine if we are split on byte boundary
  var iStart uint = uint(offset/8)
  var iEnd uint = uint((offset+count-1)/8)

  if iStart == iEnd {
    var mask uint8 = (1 << uint(8 - offset)) - 1
    mask -= (1 << uint(8 - offset - count)) - 1
    val := mask & bytes[iStart]
    val = val >> uint(8 - offset - count)
    return uint(val), nil
  } else {
    // Grab val from start byte
    var val uint8
    bitsStart := uint(8 - offset)
    {
      var mask uint8
      if bitsStart == 8 {
        mask = 0xff
      } else {
        mask = (1 << bitsStart) - 1
      }
      val = mask & bytes[iStart]
    }
    startByteValue := uint32(val) << (uint(count)-bitsStart)

    // Grab val from end byte
    bitsEnd := uint(count) - 8*(iEnd-iStart-1) - bitsStart
    {
      var mask uint8 = ^(uint8(0xff) >> bitsEnd)
      val = mask & bytes[iEnd]
      val = val >> uint(8 - bitsEnd)
    }
    endByteValue := uint32(val)

    // Grab vals from middle bytes if there are any
    var middleBytesValue uint32
    if iEnd-iStart > 1 {
      numberOfMiddleBytes := iEnd-iStart-1
      for i := uint(0); i < numberOfMiddleBytes; i++ {
        val = bytes[iStart + i + 1]
        middleBytesValue |= uint32(val) << (bitsEnd + uint((numberOfMiddleBytes-i-1)*8))
      }
    }

    final := startByteValue | middleBytesValue | endByteValue
    return uint(final), nil
  }
}

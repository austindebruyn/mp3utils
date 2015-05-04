package mp3
import (
  //"errors"
  "fmt"
  "io"
  "github.com/austindebruyn/mp3utils/mp3/utils"
)

type SideData struct {
  bytes []byte // private
  ShowBytes []byte
  MainDataPtr uint
  Chunks []ChunkMetadata
}

type ChunkMetadata struct {
  Size uint
  BigValues uint
  GlobalGain uint
  ScaleFactor uint
  WindowSwitching bool
}

// Reads the side data chunk from the file
func ReadSideData(header FrameHeader, file io.Reader) SideData {
  size := header.GetSideDataLength()
  array := make([]byte, size)
  file.Read(array)
  isMono := header.GetChannelMode() == MPEG_CM_Mono
  sideData := ParseSideData(array, isMono)
  return sideData
}

// Parses side information from a slice of raw bytes
func ParseSideData(bytes []byte, isMono bool) SideData {
  var side SideData

  side.ShowBytes = bytes

  // Parse out the pointer to the main data section
  var parsed uint32
  parsed |= uint32(bytes[0]) << 1
  parsed |= uint32(bytes[1]) & 0x80
  side.MainDataPtr = uint(parsed)

  // Dual and mono channel side data is processed much differently
  if isMono {
    panic("cant do mono yet!")
  } else {
    // Dual channel side data
    side.MainDataPtr, _ = pickbits.PickBits(bytes, 0, 9)

    var Chunk1L ChunkMetadata
    Chunk1L.Size, _ = pickbits.PickBits(bytes, 20, 12)
    Chunk1L.BigValues, _ = pickbits.PickBits(bytes, 44, 9)
    Chunk1L.GlobalGain, _ = pickbits.PickBits(bytes, 62, 8)
    Chunk1L.ScaleFactor, _ = pickbits.PickBits(bytes, 78, 4)
    bit, _ := pickbits.PickBits(bytes, 86, 1)
    Chunk1L.WindowSwitching = (bit > 0)
    var Chunk1R ChunkMetadata
    Chunk1R.Size, _ = pickbits.PickBits(bytes, 32, 12)
    Chunk1R.BigValues, _ = pickbits.PickBits(bytes, 53, 9)
    Chunk1R.GlobalGain, _ = pickbits.PickBits(bytes, 70, 8)
    Chunk1R.ScaleFactor, _ = pickbits.PickBits(bytes, 82, 4)
    bit, _ = pickbits.PickBits(bytes, 87, 1)
    Chunk1R.WindowSwitching = (bit > 0)

    side.Chunks = []ChunkMetadata{Chunk1L, Chunk1R}
  }

  return side
}

func (chunk ChunkMetadata) String() string {
  output := "ChunkMetadata{"
  output += fmt.Sprintf("Size: %0x, ", chunk.Size)
  output += fmt.Sprintf("BigValues: %d, ", chunk.BigValues)
  output += fmt.Sprintf("GlobalGain: %d, ", chunk.GlobalGain)
  output += fmt.Sprintf("ScaleFactor: %d, ", chunk.ScaleFactor)
  output += fmt.Sprintf("WindowSwitching: %v }", chunk.WindowSwitching)
  return output
}

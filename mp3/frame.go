package mp3

import (
  "fmt"
  "errors"
)

type Frame struct {
  Header FrameHeader
  Side   SideData
  Data   FrameData
  CRC    uint16
  HasCRC bool
}

var FirstFrame bool = true
var PreviousFrameMainDataEnd int

// Reads an entire MP3 frame
func ReadFrame(bytes []byte, offset int) (Frame, error) {
  var frame Frame

  header, err := ParseHeader(bytes[:4])
  if err != nil {
    return frame, err
  }

  if !header.IsValid() {
    return frame, errors.New("eof")
  }

  // Parse a CRC if it exists
  if header.HasCRC() {
    frame.HasCRC = true
    crcbytes := bytes[2:4]
    frame.CRC = (uint16(crcbytes[0]) << 8) | uint16(crcbytes[1])
  }

  side := ReadSideData(header, bytes, 0)
  data, err := ReadFrameData(header, bytes, 0)
  if err != nil {
    return frame, err
  }

  if FirstFrame {
    FirstFrame = false
    fmt.Println("first frame")
  } else {
    // Attempt to parse out main data
    StartSeek1 := offset - PreviousFrameMainDataEnd
    EndSeek1 := offset
    StartSeek2 := offset + 4 + header.GetSideDataLength()
    fl, _ := header.GetFrameLength()
    EndSeek2 := offset + int(fl) - int(side.MainDataPtr)

    if StartSeek1 >= 0 {
      // Only parse real things...
      DataLength := EndSeek1 + EndSeek2 - StartSeek1 - StartSeek2
      fmt.Printf("Main data is contained [%d,%d] [%d, %d] Length: %d\n", StartSeek1, EndSeek1, StartSeek2, EndSeek2, DataLength)
      MainData := append(bytes[StartSeek1:EndSeek2], bytes[StartSeek2:EndSeek2]...)
      fmt.Printf("hallelujah: %v\n", MainData)
    }
  }
  PreviousFrameMainDataEnd = int(side.MainDataPtr)

  frame.Header = header
  frame.Side = side
  frame.Data = data
  return frame, nil
}

// Prints the frame for debugging
func (frame Frame) String() string {

  var output string

  fr := frame.Header
  side := frame.Side
  data := frame.Data

  output += fmt.Sprintf("FRAME [header: 0x%0x]\n", fr)
  if (frame.HasCRC) {
    output += fmt.Sprintf("CRC: 0x%0x\n", frame.CRC)
  }
  version := fr.GetVersion()
  channel := fr.GetChannelMode()
  layer := fr.GetLayer()
  output += fmt.Sprintf("Version:%d ChannelMode:%d Layer:%d ", version,
    channel, layer)

  bitrate, err := fr.GetBitrate()
  if err != nil {
    return "bad frame"
  }
  output += fmt.Sprintf("Bitrate:%dkbps ", bitrate)

  sr, err := fr.GetSampleRate()
  if err != nil {
    return "bad frame"
  }
  output += fmt.Sprintf("SampleRate:%dhz ", sr)

  fl, err := fr.GetFrameLength()
  if err != nil {
    return "bad frame"
  }
  output += fmt.Sprintf("FrameLength:%d ", fl)
  output += fmt.Sprintf("SideDataLength:%d\n", fr.GetSideDataLength())

  output += fmt.Sprintf("side data: %v\n", side)
  output += fmt.Sprintf("data: %v\n", data)

  return output
}

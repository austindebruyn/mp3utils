package mp3

import (
  "fmt"
  "io"
  "errors"
)

type Frame struct {
  Header FrameHeader
  Side   SideData
  Data   FrameData
  CRC    uint16
  HasCRC bool
}

// Reads an entire MP3 frame
func ReadFrame(file io.Reader) (Frame, error) {
  var frame Frame

  header, err := ReadFrameHeader(file)
  if err != nil {
    return frame, err
  }

  if !header.IsValid() {
    return frame, errors.New("eof")
  }

  // Parse a CRC if it exists
  if header.HasCRC() {
    frame.HasCRC = true
    var crcbytes []byte = make([]byte, 2)
    file.Read(crcbytes)
    frame.CRC = (uint16(crcbytes[0]) << 8) | uint16(crcbytes[1])
  }

  side := ReadSideData(header, file)
  data, err := ReadFrameData(header, file)
  if err != nil {
    return frame, err
  }

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

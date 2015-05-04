package mp3

import (
  "fmt"
  "errors"
  "io"
)

type Chunk []byte
type FrameData []Chunk

// Determines length of the frame data and reads it into a slice
func ReadFrameData(header FrameHeader, file io.Reader) (FrameData, error) {
  frameLength, err := header.GetFrameLength()
  if err != nil {
    return nil, err
  }

  sideDataLength := header.GetSideDataLength()

  dataLength := frameLength - sideDataLength - 4
  if dataLength < 1 || dataLength > 1024 {
    message := fmt.Sprintf("calc error: bad data length %d", dataLength)
    return nil, errors.New(message)
  }
  dataBytes := make([]byte, dataLength)

  file.Read(dataBytes)
  if header.GetChannelMode() == MPEG_CM_Mono {
    // mono file, 2 chunks per frame
    length := len(dataBytes)
    frameData := []Chunk{dataBytes[:length/2], dataBytes[length/2:]}
    return FrameData(frameData), nil
  } else {
    // stereo file, 4 chunks per frame
    length := len(dataBytes)
    frameData := []Chunk{dataBytes[:length/4], dataBytes[length/4:length/2],
      dataBytes[length/2:3*length/4], dataBytes[3*length/4:]}
    return FrameData(frameData), nil
  }
}

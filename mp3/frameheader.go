package mp3

import (
  "errors"
)

const (
  MPEGVersion1 = iota
  MPEGVersion2
  MPEGVersionReserved
  MPEGVersion2_5
)

const (
  MPEG_CM_Stereo = iota
  MPEG_CM_JointStereo
  MPEG_CM_Dual
  MPEG_CM_Mono
)

const (
  MPEGLayerReserved = iota
  MPEGLayer3
  MPEGLayer2
  MPEGLayer1
)

type FrameHeader uint32

var Bitrates = [16][]int{}
// Sets up the bitrate table
func init() {
  // -1 is free, 0 is bad
  Bitrates[0] = []int{-1, -1, -1, -1, -1}
  Bitrates[1] = []int{32, 32, 32, 32, 8}
  Bitrates[2] = []int{64, 48, 40, 48, 16}
  Bitrates[3] = []int{96, 56, 48, 56, 24}
  Bitrates[4] = []int{128, 64, 56, 64, 32}
  Bitrates[5] = []int{160, 80, 64, 80, 40}
  Bitrates[6] = []int{192, 96, 80, 96, 48}
  Bitrates[7] = []int{224, 112, 96, 112, 56}
  Bitrates[8] = []int{256, 128, 112, 128, 64}
  Bitrates[9] = []int{288, 160, 128, 144, 80}
  Bitrates[10] = []int{320, 192, 160,	160, 96}
  Bitrates[11] = []int{352, 224, 192,	176, 112}
  Bitrates[12] = []int{384, 256, 224,	192, 128}
  Bitrates[13] = []int{416, 320, 256,	224, 144}
  Bitrates[14] = []int{448, 384, 320,	256, 160}
  Bitrates[15] = []int{0, 0, 0, 0, 0}
}

// Parses a 4-byte slice into an MPEG frame header
func ParseHeader(bytes []byte) (FrameHeader, error) {

  if len(bytes) != 4 {
    return 0, errors.New("incorrect length")
  }

  var header FrameHeader
  header |= FrameHeader(bytes[0]) << 24
  header |= FrameHeader(bytes[1]) << 16
  header |= FrameHeader(bytes[2]) << 8
  header |= FrameHeader(bytes[3])

  return header, nil
}

// Returns whether or not the MPEG frame looks valid
func (header FrameHeader) IsValid() bool {
  if (header & 0xffe00000) == 0xffe00000 {
    return true
  }
  return false
}

// Returns the MPEG version number
func (header FrameHeader) GetVersion() int {
  switch (header & 0x00180000) >> 19 {
  case 0:
    return MPEGVersion2_5
  case 1:
    return MPEGVersionReserved
  case 2:
    return MPEGVersion2
  case 3:
    return MPEGVersion1
  }
  panic("shouldn't be here")
}

// Returns the MPEG layer number
func (header FrameHeader) GetLayer() int {
  switch (header & 0x00060000) >> 17 {
  case 0:
    return MPEGLayerReserved
  case 1:
    return MPEGLayer3
  case 2:
    return MPEGLayer2
  case 3:
    return MPEGLayer1
  }
  panic("shouldn't be here")
}

// Determines if there is a CRC block for this frame
func (header FrameHeader) HasCRC() bool {
  value := header & 0x00010000
  if value == 0 {
    return true
  }
  return false
}

// Returns the bitrate in kbps
func (header FrameHeader) GetBitrate() (int, error) {
  layer := header.GetLayer()
  value := (header & 0x0000f000) >> 12

  if version := header.GetVersion(); version == MPEGVersion1 {
    // V1,L1  V1,L2  V1,L3
    switch (layer) {
    case MPEGLayer1:
      return Bitrates[value][0], nil
    case MPEGLayer2:
      return Bitrates[value][1], nil
    case MPEGLayer3:
      return Bitrates[value][2], nil
    }
  } else {
    // V2,L1  V2,L2/L3
    if layer == MPEGLayer1 {
      return Bitrates[value][3], nil
    }
    return Bitrates[value][4], nil
  }
  return 0, errors.New("unspecified bitrate")
}

// Return the sample rate in hz
func (header FrameHeader) GetSampleRate() (int, error) {
  value := (header & 0x00000c00) >> 10

  if value < 3 {
    switch header.GetVersion() {
    case MPEGVersion1:
      choices := []int{44100, 48000, 32000}
      return choices[value], nil
    case MPEGVersion2:
      choices := []int{22050, 24000, 16000}
      return choices[value], nil
    case MPEGVersion2_5:
      choices := []int{11025, 11000, 8000}
      return choices[value], nil
    }
  }
  return 0, errors.New("unspecified sample rate")
}

// Return the channel mode
func (header FrameHeader) GetChannelMode() int {
  switch (header & 0x000000c0) >> 6 {
  case 0:
    return MPEG_CM_Stereo
  case 1:
    return MPEG_CM_JointStereo
  case 2:
    return MPEG_CM_Dual
  case 3:
    return MPEG_CM_Mono
  }
  panic("shouldn't be here")
}

// Returns the size of this frame in bytes
func (header FrameHeader) GetFrameLength() (int, error) {
  var SamplesPerFrame int
  if header.GetVersion() == MPEGVersion1 {
    SamplesPerFrame = 1152
  } else {
    SamplesPerFrame = 576
  }
  var BitsPerFrame int = SamplesPerFrame / 8

  sr, err := header.GetSampleRate()
  if err != nil {
    return 0, err
  }

  br, err := header.GetBitrate()
  if err != nil {
    return 0, err
  }

  // Detect padding bit
  extraPadding := 0
  if (header & 0x00000200) == 0x00000200 {
    if header.GetLayer() == MPEGLayer1 {
      extraPadding = 4
    } else {
      extraPadding = 1
    }
  }

  var size int = 1000*br*BitsPerFrame/sr + extraPadding
  return size, nil
}

// Returns the size of the side data in bytes
func (header FrameHeader) GetSideDataLength() int {
  if header.GetChannelMode() < 3 {
    return 32
  }
  return 17
}

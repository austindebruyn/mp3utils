package pickbits

import (
  "testing"
)

func TestOutOfRange(t *testing.T) {
  if _, err := PickBits(make([]byte, 0), 1, 1); err == nil {
    t.Fail()
  }
  if _, err := PickBits(make([]byte, 2), 1, 16); err == nil {
    t.Fail()
  }
  if _, err := PickBits(make([]byte, 10), 79, 2); err == nil {
    t.Fail()
  }
  if _, err := PickBits(make([]byte, 1), 1, 8); err == nil {
    t.Fail()
  }
}

func TestOneByte(t *testing.T) {
  // 10[10 1]010
  val, err := PickBits([]byte{0xaa}, 2, 3)
  if err != nil {
    t.Fail()
  }
  // 101 = 5
  if val != 5 {
    t.Fail()
  }
}

func TestTwoByte(t *testing.T) {
  // 1111 111[1 00]00 0000
  val, err := PickBits([]byte{0xff, 0x00}, 7, 3)
  if err != nil {
    t.Fail()
  }
  // [100] = 4
  if val != 4 {
    t.Fail()
  }
  // [0001 1011 0]111 1111
  val, err = PickBits([]byte{0x1b, 0x7f}, 0, 9)
  if err != nil {
    t.Fail()
  }
  // [000110110] = 54
  if val != 54 {
    t.Fail()
  }
}

func TestFourByte(t *testing.T) {
  // [1010 1010 1010 1010]
  val, err := PickBits([]byte{0xaa, 0xaa, 0xaa, 0xaa}, 0, 32)
  if err != nil {
    t.Fail()
  }
  // [1010101010101010] = 2863311530
  if val != 2863311530 {
    t.Fail()
  }
}

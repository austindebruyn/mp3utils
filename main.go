package main

import (
  "fmt"
  "os"
  "github.com/austindebruyn/mp3utils/mp3"
)

func main() {
  if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
    fmt.Printf("Usage: %s <filename>\n", os.Args[0])
    os.Exit(1)
  }

  filename := os.Args[1]

  file, err := os.Open(filename)
  if err != nil {
    fmt.Printf("Couldn't open %s.\n", filename)
    os.Exit(1)
  }
  defer file.Close()

  for i := 0; i < 20; i++ {
    frame, err := mp3.ReadFrame(file)
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
    fmt.Println(frame)
  }
}

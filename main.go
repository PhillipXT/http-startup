package main

import (
    "fmt"
    "io"
    "log"
    "os"
)

func main() {
    file, err := os.Open("messages.txt")
    if err != nil {
        log.Fatal(err)
    }

    for {
        buffer := make([]byte, 8)
        n, err := file.Read(buffer)
        if err == io.EOF {
            break
        }
        if n > 0 {
            fmt.Printf("read: %s\n", string(buffer))
        }
    }
}

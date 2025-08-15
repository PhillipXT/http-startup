package main

import (
    "fmt"
    "io"
    "log"
    "os"
)

const filePath = "messages.txt"

func main() {

    file, err := os.Open(filePath)
    if err != nil {
        log.Fatalf("could not open %s: %s\n", filePath, err)
    }

    defer file.Close()

    fmt.Printf("Reading data from %s:\n", filePath)

    for {
        buffer := make([]byte, 8)
        n, err := file.Read(buffer)
        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Printf("error: %s\n", err.Error())
            break
        }
        fmt.Printf("\tread: %s\n", string(buffer[:n]))
    }
}

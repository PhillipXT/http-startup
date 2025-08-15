package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "strings"
)

const filePath = "messages.txt"

func main() {

    file, err := os.Open(filePath)
    if err != nil {
        log.Fatalf("could not open %s: %s\n", filePath, err)
    }

    defer file.Close()

    fmt.Printf("Reading data from %s:\n", filePath)

    line := ""

    for {
        buffer := make([]byte, 8)
        n, err := file.Read(buffer)
        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Printf("error: %s\n", err.Error())
            break
        }
        parts := strings.Split(string(buffer[:n]), "\n")
        for i := 0; i < len(parts) - 1; i++ {
            fmt.Printf("read: %s%s\n", line, parts[i])
            line = ""
        }
        line += parts[len(parts)-1]
    }
}

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

    lines := getLinesChannel(file)
    for line := range lines {
        fmt.Println("read:", line)
    }

}

func getLinesChannel(f io.ReadCloser) <-chan string {

    ch := make(chan string)

    go func() {
        defer close(ch)
        line := ""
        for {
            buffer := make([]byte, 8)
            n, err := f.Read(buffer)
            if err == io.EOF {
                break
            } else if err != nil {
                fmt.Printf("error: %s\n", err.Error())
                break
            }
            parts := strings.Split(string(buffer[:n]), "\n")
            for i := 0; i < len(parts) - 1; i++ {
                ch <- fmt.Sprintf("%s%s", line, parts[i])
                line = ""
            }
            line += parts[len(parts)-1]
        }
    }()

    return ch
}

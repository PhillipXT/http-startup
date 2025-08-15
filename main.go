package main

import (
    "fmt"
    "io"
    "log"
    "net"
    "strings"
)

const filePath = "messages.txt"

func main() {

    listener, err := net.Listen("tcp", "localhost:42069")
    if err != nil {
        log.Fatal(err)
    }

    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println("Connection accepted:", conn)

        lines := getLinesChannel(conn)

        for line := range lines {
            fmt.Println("read:", line)
        }

        fmt.Println("Connection closed")
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

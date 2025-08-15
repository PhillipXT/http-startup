package main

import (
    "fmt"
    "io"
    "log"
    "net"
    "strings"
)

const filePath = "messages.txt"
const port = ":42069"

// Run the application, then print to the port from the command line:
// printf "Can you hear me now?\r\n" | nc -w 1 127.0.0.1 42069

func main() {

    listener, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("error listening for TCP traffic: %s\n", err)
    }

    defer listener.Close()

    fmt.Println("Listening for TCP traffic on", port)

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal("error: %s\n", err)
        }

        fmt.Println("Accepted connection from:", conn.RemoteAddr())

        lines := getLinesChannel(conn)

        for line := range lines {
            fmt.Println(line)
        }

        fmt.Printf("Connection to %s closed\n", conn.RemoteAddr())
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

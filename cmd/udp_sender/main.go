package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "net"
)

const (
    port = ":42069"
)

func main() {
    udp, err := net.ResolveUDPAddr("udp", port)
    if err != nil {
        log.Fatalf("Error resolving UDP address: %s\n", err)
    }

    conn, err := net.DialUDP("udp", nil, udp)
    if err != nil {
        log.Fatalf("Unable to dial into udp network: %s\n", err)
    }

    defer conn.Close()

    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Printf("> ")
        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading input: ", err)
        }
        _, err = conn.Write([]byte(input))
        if err != nil {
            fmt.Println("Error writing input: ", err)
        }
    }
}

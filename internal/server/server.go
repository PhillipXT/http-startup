package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/PhillipXT/http-startup/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %s", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	h := response.GetDefaultHeaders(0)

	err := response.WriteStatusLine(conn, response.StatusCodeOK)
	if err != nil {
		log.Printf("Error writing status line: %s", err)
	}

	err = response.WriteHeaders(conn, h)
	if err != nil {
		log.Printf("Error writing headers: %s", err)
	}

	//res := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n"
	//conn.Write([]byte(res))
}

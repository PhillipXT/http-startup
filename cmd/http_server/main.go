package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PhillipXT/http-startup/internal/request"
	"github.com/PhillipXT/http-startup/internal/response"
	"github.com/PhillipXT/http-startup/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		body := "<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>"
		h := response.GetDefaultHeaders(len(body))
		h.Set("Content-Type", "text/html")
		w.WriteStatusLine(response.StatusCodeBadRequest)
		w.WriteHeaders(h)
		w.WriteBody([]byte(body))
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		body := "<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>"
		h := response.GetDefaultHeaders(len(body))
		h.Set("Content-Type", "text/html")
		w.WriteStatusLine(response.StatusCodeInternalServerError)
		w.WriteHeaders(h)
		w.WriteBody([]byte(body))
		return
	}
	body := "<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>"
	h := response.GetDefaultHeaders(len(body))
	h.Set("Content-Type", "text/html")
	w.WriteStatusLine(response.StatusCodeOK)
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))
	return
}

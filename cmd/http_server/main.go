package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	target := req.RequestLine.RequestTarget

	if strings.HasPrefix(target, "/httpbin") {
		proxyHandler(w, req)
		return
	}

	if target == "/yourproblem" {
		handler400(w, req)
		return
	}

	if target == "/myproblem" {
		handler500(w, req)
		return
	}

	if target == "/video" {
		videoHandler(w, req)
		return
	}

	handler200(w, req)
}

func videoHandler(w *response.Writer, req *request.Request) {

	log.Println("Processing video request")

	bytes, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		handler500(w, req)
		return
	}

	h := response.GetDefaultHeaders(len(bytes))
	h.Set("Content-Type", "video/mp4")

	w.WriteStatusLine(response.StatusCodeOK)
	w.WriteHeaders(h)
	w.WriteBody(bytes)
}

func proxyHandler(w *response.Writer, req *request.Request) {

	log.Println("Processing chunked response")

	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := "https://httpbin.org/" + target

	h := response.GetDefaultHeaders(0)
	h.Delete("Content-Length")
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Trailer", "X-Content-Sha256, X-Content-Length")

	r, err := http.Get(url)
	if err != nil {
		handler500(w, req)
		return
	}

	w.WriteStatusLine(response.StatusCodeOK)
	w.WriteHeaders(h)

	buffer := make([]byte, 64)
	body := []byte{}

	for {
		n, err := r.Body.Read(buffer)
		if n > 0 {
			body = append(body, buffer[:n]...)
			_, err := w.WriteChunkedBody(buffer[:n])
			if err != nil {
				fmt.Println("Error writing chunked body:", err)
				break
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
	}
	log.Println("Chunked response complete")

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Println("Error completing chunked body:", err)
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(body))

	t := map[string]string{}
	t["X-Content-Sha256"] = hash
	t["X-Content-Length"] = fmt.Sprintf("%d", len(body))

	err = w.WriteTrailers(t)
	if err != nil {
		fmt.Println("Error writing trailers:", err)
	}
}

func handler200(w *response.Writer, req *request.Request) {
	body := "<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>"
	h := response.GetDefaultHeaders(len(body))
	h.Set("Content-Type", "text/html")
	w.WriteStatusLine(response.StatusCodeOK)
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))
}

func handler400(w *response.Writer, req *request.Request) {
	body := "<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>"
	h := response.GetDefaultHeaders(len(body))
	h.Set("Content-Type", "text/html")
	w.WriteStatusLine(response.StatusCodeBadRequest)
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))
}

func handler500(w *response.Writer, req *request.Request) {
	body := "<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>"
	h := response.GetDefaultHeaders(len(body))
	h.Set("Content-Type", "text/html")
	w.WriteStatusLine(response.StatusCodeInternalServerError)
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))
}

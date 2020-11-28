package main

import (
	"crypto/tls"
	"fmt"
	"time"
)

/**
functional option
*/
type Server struct {
	Addr string
	Port int
	Protocol string
	Timeout time.Duration
	MaxCounts int
	TLS *tls.Config
}

type Option func(*Server)

func Protocol(p string) Option {
	return func(server *Server) {
		server.Protocol = p
	}
}

func Timeout(timeout time.Duration) Option {
	return func(server *Server) {
		server.Timeout = timeout
	}
}

func MaxCounts(maxCounts int) Option {
	return func(server *Server) {
		server.MaxCounts = maxCounts
	}
}

func TLS(tls *tls.Config) Option {
	return func(server *Server) {
		server.TLS = tls
	}
}

func NewServer(addr string, port int, options ...func(*Server)) (*Server, error) {
	// you can setup a default config
	srv := Server{
		Addr:addr,
		Port:port,
	}
	for _, option := range options {
		option(&srv)
	}

	return &srv, nil
}

func main() {
	s1, _ := NewServer("localhost", 1111)
	s2, _ := NewServer("localhost", 2222, Protocol("tcp"))
	fmt.Println(s1)
	fmt.Println(s2)
}

package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/sambakker4/httpfromtcp/internal/request"
	"github.com/sambakker4/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Server struct {
	Listener    net.Listener
	isClosed    *atomic.Bool
	HandlerFunc Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return &Server{}, err
	}

	isClosed := &atomic.Bool{}
	isClosed.Store(false)

	server := &Server{
		Listener:    listener,
		isClosed:    isClosed,
		HandlerFunc: handler,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	err := s.Listener.Close()
	return err
}

func (s *Server) listen() {
	for !s.isClosed.Load() {
		connection, err := s.Listener.Accept()
		if s.isClosed.Load() {
			break
		}

		if err != nil {
			log.Printf("connection error: %s\n", err.Error())
			continue
		}
		s.handle(connection)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("request error: %s", err.Error())
		return
	}

	writer := response.Writer{
		Writer: conn,
	}

	s.HandlerFunc(&writer, req)
}

package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/sambakker4/httpfromtcp/internal/request"
	"github.com/sambakker4/httpfromtcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

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
		Listener: listener,
		isClosed: isClosed,
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

func WriteHandlerError(w io.Writer, handlerError *HandlerError) {
	err := response.WriteStatusLine(w, handlerError.StatusCode)
	if err != nil {
		log.Printf("error: %s\n", err.Error())
	}

	headers := response.GetDefaultHeaders(len(handlerError.Message))
	err = response.WriteHeaders(w, headers)
	if err != nil {
		log.Printf("error: %s\n", err.Error())
	}

	_, err = w.Write([]byte(handlerError.Message))
	if err != nil {
		log.Printf("error: %s\n", err.Error())
	}

}

func (s *Server) handle(conn net.Conn) {
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("request error: %s", err.Error())
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	handlerErr := s.HandlerFunc(buf, req)

	if handlerErr != nil {
		WriteHandlerError(conn, handlerErr)
	} else {
		headers := response.GetDefaultHeaders(len(buf.Bytes()))
		err = response.WriteStatusLine(conn, response.Success)
		if err != nil {
			log.Printf("error writing status line: %s\n", err.Error())
		}

		err = response.WriteHeaders(conn, headers)
		if err != nil {
			log.Printf("error writing headers: %s\n", err.Error())
		}

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Printf("error writing body: %s\n", err.Error())
		}
	}

	err = conn.Close()
	if err != nil {
		log.Printf("closing connection error: %s\n", err.Error())
	}
}

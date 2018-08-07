package http

import (
	"net"
	"net/http"

	"github.com/Tinee/prog-image"
)

// Server represents an HTTP server.
type Server struct {
	ln      net.Listener
	handler *Handler

	Addr string
}

// NewServer returns a new instance of Server.
func NewServer(addr string, h *Handler) (*Server, error) {
	if h == nil {
		return nil, progimage.ErrInvalidArgument
	}
	return &Server{
		Addr:    addr,
		handler: h,
	}, nil
}

// Open opens a socket and serves the HTTP server.
func (s *Server) Open() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.ln = ln

	go func() { http.Serve(s.ln, s.handler) }()

	return nil
}

// Close the underlaying socket.
func (s *Server) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}
	return nil
}

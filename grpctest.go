package grpctest

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type Server struct {
	ctx      context.Context
	server   *grpc.Server
	listener *bufconn.Listener
}

func NewServer() *Server {
	return &Server{
		ctx:      context.Background(),
		server:   grpc.NewServer(),
		listener: bufconn.Listen(bufSize),
	}
}

func (s *Server) WithContext(ctx context.Context) *Server {
	s.ctx = ctx
	return s
}

func (s *Server) RunServer(t *testing.T, fn func(*grpc.Server)) {
	t.Helper()
	fn(s.server)
	go func() {
		t.Helper()
		if err := s.server.Serve(s.listener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func (s *Server) ClientConn(t *testing.T) *grpc.ClientConn {
	t.Helper()
	conn, err := grpc.DialContext(s.ctx, "bufnet", grpc.WithContextDialer(s.dialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("clientConn: %v", err)
	}
	return conn
}

func (s *Server) Close() {
	s.server.Stop()
	s.listener.Close()
}

func (s *Server) dialer(_ context.Context, _ string) (net.Conn, error) {
	return s.listener.Dial()
}

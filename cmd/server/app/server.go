package app

import (
	"fmt"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/ruijzhan/dungeons/cmd/server/http"
	"github.com/ruijzhan/dungeons/pkg/resolve"
	"google.golang.org/grpc"
)

type ServerOption struct {
	//http server options
	HttpListenAddr string

	//grpc server options
	GrpcListenAddr string
}

type Server struct {
	opt ServerOption

	//gin for http server
	*gin.Engine

	dns resolve.Resolver

	//grpc server
	rpc *grpc.Server
}

// NewServer creates a new HTTP server with custom options.
func NewServer(opt ServerOption) *Server {
	srv := &Server{
		opt:    opt,
		Engine: gin.New(),
		dns:    resolve.New(),
		rpc:    grpc.NewServer(),
	}
	srv.addRoutes()
	return srv
}

// Run starts the server and listens for incoming HTTP and gRPC requests.
func (s *Server) Run() (err error) {

	errc := make(chan error, 1)

	go func() {
		if err := s.rpc.Serve(listen(s.opt.GrpcListenAddr)); err != nil {
			errc <- fmt.Errorf("GRPC server stopped unexpectedly: %v", err)
		}
	}()

	go func() {
		if err := s.Engine.RunListener(listen(s.opt.HttpListenAddr)); err != nil {
			errc <- fmt.Errorf("http server stopped unexpectedly: %v", err)
		}
	}()

	return <-errc
}

/*
listen listens on the specified TCP network address addr and returns a net.Listener object.
Parameters:

- addr (string): the address to listen on.
Returns:

- ln (net.Listener): the net.Listener object that is listening on the specified address.
*/
func listen(addr string) (ln net.Listener) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on tcp port %v", err)
	}
	return
}

// addRoutes adds HTTP handlers to the server's router.
func (s *Server) addRoutes() {
	http.AddPing(s.Engine)
	http.AddCheckHost(s.Engine, s.dns)
}

// Package server provides all the implementation for the RPC handlers
package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/auburnhacks/sponsor/pkg/utils"
	api "github.com/auburnhacks/sponsor/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	// RPCAddr is a flag that specifies the addres at which the gRPC server will run
	RPCAddr *string
	// GatewayAddr is a flag that specifies the address at which the gateway will run
	GatewayAddr *string
)

const (
	tokenIssuer = "sponsor_auburnhacks"
)

// Server is a server that defines an interface for this package
type Server interface {
	// Serve serves the server
	Serve() error
	// Stop will attempt a graceful shutdown
	Stop() error
}

// New takes a private key for encryption and returns a struct
// that implements the Server interface
func New(privKey []byte) Server {
	log.Debugf("gateway address: %s", *GatewayAddr)
	log.Debugf("rpc server addr: %s", *RPCAddr)
	return &rpcServer{
		privKey: privKey,
	}
}

// rpcServer is the grpc server struct used for all handler calls
type rpcServer struct {
	privKey []byte
	srvErr  chan error
	grpcSrv *grpc.Server
	gwSrv   *http.Server
}

func (s *rpcServer) Serve() error {
	l, err := net.Listen("tcp", *RPCAddr)
	if err != nil {
		return err
	}
	go s.serveGRPC(l)
	go s.serveGateway(*GatewayAddr, *RPCAddr)
	return <-s.srvErr
}

func (s *rpcServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := s.gwSrv.Shutdown(ctx); err != nil {
		return err
	}
	s.grpcSrv.GracefulStop()
	return nil
}

func (s *rpcServer) serveGRPC(l net.Listener) {
	s.grpcSrv = grpc.NewServer(grpc.UnaryInterceptor(utils.UnaryAuthInterceptor))
	api.RegisterSponsorServiceServer(s.grpcSrv, s)
	if err := s.grpcSrv.Serve(l); err != nil {
		log.Errorf("error while serving grpc server: %v", err)
		s.srvErr <- err
	}
}

func (s *rpcServer) serveGateway(listenAddr, serviceEndpoint string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard,
			&runtime.JSONPb{OrigName: false}),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := api.RegisterSponsorServiceHandlerFromEndpoint(ctx, mux, serviceEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	s.gwSrv = &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}
	if err := s.gwSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Errorf("error while serving gateway: %v", err)
		s.srvErr <- err
	}
}

// Package server provides all the implementation for the RPC handlers
package server

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/auburnhacks/sponsor/pkg/utils"
	api "github.com/auburnhacks/sponsor/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	signingKey = []byte("supersecret")
)

// SponsorServer is a struct that implements the SponsorServiceServer interface
// that is auto-gernerated by gRPC
type SponsorServer struct {
	DB   *mongo.Client
	quit chan struct{}
	tWg  sync.WaitGroup
}

// ListenAndServe is a helper func that invokes the and serves the gRPC server
// and the gateway in separate goroutines
func ListenAndServe(srv *SponsorServer, l net.Listener, listenAddr, serviceEndpoint *string) {
	go srv.serveGRPC(l)
	go srv.serveGateway(listenAddr, serviceEndpoint)
}

// NewSponsorServer is a constructer that returns an instance of the SponsorServer
func NewSponsorServer() *SponsorServer {
	return &SponsorServer{
		quit: make(chan struct{}, 2),
	}
}

// Shutdown is a function that wait on the quit channel and signals the gateway
// and gRPC server to shutdown
func (s *SponsorServer) Shutdown() {
	s.quit <- struct{}{}
	s.tWg.Wait()
	s.quit <- struct{}{}
	s.tWg.Wait()
	close(s.quit)
}

func (s *SponsorServer) serveGRPC(l net.Listener) {
	srv := grpc.NewServer(grpc.UnaryInterceptor(utils.UnaryAuthInterceptor))
	api.RegisterSponsorServiceServer(srv, s)

	go func() {
		log.Println("serving grpc")
		if err := srv.Serve(l); err != nil {
			log.Fatalf("error serving: %v", err)
		}
	}()
	<-s.quit
	log.Infof("terminating rpc server")
	s.tWg.Add(1)
	srv.GracefulStop()
	s.tWg.Done()
}

func (s *SponsorServer) serveGateway(listenAddr, serviceEndpoint *string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard,
		&runtime.JSONPb{OrigName: false}))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := api.RegisterSponsorServiceHandlerFromEndpoint(ctx, mux, *serviceEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{
		Addr:    *listenAddr,
		Handler: mux,
	}
	go func() {
		log.Info("serving gateway")
		log.Fatal(srv.ListenAndServe())
	}()
	<-s.quit
	log.Println("terminating gateway")
	s.tWg.Add(1)
	srv.Shutdown(ctx)
	s.tWg.Done()
}

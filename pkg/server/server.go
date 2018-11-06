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

const (
	tokenIssuer = "sponsor_auburnhacks"
)

// SponsorServer is a struct that implements the SponsorServiceServer interface
// that is auto-gernerated by gRPC
type SponsorServer struct {
	DB     *mongo.Client
	quit   chan struct{}
	tWg    sync.WaitGroup
	jwtKey []byte
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

// WithKey is a function that modifies the SponsorServer and return the
// instance
func (ss *SponsorServer) WithKey(key []byte) *SponsorServer {
	ss.jwtKey = key
	return ss
}

// Shutdown is a function that wait on the quit channel and signals the gateway
// and gRPC server to shutdown
func (ss *SponsorServer) Shutdown() {
	ss.quit <- struct{}{}
	ss.tWg.Wait()
	ss.quit <- struct{}{}
	ss.tWg.Wait()
	close(ss.quit)
}

func (ss *SponsorServer) serveGRPC(l net.Listener) {
	srv := grpc.NewServer(grpc.UnaryInterceptor(utils.UnaryAuthInterceptor))
	api.RegisterSponsorServiceServer(srv, ss)

	go func() {
		log.Println("serving grpc")
		if err := srv.Serve(l); err != nil {
			log.Fatalf("error serving: %v", err)
		}
	}()
	<-ss.quit
	log.Infof("terminating rpc server")
	ss.tWg.Add(1)
	srv.GracefulStop()
	ss.tWg.Done()
}

func (ss *SponsorServer) serveGateway(listenAddr, serviceEndpoint *string) {
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
	<-ss.quit
	log.Println("terminating gateway")
	ss.tWg.Add(1)
	srv.Shutdown(ctx)
	ss.tWg.Done()
}

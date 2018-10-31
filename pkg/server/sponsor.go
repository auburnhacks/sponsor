package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/auburnhacks/sponsor/pkg/utils"
	api "github.com/auburnhacks/sponsor/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mongodb/mongo-go-driver/mongo"
	"google.golang.org/grpc"
)

func ListenAndServe(srv *SponsorServer, l net.Listener, listenAddr, serviceEndpoint *string) {
	go srv.serveGRPC(l)
	go srv.serveGateway(listenAddr, serviceEndpoint)
}

func NewSponsorServer() *SponsorServer {
	return &SponsorServer{
		quit: make(chan struct{}, 2),
	}
}

type SponsorServer struct {
	DB   *mongo.Client
	quit chan struct{}
	tWg  sync.WaitGroup
}

func (s *SponsorServer) CreateAdmin(ctx context.Context, req *api.CreateAdminRequest) (*api.CreateAdminResponse, error) {
	return &api.CreateAdminResponse{}, nil
}

func (s *SponsorServer) GetAdmin(ctx context.Context, req *api.GetAdminRequest) (*api.GetAdminResponse, error) {
	return &api.GetAdminResponse{}, nil
}

func (s *SponsorServer) DeleteAdmin(ctx context.Context, req *api.DeleteAdminRequest) (*api.DeleteAdminResponse, error) {
	return &api.DeleteAdminResponse{}, nil
}

func (s *SponsorServer) LoginAdmin(ctx context.Context, req *api.LoginAdminRequest) (*api.LoginAdminResponse, error) {
	return &api.LoginAdminResponse{}, nil
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
	log.Println("terminating rpc server")
	s.tWg.Add(1)
	srv.GracefulStop()
	s.tWg.Done()
}

func (s *SponsorServer) serveGateway(listenAddr, serviceEndpoint *string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard,
		&runtime.JSONPb{OrigName: true}))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := api.RegisterSponsorServiceHandlerFromEndpoint(ctx, mux, *serviceEndpoint, opts); err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{
		Addr:    *listenAddr,
		Handler: mux,
	}
	go func() {
		log.Println("serving gateway")
		log.Fatal(srv.ListenAndServe())
	}()
	<-s.quit
	log.Println("terminating gateway")
	s.tWg.Add(1)
	srv.Shutdown(ctx)
	s.tWg.Done()
}

func (s *SponsorServer) Shutdown() {
	s.quit <- struct{}{}
	s.tWg.Wait()
	s.quit <- struct{}{}
	s.tWg.Wait()
	close(s.quit)
}

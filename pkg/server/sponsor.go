package server

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/auburnhacks/sponsor/pkg/admin"
	"github.com/auburnhacks/sponsor/pkg/utils"
	api "github.com/auburnhacks/sponsor/proto"
	"github.com/dgrijalva/jwt-go"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	signingKey = []byte("supersecret")
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
	admin := admin.New(req.Name, req.Email, req.PasswordPlainText)
	if err := admin.Register(); err != nil {
		log.Errorf("error while registering admin: %v", err)
		return nil, err
	}
	log.Debugf("%+v", admin)
	return &api.CreateAdminResponse{
		Admin: &api.Admin{
			Email:   admin.Email,
			AdminID: int64(admin.ID),
			ACL:     admin.ACL,
		},
	}, nil
}

func (s *SponsorServer) GetAdmin(ctx context.Context, req *api.GetAdminRequest) (*api.GetAdminResponse, error) {
	return &api.GetAdminResponse{}, nil
}

func (s *SponsorServer) DeleteAdmin(ctx context.Context, req *api.DeleteAdminRequest) (*api.DeleteAdminResponse, error) {
	return &api.DeleteAdminResponse{}, nil
}

func (s *SponsorServer) LoginAdmin(ctx context.Context, req *api.LoginAdminRequest) (*api.LoginAdminResponse, error) {
	// TODO: look up database for the user
	admin, err := admin.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	log.Debugf("%+v", admin)
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Date(2019, 9, 30, 0, 0, 0, 0, time.UTC).Unix(),
		Issuer:    "sponsor_auburnhacks",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(signingKey)
	if err != nil {
		return nil, err
	}
	return &api.LoginAdminResponse{
		Token: tokenStr,
		Admin: &api.Admin{
			AdminID: int64(admin.ID),
			Name:    admin.Name,
			Email:   admin.Email,
			ACL:     admin.ACL,
		},
	}, nil
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

func (s *SponsorServer) Shutdown() {
	s.quit <- struct{}{}
	s.tWg.Wait()
	s.quit <- struct{}{}
	s.tWg.Wait()
	close(s.quit)
}

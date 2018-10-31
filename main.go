package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/auburnhacks/sponsor/pkg/server"
	"github.com/mongodb/mongo-go-driver/mongo"
)

var (
	sponsorService *string
	dbURI          *string
	listenAddr     *string
)

func init() {
	sponsorService = flag.String("sponsor_endpoint", "localhost:10000", "hostport for sponsor service")
	dbURI = flag.String("db_uri", "mongodb://localhost:27017", "mongoDB connection URI")
	listenAddr = flag.String("listen_addr", "localhost:8080", "listen_addr for grpc gateway")

	flag.Parse()
}

func main() {
	// Connect to the database
	client, err := mongo.NewClient(*dbURI)
	if err != nil {
		log.Fatalf("error creating mongo client: %v\n", err)
	}
	if err := client.Connect(context.TODO()); err != nil {
		log.Fatalf("error connecting to mongo database: %v\n", err)
	}

	srv := server.NewSponsorServer()
	srv.DB = client

	// gRPC server listener
	l, err := net.Listen("tcp", *sponsorService)
	if err != nil {
		log.Fatalf("error create listener: %v", err)
	}
	go func() {
		log.Printf("server running on pid: %d\n", os.Getpid())
		server.ListenAndServe(srv, l, listenAddr, sponsorService)
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	signal := <-quit
	log.Printf("received %v signal, terminating server", signal)
	srv.Shutdown()
	os.Exit(0)
}

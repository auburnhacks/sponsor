package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/auburnhacks/sponsor/pkg/db"
	"github.com/auburnhacks/sponsor/pkg/server"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var (
	sponsorService *string
	listenAddr     *string
	dbHost         *string
	dbPort         *int
	dbUser         *string
	dbPassword     *string
	dbName         *string
	dbSSLMode      *bool
	debug          *bool
)

func init() {
	/*
		DEFAULT CONNECTION STRING (DEVELOPMENT)
		===================================================================================================
		host=localhost port=5432 user=dev_user password=test123 dbname=auburnhacks_sponsors sslmode=disable
		===================================================================================================
		* If you are using docker run the following command to start a lock postgres database running at
		localhost:5432
		docker run \--name pg-sponsor-local \
		-p 5432:5432 \
		-e POSTGRES_PASSWORD=test123 \
		-e POSTGRES_USER=dev_user \
		-e POSTGRES_DB=auburnhacks_sponsors \
		-v ~/data/pg-data:/var/lib/postgresql/data \
		-d \
		postgres
	*/
	// database flags
	dbHost = flag.String("db_host", "localhost", "hostname for postgres database")
	dbPort = flag.Int("db_port", 5432, "postgres database port number")
	dbUser = flag.String("db_user", "dev_user", "username for the postgres database.[make sure this user had all privilages]")
	dbPassword = flag.String("db_password", "test123", "postgres database password. Make sure its secure")
	dbName = flag.String("db_name", "auburnhacks_sponsors", "database name that has to be accessed after connecting.")
	dbSSLMode = flag.Bool("db_ssl_mode", false, "use this flag to connecto to the database using ssl")

	// gRPC and gateway configuration flags
	sponsorService = flag.String("sponsor_endpoint", "localhost:10000", "hostport for sponsor service")
	listenAddr = flag.String("listen_addr", "localhost:8080", "listen_addr for grpc gateway")

	// application related flags
	debug = flag.Bool("debug", false, "use this flag to set debug logging [DONT USE IN PRODUCTION]")

	flag.Parse()
}

func main() {
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
	// building the DB connection string
	sslMode := "disable" // change to enable in production
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=%s", *dbHost, *dbPort, *dbUser, *dbPassword,
		*dbName, sslMode)

	log.Debugf("connecting to database with info: %s", psqlInfo)

	dbConn, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	// Initialize the global database connection with the opened connection
	db.Conn = dbConn

	quit := make(chan os.Signal)

	// Run all database migrations
	log.Info("running database migrations...")
	if err := db.MigrateDB(quit); err != nil {
		log.Fatal(err)
	}

	srv := server.NewSponsorServer()

	// gRPC server listener
	l, err := net.Listen("tcp", *sponsorService)
	if err != nil {
		log.Fatalf("error create listener: %v", err)
	}
	go func() {
		log.Infof("server running on pid: %d", os.Getpid())
		server.ListenAndServe(srv, l, listenAddr, sponsorService)
	}()
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	signal := <-quit
	log.Infof("received %v signal, terminating server", signal)
	db.Conn.Close()
	srv.Shutdown()
	os.Exit(0)
}

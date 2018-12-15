package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/auburnhacks/sponsor/pkg/auth"
	"github.com/auburnhacks/sponsor/pkg/db"
	"github.com/auburnhacks/sponsor/pkg/participant"
	"github.com/auburnhacks/sponsor/pkg/server"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var (
	dbHost          *string
	dbPort          *int
	dbUser          *string
	dbPassword      *string
	dbName          *string
	dbSSLMode       *bool
	debug           *bool
	quillMongoURI   *string
	resumesMongoURI *string
	syncDuration    *time.Duration
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

	// application related flags
	debug = flag.Bool("debug", false, "use this flag to set debug logging [DONT USE IN PRODUCTION]")
	quillMongoURI = flag.String("quill_db_uri", "mongodb://localhost:27017/quill", "database URI for quill")
	resumesMongoURI = flag.String("resumes_db_uri", "mongodb://localhost:27017/resumes", "database uri for resumes")
	syncDuration = flag.Duration("sync_duration", 1*time.Minute, "sleep duration every sync period")

	server.GatewayAddr = flag.String("gateway_addr", "localhost:8080", "grpc gateway listen addr")
	server.RPCAddr = flag.String("rpc_addr", "localhost:10000", "grpc server listening addr")

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
	log.Info("running database migrations")
	if err := db.Migrate(quit); err != nil {
		log.Fatalf("error migrating database: %v", err)
	}
	// Read the signing key for JWT tokens
	key, err := auth.LoadJWTKey(filepath.Join(".", "jwt_key_dev"))
	if err != nil {
		log.Fatalf("error loading jwt key: %v", err)
	}
	s := server.New(key)
	go func() {
		log.Infof("starting server pid: %d", os.Getpid())
		if err := s.Serve(); err != nil {
			log.Fatalf("error while serving: %v", err)
		}
	}()
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	pollerQuit := make(chan struct{})
	// Start syncing the participants from the external database
	err = participant.Sync(*quillMongoURI, *resumesMongoURI, *syncDuration, pollerQuit)
	if err != nil {
		log.Errorf("error while syncing from external db: %v", err)
		pollerQuit <- struct{}{}
		os.Exit(1)
	}
	// blocks forever until SIGTERM signal is received for container
	signal := <-quit
	pollerQuit <- struct{}{}
	log.Infof("received %v signal, terminating server", signal)
	// Dropping migrations
	if *debug {
		if err := db.Drop(); err != nil {
			log.Fatalf("error dropping database: %v", err)
		}
	}
	db.Conn.Close()
	if err := s.Stop(); err != nil {
		log.Infof("error stopping server: %v", err)
	}
	os.Exit(0)
}

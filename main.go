package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/axamon/hextest/database/psql"
	redisdb "github.com/axamon/hextest/database/redis"
	"github.com/axamon/hextest/ticket"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var version = "0.1.0"

func main() {

	// dbType is a flag used to choose which backend database to use.
	dbType := flag.String("database", "redis", "database type [redis, psql]")
	redisAddress := flag.String("redis", "localhost:6379", "Address of redis server")
	port := flag.Int("port", 3000, "tcp port to use")

	// parses the flag.
	flag.Parse()

	// ticketRepo idebtifies which repository to use.
	var ticketRepo ticket.Repository

	// choses which db to use as repository.
	switch *dbType {
	case "psql":
		pconn := postgresConnection("postgresql://postgres@localhost/ticket?sslmode=disable")
		defer pconn.Close()
		ticketRepo = psql.NewPostgresTicketRepository(pconn)
	case "redis":
		rconn := redisConnection(*redisAddress)
		defer rconn.Close()
		ticketRepo = redisdb.NewRedisTicketRepository(rconn)
	default:
		panic("Unknown database")
	}

	ticketService := ticket.NewService(ticketRepo)
	ticketHandler := ticket.NewTicketHandler(ticketService)

	/* HTTP ROUTES */
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tickets", ticketHandler.GetAll).Methods("GET")
	router.HandleFunc("/tickets/{id}", ticketHandler.GetByID).Methods("GET")
	router.HandleFunc("/tickets/delete/{id}", ticketHandler.DeleteByID).Methods("GET")
	router.HandleFunc("/tickets", ticketHandler.Create).Methods("POST")

	// main handle router
	http.Handle("/", accessControl(router))
	/* HTTP ROUTES END */
	id := uuid.New().String()
	// register microservice on Consul.
	registerService(id, "ticket", version, "127.0.0.1", *port, "5m", "30s", "2s")

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :", *port)
		errs <- http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("terminated %s", <-errs)

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

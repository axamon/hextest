package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/axamon/hextest/database/psql"
	redisdb "github.com/axamon/hextest/database/redis"
	"github.com/axamon/hextest/ticket"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	// dbType is a flag used to choose which backend database to use.
	dbType := flag.String("database", "redis", "database type [redis, psql]")
	redisAddress := flag.String("redis", "localhost:6379", "Address of redis server")

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

	// register microservice on Consul.
	registerService()

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :3000")
		errs <- http.ListenAndServe(":3000", nil)
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

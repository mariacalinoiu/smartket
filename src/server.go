package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/mariacalinoiu/smartket/src/datasources"
	"github.com/mariacalinoiu/smartket/src/handlers"
)

type server struct {
	mux    *http.ServeMux
	logger *log.Logger
}

type option func(*server)

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.log("Method: %s, Path: %s", r.Method, r.URL.Path)
	s.mux.ServeHTTP(w, r)
}

func (s *server) log(format string, v ...interface{}) {
	s.logger.Printf(format+"\n", v...)
}

func logWith(logger *log.Logger) option {
	return func(s *server) {
		s.logger = logger
	}
}

func setup(logger *log.Logger, db datasources.DBClient) *http.Server {
	server := newServer(db, logWith(logger))
	return &http.Server{
		Addr:         ":8081",
		Handler:      server,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  600 * time.Second,
	}
}

func newServer(db datasources.DBClient, options ...option) *server {
	s := &server{logger: log.New(ioutil.Discard, "", 0)}

	for _, o := range options {
		o(s)
	}

	s.mux = http.NewServeMux()

	s.mux.HandleFunc("/departments",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleDepartments(w, r, db, s.logger)
		},
	)
	s.mux.HandleFunc("/categories",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleCategories(w, r, db, s.logger)
		},
	)
	s.mux.HandleFunc("/products",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleProducts(w, r, db, s.logger)
		},
	)
	s.mux.HandleFunc("/orders",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleOrdersAdd(w, r, db, s.logger)
		},
	)
	s.mux.HandleFunc("/orders/delete",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleOrdersDelete(w, r, db, s.logger)
		},
	)
	s.mux.HandleFunc("/orders/update",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleOrdersUpdate(w, r, db, s.logger)
		},
	)

	return s
}

func main() {
	logger := log.New(os.Stdout, "", 0)
	db := datasources.GetClient("user", "password", "onlinestore")
	hs := setup(logger, db)

	logger.Printf("Listening on http://localhost%s\n", hs.Addr)
	go func() {
		if err := hs.ListenAndServe(); err != nil {
			logger.Println(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	logger.Println("Shutting down webserver.")
	os.Exit(0)
}

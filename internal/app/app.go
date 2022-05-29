package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"l0/internal/nats"
	store "l0/internal/repository"

	"github.com/nats-io/stan.go"
)

type App struct {
	Sc         *stan.Conn
	HTTPClient *http.Client
	DB         *sql.DB
	Rep        *store.Repository
}

func CreateApp() (*App, error) {
	sc, err := nats.NewNats()

	if err != nil {
		return nil, err
	}

	return &App{
		Sc: sc,
	}, nil
}

func (a *App) Run() error {
	a.RunNats()

	err := a.RunDB()
	if err != nil {
		return fmt.Errorf("error on creating DB connection: %v", err)
	}

	err = a.CreateRepository()
	if err != nil {
		return fmt.Errorf("error on creating repository: %v", err)
	}

	go a.RunHTTP()

	fmt.Printf("HTTP started!")

	return nil
}

func (a *App) RunDB() error {
	dbConn, err := store.NewDBConnect()

	if err != nil {
		return fmt.Errorf("error on creating store: %v", err)
	}

	fmt.Println("Postgres is connected...")

	a.DB = dbConn
	return nil
}

func (a *App) RunHTTP() {
	http.HandleFunc("/", a.HandlerWebApp)
	http.HandleFunc("/order/get", a.HandlerGetData)
	http.HandleFunc("/order/add", a.HandlerAddData)

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		fmt.Printf("error on creating http-server: %v", err)
		return
	}

	fmt.Printf("HTTP server created")
}

func (a *App) CreateRepository() error {
	rep := store.NewRepository(a.DB)

	err := rep.LoadAllRepositories()
	if err != nil {
		return fmt.Errorf("error on loading repositories: %v", err)
	}

	a.Rep = rep
	return nil
}

func (a *App) RunNats() {
	sc := *a.Sc

	_, err := sc.Subscribe("json-parser", a.HandleJSON)

	if err != nil {
		fmt.Printf("\nerror on subscribe: %v", err)
		return
	}

	fmt.Println("Nats streaming is listening...")
}

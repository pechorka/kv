package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/pechorka/kv/app/kv-api/handlers"
	"github.com/pechorka/kv/internal/store"
	"github.com/pkg/errors"
)

func main() {
	log := log.New(os.Stdout, "KV ", log.LstdFlags|log.Lshortfile)
	if err := run(log); err != nil {
		log.Println("остановка сервиса", "error: ", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {
	log.Println("main : запуск сервиса")
	defer log.Println("main : остановка сервиса")

	var port int
	flag.IntVar(&port, "port", 8080, "порт на котором будет запущен сервис")
	flag.Parse()

	store := store.NewMemory()

	api := http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: handlers.API(log, store),
	}

	serverErr := make(chan error)

	go func() {
		log.Println("api слушает на ", port)
		serverErr <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return errors.Wrap(err, "при запуске сервера")

	case <-shutdown:
		log.Println("main : начало остановки сервиса")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : не получилось остановить сервис за 5 секунд: %v", err)
			err = api.Close()
		}

		if err != nil {
			return errors.Wrap(err, "не получилось остановить сервис")
		}
	}

	return nil
}

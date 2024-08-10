package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/adettelle/accumulative-loyalty-system/internal/gophermart/config"
	"github.com/adettelle/accumulative-loyalty-system/internal/gophermart/database"
	"github.com/adettelle/accumulative-loyalty-system/internal/gophermart/server/api"
	"github.com/adettelle/accumulative-loyalty-system/pkg/mware"
	"github.com/go-chi/chi/v5"
)

func main() {
	var uri string

	config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if config.DBUri != "" {
		uri = config.DBUri
	}

	db, err := database.Connect(uri)
	if err != nil {
		log.Fatal(err)
	}

	err = database.CreateTable(db, context.Background())
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	storage := &api.DBStorage{
		DB:  db,
		Ctx: context.Background(),
	}

	fmt.Println("Starting server")
	address := config.Address //"localhost:8080"

	r := chi.NewRouter()

	// конфигурирование сервера
	r.Post("/api/user/login", storage.Login)
	r.Post("/api/user/orders", mware.AuthMwr(storage.AddOrder))
	r.Get("/api/user/orders", mware.AuthMwr(storage.GetOrders))
	r.Get("/api/user/balance", mware.AuthMwr(storage.GetBalance))
	r.Post("/api/user/balance/withdraw", mware.AuthMwr(storage.PostWithdraw))
	r.Get("/api/user/withdrawals", mware.AuthMwr(storage.GetWithdrawals))

	err = http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal(err)
	}
}

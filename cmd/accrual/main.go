package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/adettelle/accumulative-loyalty-system/internal/accrual/accrdb"
	"github.com/adettelle/accumulative-loyalty-system/internal/accrual/accrserver/accrapi"
	"github.com/go-chi/chi/v5"
)

func main() {
	ps := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		`localhost`, `5433`, `postgres`, `password`, `accrual-service`)

	db, err := accrdb.Connect(ps)
	if err != nil {
		log.Fatal(err)
	}

	err = accrdb.CreateTable(db, context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	storage := &accrapi.DBStorage{
		DB:  db,
		Ctx: context.Background(),
	}

	fmt.Println("Starting server")
	address := "localhost:8081"

	r := chi.NewRouter()

	// конфигурирование сервера
	r.Get("/api/orders/{number}", storage.GetOrderByNumber)
	r.Post("/api/orders", storage.RegisterOrder)
	r.Post("/api/goods", storage.RegisterRewardType)

	err = http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal(err)
	}
}

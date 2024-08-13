package accrapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/adettelle/accumulative-loyalty-system/internal/accrual/accrmodel"
	"github.com/adettelle/accumulative-loyalty-system/pkg/luhn"
	"github.com/go-chi/chi/v5"
)

type DBStorage struct {
	Ctx context.Context
	DB  *sql.DB
}

type OrderResponse struct {
	Number string  `json:"order"`
	Status string  `json:"status"`
	Amount float64 `json:"accrual"`
}

func NewOpderResp(order accrmodel.Order) OrderResponse {
	return OrderResponse{
		Number: order.Number,
		Status: order.Status,
		Amount: order.Amount,
	}
}

// 204 — заказ не зарегистрирован в системе расчёта - это на ошибку emptyRow или если статус не проставлен???
// 429 — превышено количество запросов к сервису.
func (s *DBStorage) GetOrderByNumber(w http.ResponseWriter, r *http.Request) {
	numberToSearch := chi.URLParam(r, "number")

	// проверка Luhn
	if !luhn.CheckLuhn(numberToSearch) {
		w.WriteHeader(http.StatusUnprocessableEntity) // неверный формат номера заказа
		return
	}

	number, err := strconv.Atoi(numberToSearch)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// ord := OrderResponse{}

	order, err := accrmodel.GetOrderByNumber(s.DB, s.Ctx, number)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError) // ошибка с БД
		return
	}

	resp, err := json.Marshal(NewOpderResp(order))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *DBStorage) RegisterOrder(w http.ResponseWriter, r *http.Request) {

}

func (s *DBStorage) RegisterRewardType(w http.ResponseWriter, r *http.Request) {

}

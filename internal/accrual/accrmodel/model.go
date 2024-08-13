package accrmodel

import (
	"context"
	"database/sql"
)

const (
	RewardTypePercent = "percent"
	RewardTypePoints  = "points"
	StatusRegistered  = "registered"
	StatusProcessing  = "processing"
	StatusInvalid     = "invalid"
	StatusProcessed   = "processed"
)

type Order struct {
	ID     int
	Number string
	// Product string
	// Price   float64
	Status string
	Amount float64
}

// Чтение строки по заданному номеру заказа.
// Из таблицы должна вернуться только одна строка.
func GetOrderByNumber(db *sql.DB, ctx context.Context, number int) (Order, error) {
	ord := Order{}

	sqlSt := `select * from order_item where order_number = $1;`

	row := db.QueryRowContext(ctx, sqlSt, number)

	err := row.Scan(&ord.ID, &ord.Number, &ord.Status, &ord.Amount)

	return ord, err
}

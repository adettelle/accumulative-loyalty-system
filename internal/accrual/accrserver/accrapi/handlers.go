package accrapi

import (
	"context"
	"database/sql"
	"net/http"
)

type DBStorage struct {
	Ctx context.Context
	DB  *sql.DB
}

func (s *DBStorage) RegisterOrder(w http.ResponseWriter, r *http.Request) {

}

func (s *DBStorage) RegisterRewardType(w http.ResponseWriter, r *http.Request) {

}

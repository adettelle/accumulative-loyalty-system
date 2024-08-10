package accrdb

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func CreateTable(db *sql.DB, ctx context.Context) error {
	// rewardType := `create type reward_type_enum as enum ('percent', 'points');`
	// _, err := db.ExecContext(ctx, rewardType, accrmodel.RewardTypePercent, accrmodel.RewardTypePoints)
	// if err != nil {
	// 	return err
	// }

	sqlStReward := `create table if not exists reward 
		(id serial primary key, 
		"match" varchar(60) not null,
		reward_type reward_type_enum not null);`

	_, err := db.ExecContext(ctx, sqlStReward)
	if err != nil {
		return err
	}

	sqlStOrderItem := `create table if not exists order_item 
		(id serial primary key, 
		order_number integer,
		price double precision,
		reward_processed boolean,
		reward_amount double precision);`

	_, err = db.ExecContext(ctx, sqlStOrderItem)
	if err != nil {
		return err
	}

	return nil
}

func Connect(dbParams string) (*sql.DB, error) {
	log.Println("Connecting to DB accrual service", dbParams)
	db, err := sql.Open("pgx", dbParams)
	if err != nil {
		return nil, err
	}
	return db, nil
}

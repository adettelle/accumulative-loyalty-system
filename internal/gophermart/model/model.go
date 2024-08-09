package model

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	StatusNew             = "new"
	StatusProcessing      = "processing"
	StatusInvalid         = "invalid"
	StatusProcessed       = "processed"
	TransactionAccrual    = "accrual"
	TransactionWithdrawal = "withdrawal"
)

// func CreateTable(db *sql.DB, ctx context.Context) error { // user
// 	sqlStCustomer := `create table if not exists customer
// 		(id serial primary key,
// 		first_name varchar(30) not null,
// 		last_name varchar(30) not null,
// 		email varchar(100) not null,
// 		phone varchar(30) not null,
// 		login varchar(100),
// 		password varchar(255),
// 		created_at timestamp not null default now(),
// 		unique(phone, email));`

// 	_, err := db.ExecContext(ctx, sqlStCustomer)
// 	if err != nil {
// 		return err
// 	}

// 	// statusType := `create type status_type_enum as enum ('new', 'processing', 'invalid', 'processed');`
// 	// _, err = db.ExecContext(ctx, statusType)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	sqlStOrder := `create table if not exists "order"
// 		(id serial primary key,
// 		customer_id integer,
// 		number text,
// 		status status_type_enum not null,
// 		created_at timestamp not null default now(),
// 		foreign key (customer_id) references customer (id),
// 		unique(number, customer_id));`

// 	_, err = db.ExecContext(ctx, sqlStOrder)
// 	if err != nil {
// 		return err
// 	}

// 	// начисления и списания
// 	// transactionType := `create type transaction_type_enum as enum ('accrual', 'withdrawal');`
// 	// _, err = db.ExecContext(ctx, transactionType)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	sqlStLoyalty := `create table if not exists loyalty_system
// 		(id serial primary key,
// 		customer_id integer,
// 		order_id integer,
// 		points double precision,
// 		transacton transaction_type_enum,
// 		created_at timestamp not null default now(),
// 		unique(customer_id, order_id),
// 		foreign key (customer_id) references customer (id),
// 		foreign key (order_id) references "order" (id));`

// 	_, err = db.ExecContext(ctx, sqlStLoyalty)
// 	if err != nil {
// 		return err
// 	}

// 	sqlStProduct := `create table if not exists product
// 		(id serial primary key,
// 		title varchar(60),
// 		price integer);`

// 	_, err = db.ExecContext(ctx, sqlStProduct)
// 	if err != nil {
// 		return err
// 	}

// 	sqlStOrderProduct := `create table if not exists order_product
// 		(id serial primary key,
// 		order_id integer references "order" (id),
// 		product_id integer references product (id),
// 		amount integer);`

// 	_, err = db.ExecContext(ctx, sqlStOrderProduct)
// 	if err != nil {
// 		return err
// 	}

// 	// rewardType := `create type reward_type_enum as enum ('percent', 'points');`
// 	// _, err = db.ExecContext(ctx, rewardType)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	sqlStReward := `create table if not exists reward
// 		(id serial primary key,
// 		title varchar(60),
// 		product_id integer references product (id),
// 		description varchar(255),
// 		reward_type reward_type_enum not null);`

// 	_, err = db.ExecContext(ctx, sqlStReward)
// 	if err != nil {
// 		return err
// 	}

// 	sqlStRewardSystem := `create table if not exists reward_system
// 		(id serial primary key,
// 		order_id integer references "order" (id),
// 		points double precision);`

// 	_, err = db.ExecContext(ctx, sqlStRewardSystem)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func Connect(dbParams string) (*sql.DB, error) {
// 	log.Println("Connecting to DB", dbParams)
// 	db, err := sql.Open("pgx", dbParams)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// defer db.Close()

// 	// делаем запрос
// 	row := db.QueryRowContext(context.Background(), "SELECT 1;")
// 	// готовим переменную для чтения результата
// 	var id int64
// 	err = row.Scan(&id) // разбираем результат
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(id)

// 	return db, nil
// }

// если возврат (0, err) - это значит, что юзера с таким заказом нет
func GetUserByOrderOld(numOrder string, db *sql.DB, ctx context.Context) (int, error) {
	sqlSt := `select customer_id from "order" where "number" = $1;`
	row := db.QueryRowContext(ctx, sqlSt, numOrder)

	var id int

	err := row.Scan(&id)
	log.Println("err:", err)
	log.Println("id:", id)
	// if no rows there is no session
	if err == sql.ErrNoRows {
		return 0, nil
	}
	// an error other than no rows was returned, return with error
	return id, err
}

// GetUserByOrder возвращает id юзера и ошибку
func GetUserByOrder(numOrder string, db *sql.DB, ctx context.Context) (int, error) {
	sqlSt := `select customer_id from "order" where "number" = $1;`
	row := db.QueryRowContext(ctx, sqlSt, numOrder)

	var id int
	err := row.Scan(&id)

	return id, err
}

func GetLoginByID(id int, db *sql.DB, ctx context.Context) (string, error) {
	sqlSt := `select login from customer where id = $1;`
	row := db.QueryRowContext(ctx, sqlSt, id)

	var login string
	err := row.Scan(&login)

	return login, err
}

func OrderExists(numOrder string, db *sql.DB, ctx context.Context) (bool, error) {
	sqlSt := `select count(id) > 0 as order_exists from "order" where "number" = $1;`
	row := db.QueryRowContext(ctx, sqlSt, numOrder)

	var ordExists bool

	err := row.Scan(&ordExists)
	log.Println("error in OrderExists:", err)
	log.Println("ordExists:", ordExists)
	return ordExists, err
}

// // проверяем, есть ли пользователь с таким номером заказа
// func UserHasOrder(numOrder string, userId int, db *sql.DB, ctx context.Context) (bool, error) {
// 	ordExists, err := OrderExists(numOrder, db, ctx)
// 	if !ordExists {
// 		return false, err
// 	}

// 	userIdByGet, err := GetUserByOrder(numOrder, db, ctx)
// 	if userIdByGet != userId { // такой номера заказа уже есть у другого пользователя
// 		return
// 	}

// 	// если 0, err - это значит, что юзера с таким заказом нет
// 	userIdByGet, err := GetUserByOrder(numOrder, db, ctx)
// 	if userIdByGet == 0 { // такого номера заказа у пользователя нет
// 		log.Printf("There is no user with order number %s", numOrder)
// 		return userIdByGet, false, err
// 	}
// 	if userIdByGet != userId { // такой номера заказа уже есть у другого пользователя
// 		log.Printf("There is a user %d with order number %s", userIdByGet, numOrder)
// 		return userIdByGet, false, err
// 	}
// 	return userIdByGet, true, nil // такой номера заказа есть у проверяемого пользователя
// }

// ????????????????/
// проверяем, есть ли пользователь с таким номером заказа
// UserHasOrder возвращает id юзера, bool, err
// 0, false, nil - такого номера заказа ни у кого нет
// id, false, err - такой номера заказа уже есть у другого пользователя
// id, true, err - такой номера заказа есть у проверяемого пользователя
func UserHasOrder(numOrder string, userId int, db *sql.DB, ctx context.Context) (int, bool, error) {
	// если 0, err=nil - это значит, что юзера с таким заказом нет
	userIdByGet, err := GetUserByOrder(numOrder, db, ctx)
	if err == nil { // такого номера заказа ни у кого нет
		log.Printf("There is no user with order number %s", numOrder)
		return userIdByGet, false, nil
	}
	// if userIdByGet == 0 { // такого номера заказа у пользователя нет
	// 	log.Printf("There is no user with order number %s", numOrder)
	// 	return userIdByGet, false, err
	// }
	if userIdByGet != userId { // такой номера заказа уже есть у другого пользователя
		log.Printf("There is a user %d with order number %s", userIdByGet, numOrder)
		return userIdByGet, false, err
	}
	return userIdByGet, true, err // такой номера заказа есть у проверяемого пользователя
}

// func AddOrder(numOrder string, userId int, db *sql.DB, ctx context.Context) error {
// 	log.Println("Writing to DB")
// 	sqlSt := `insert into "order" (customer_id, number, status)
// 			values ($1, $2, $3)` // on conflict

// 	_, err := db.ExecContext(ctx, sqlSt, 1, numOrder, StatusNew)
// 	if err != nil {
// 		log.Println("error in adding order:", err)
// 		return err
// 	}
// 	return nil
// }

type Order struct {
	Id       string // почему string???????????????
	Number   string
	Status   string
	Points   float64
	CratedAt time.Time
}

type Customer struct {
	Id        int // почему string???????????????
	FirstName string
	LastName  string
	Email     string
	Phone     string
	// Roles
	// IsDeleted
}

func GetOrders(db *sql.DB, ctx context.Context) ([]Order, error) {
	orders := make([]Order, 0)

	sqlSt := `select ord.id, "number", status, ls.points, ord.created_at 
		from "order" ord
		join loyalty_system ls 
		on ls.order_id = ord.id`

	rows, err := db.QueryContext(ctx, sqlSt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var ord Order
		err := rows.Scan(&ord.Id, &ord.Number, &ord.Status, &ord.Points, &ord.CratedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, ord)
	}
	return orders, nil
}

// GetAccrualPoints показывает количество набранных баллов пользователя
func GetAccrualPoints(userID int, db *sql.DB, ctx context.Context) (float64, error) {
	// userID := 1
	sqlSt := `select ls.points from loyalty_system ls 
		join customer c 
		on c.id = ls.customer_id 
		where c.id = $1 and ls.transacton = $2;` // 'accrual'

	row := db.QueryRowContext(ctx, sqlSt, userID, TransactionAccrual)

	var pointsAccrual float64

	err := row.Scan(&pointsAccrual)
	if err != nil {
		log.Printf("Error in getting balance of user %d", userID)
		return 0, err
	}

	return pointsAccrual, nil
}

// GetWithdrawalPoints показывает количество потраченных баллов пользователя
func GetWithdrawalPoints(userID int, db *sql.DB, ctx context.Context) (float64, error) {
	sqlSt := `select ls.points from loyalty_system ls 
		join customer c 
		on c.id = ls.customer_id 
		where c.id = $1 and ls.transacton = $2;` // 'withdrawal'

	row := db.QueryRowContext(ctx, sqlSt, userID, TransactionWithdrawal)

	var pointsWithdrawal float64

	err := row.Scan(&pointsWithdrawal)
	if err != nil {
		log.Printf("Error in getting balance of user %d", userID)
		return 0, err
	}

	return pointsWithdrawal, nil
}

// как передать пользователя???????????????
// Withdraw списывает баллы sum с номера счета order у зарегистрированного пользователя
func Withdraw(order string, sum float64, db *sql.DB, ctx context.Context) error {
	userID := 1
	sqlSt := `update loyalty_system ls
		set points = points - $1 
		where customer_id = $2 and transacton = 'accrual';`

	_, err := db.ExecContext(ctx, sqlSt, sum, userID)
	if err != nil {
		return err
	}

	return nil
}

// транзакция списания
type TransactionW struct {
	// Id       string
	OrderNumber string
	// Status   string
	Points   float64
	CratedAt time.Time
}

// AllWithdrawals показывает все транзакции с выводом средств
func AllWithdrawals(db *sql.DB, ctx context.Context) ([]TransactionW, error) {
	transactions := make([]TransactionW, 0)
	sqlSt := `select ord."number", ls.points, ls.created_at 
		from loyalty_system ls 
		join "order" ord
		on ord.id = ls.order_id
		where ls.transacton = $1
		order by created_at desc;`

	rows, err := db.QueryContext(ctx, sqlSt, TransactionWithdrawal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var tr TransactionW
		err := rows.Scan(&tr.OrderNumber, &tr.Points, &tr.CratedAt)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, tr)
	}
	return transactions, nil
}

func GetCustomerByLogin(login string, db *sql.DB, ctx context.Context) (*Customer, error) {
	sqlSt := `select id, first_name, last_name, email, phone from customer where login = $1;`

	row := db.QueryRowContext(ctx, sqlSt, login)

	var customer Customer

	err := row.Scan(&customer.Id, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // считаем, что это не ошибка, просто не нашли пользователя
		}
		return nil, err
	}
	return &customer, nil
}

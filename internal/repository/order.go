package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Order struct {
	Order_uid          string    `json:"order_uid"`
	Track_number       string    `json:"track_number"`
	Entry              string    `json:"entry"`
	Delivery           Delivery  `json:"delivery"`
	Payment            Payment   `json:"payment"`
	Items              []Item    `json:"items"`
	Locale             string    `json:"locale"`
	Internal_signature string    `json:"internal_signature"`
	Customer_id        string    `json:"customer_id"`
	Delivery_service   string    `json:"delivery_service"`
	Shardkey           string    `json:"shardkey"`
	Sm_id              int       `json:"sm_id"`
	Date_created       time.Time `json:"date_created"`
	Oof_shard          string    `json:"oof_shard"`
}

type RepositoryOrder struct {
	elements map[string]Order
	db       *sql.DB
	Rep      *Repository
	mx       *sync.RWMutex
}

const (
	queryCreateOrder = `
		INSERT into orders 
			(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, delivery)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	queryInsertPayment = `
		INSERT into payment
			(transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	querySelectOrder = `
		SELECT 
			order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, delivery
		FROM 
			orders
	`
)

func NewOrderRep(s *sql.DB) *RepositoryOrder {
	return &RepositoryOrder{
		db:       s,
		elements: map[string]Order{},
		mx:       &sync.RWMutex{},
	}
}

func (r *RepositoryOrder) LoadAllElements() error {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.Rep.Items.mx.RLock()
	defer r.Rep.Items.mx.RUnlock()
	r.Rep.Payments.mx.RLock()
	defer r.Rep.Payments.mx.RUnlock()

	ctx := context.Background()
	rows, err := r.db.QueryContext(ctx, querySelectOrder)

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		o := Order{}
		var val string
		err := rows.Scan(&o.Order_uid, &o.Track_number, &o.Entry, &o.Locale, &o.Internal_signature, &o.Customer_id, &o.Delivery_service, &o.Shardkey, &o.Sm_id, &o.Date_created, &o.Oof_shard, &val)
		if err != nil {
			return fmt.Errorf("error on rows scan: %v", err)
		}

		if err = json.Unmarshal([]byte(val), &o.Delivery); err != nil {
			return fmt.Errorf("error on json unmarshal delivery: %v", err)
		}

		arr, ok := r.Rep.Items.elements[o.Track_number]
		if !ok {
			o.Items = make([]Item, 0)
		} else {
			o.Items = arr
		}

		pay, ok := r.Rep.Payments.elements[o.Order_uid]
		if !ok {
			return fmt.Errorf("payment is not created: %v", o.Order_uid)
		}

		o.Payment = pay
		r.elements[o.Order_uid] = o
	}

	return nil
}

func (r *RepositoryOrder) CreateByJSON(ctx context.Context, jorder []byte) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	newOrder := &Order{}

	if err := json.Unmarshal(jorder, newOrder); err != nil {
		return err
	}

	return r.Insert(ctx, *newOrder)
}

func (r *RepositoryOrder) Insert(ctx context.Context, o Order) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return fmt.Errorf("error on creating order: %w", err)
	}

	val, err := json.Marshal(o.Delivery)
	if err != nil {
		return fmt.Errorf("error on marshalling delivery: %v", err)
	}

	_, err = tx.ExecContext(ctx, queryCreateOrder, o.Order_uid, o.Track_number, o.Entry, o.Locale, o.Internal_signature, o.Customer_id, o.Delivery_service, o.Shardkey, o.Sm_id, o.Date_created, o.Oof_shard, val)

	if err != nil {
		return fmt.Errorf("error on create order: %w", err)
	}

	payment := o.Payment

	_, err = tx.ExecContext(ctx, queryInsertPayment, payment.Transaction, payment.Request_id, payment.Currency, payment.Provider, payment.Amount, payment.Payment_dt, payment.Bank, payment.Delivery_cost, payment.Goods_total, payment.Custom_fee)

	if err != nil {
		return err
	}

	for _, v := range o.Items {
		if err = r.Rep.Items.Insert(ctx, tx, v); err != nil {
			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	r.elements[o.Order_uid] = o

	return nil
}

func (r *RepositoryOrder) GetJSONByID(id []byte) ([]byte, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	el, err := r.getOrderByID(string(id))
	if err != nil {
		return []byte{}, err
	}

	val, err := json.Marshal(*el)
	if err != nil {
		return []byte{}, err
	}

	return val, nil
}

func (r *RepositoryOrder) getOrderByID(id string) (*Order, error) {
	order, ok := r.elements[id]
	if !ok {
		return nil, fmt.Errorf("order by %v id doesn't exists", id)
	}

	return &order, nil
}

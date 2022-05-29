package store

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

type Payment struct {
	Transaction   string `json:"transaction"`
	Request_id    string `json:"request_id"`
	Currency      string `json:"currency"`
	Provider      string `json:"provider"`
	Amount        int    `json:"amount"`
	Payment_dt    int    `json:"payment_dt"`
	Bank          string `json:"bank"`
	Delivery_cost int    `json:"delivery_cost"`
	Goods_total   int    `json:"goods_total"`
	Custom_fee    int    `json:"custom_fee"`
}

const (
	querySelectPayments = `
		SELECT 
			transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
	 	FROM 
		 	payment`
)

type RepositoryPayments struct {
	db       *sql.DB
	elements map[string]Payment
	Rep      *Repository
	mx       *sync.RWMutex
}

func NewPaymentsRep(s *sql.DB) *RepositoryPayments {
	return &RepositoryPayments{
		db:       s,
		elements: map[string]Payment{},
		mx:       &sync.RWMutex{},
	}
}

func (r *RepositoryPayments) LoadAllElements() error {
	r.mx.Lock()
	defer r.mx.Unlock()

	ctx := context.Background()
	rows, err := r.db.QueryContext(ctx, querySelectPayments)

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		p := Payment{}
		err := rows.Scan(&p.Transaction, &p.Request_id, &p.Currency, &p.Provider, &p.Amount, &p.Payment_dt, &p.Bank, &p.Delivery_cost, &p.Goods_total, &p.Custom_fee)
		if err != nil {
			return fmt.Errorf("error on rows scan: %v", err)
		}

		r.elements[p.Transaction] = p
	}

	return nil
}

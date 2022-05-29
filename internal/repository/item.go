package store

import (
	"context"
	"database/sql"
	"sync"

	_ "github.com/lib/pq"
)

type Item struct {
	Chrt_id      int    `json:"chrt_id"`
	Track_number string `json:"track_number"`
	Price        int    `json:"price"`
	Rid          string `json:"rid"`
	Name         string `json:"name"`
	Sale         int    `json:"sale"`
	Size         string `json:"size"`
	Total_price  int    `json:"total_price"`
	Nm_id        int    `json:"nm_id"`
	Brand        string `json:"brand"`
	Status       int    `json:"status"`
}

type RepositoryItems struct {
	db       *sql.DB
	elements map[string][]Item
	Rep      *Repository
	mx       *sync.RWMutex
}

const (
	queryInsertItem = `
		INSERT INTO 
			items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	querySelectAllItems = `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items`
)

func NewItemsRep(s *sql.DB) *RepositoryItems {
	return &RepositoryItems{
		db:       s,
		elements: map[string][]Item{},
		mx:       &sync.RWMutex{},
	}
}

func (r *RepositoryItems) LoadAllElements() error {
	r.mx.Lock()
	defer r.mx.Unlock()
	ctx := context.Background()
	rows, err := r.db.QueryContext(ctx, querySelectAllItems)

	if err != nil {
		return err
	}

	for rows.Next() {
		item := Item{}

		err := rows.Scan(&item.Chrt_id, &item.Track_number, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.Total_price, &item.Nm_id, &item.Brand, &item.Status)
		if err != nil {
			return err
		}

		r.AppendElementByTrackNumber(item)
	}

	return nil
}

func (r *RepositoryItems) AppendElementByTrackNumber(v Item) {
	key := v.Track_number
	arr, ok := r.elements[key]
	if !ok {
		r.elements[key] = make([]Item, 0, 1)
		r.elements[key] = append(r.elements[key], v)
		return
	}
	r.elements[key] = append(arr, v)
}

func (r *RepositoryItems) Insert(ctx context.Context, tx *sql.Tx, v Item) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	_, err := tx.ExecContext(ctx, queryInsertItem, v.Chrt_id, v.Track_number, v.Price, v.Rid, v.Name, v.Sale, v.Size, v.Total_price, v.Nm_id, v.Brand, v.Status)

	if err != nil {
		return err
	}

	r.AppendElementByTrackNumber(v)

	return nil
}

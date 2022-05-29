package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Repository struct {
	Items    *RepositoryItems
	Orders   *RepositoryOrder
	Payments *RepositoryPayments
}

func NewDBConnect() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewRepository(s *sql.DB) *Repository {
	globalRep := &Repository{}

	itemsRep := NewItemsRep(s)
	orderRep := NewOrderRep(s)
	paymentRep := NewPaymentsRep(s)

	itemsRep.Rep = globalRep
	orderRep.Rep = globalRep
	paymentRep.Rep = globalRep

	globalRep.Items = itemsRep
	globalRep.Orders = orderRep
	globalRep.Payments = paymentRep

	return globalRep
}

func (g *Repository) GetOrderRep() *RepositoryOrder {
	return g.Orders
}

func (g *Repository) GetItemRep() *RepositoryItems {
	return g.Items
}

func (g *Repository) LoadAllRepositories() error {
	err := g.Items.LoadAllElements()
	if err != nil {
		return fmt.Errorf("error on loading items: %v", err)
	}

	err = g.Payments.LoadAllElements()
	if err != nil {
		return fmt.Errorf("error on loading payments: %v", err)
	}

	err = g.Orders.LoadAllElements()
	if err != nil {
		return fmt.Errorf("error on loading orders: %v", err)
	}
	return nil
}

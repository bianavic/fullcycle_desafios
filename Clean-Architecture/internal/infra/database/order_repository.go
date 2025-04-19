package database

import (
	"database/sql"
	"github.com/bianavic/fullcycle_clean-architecture/internal/entity"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	smt, err := r.DB.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = smt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) ListOrders() ([]entity.Order, error) {
	rows, err := r.DB.Query("SELECT id, price, tax, final_price FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		scanErr := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
		if scanErr != nil {
			return nil, scanErr
		}

		orders = append(orders, order)
	}
	return orders, nil
}

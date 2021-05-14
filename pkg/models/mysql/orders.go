package mysql

import (
	"database/sql"
	"errors"

	"ProjectGoLive/pkg/models"
)

// Define a ProductModel type which wraps a sql.DB connection pool.
type OrderModel struct {
	DB *sql.DB
}

// Create inserts a new row into the orders table
// using information from a CartItem.
func (m *OrderModel) Create(c *models.CartItem) error {
	stmt :=
		`INSERT INTO orders (UserID, ProductID,	Qty, SellerID, Status)
		 VALUES (?, ?, ?, (SELECT SellerID FROM Product WHERE ProductID = ?), 0);`

	_, err := m.DB.Exec(stmt, c.UserID, c.Product.ProductID, c.Qty, c.Product.ProductID)
	if err != nil {
		return err
	}

	return nil
}

// GetAll retrieves all the orders which has the
// specified UserID column value. GetAll returns
// different resuls depending if the userID provided
// refers to a seller.
func (m *OrderModel) GetAll(userID string, isSeller bool) ([]*models.Orders, error) {
	var stmt string
	if !isSeller {
		stmt = `SELECT 
					Orders.OrderID,	Orders.UserID, Product.Name,
					Orders.Qty,	Orders.SellerID, Orders.Status
				FROM Orders 
				LEFT JOIN Product ON Orders.ProductID = Product.ProductID
				WHERE Orders.UserID=?`
	} else {
		stmt = `SELECT
					Orders.OrderID, Orders.UserID, Product.Name,
					Orders.Qty, Orders.SellerID, Orders.Status
				FROM Orders 
				LEFT JOIN Product ON Orders.ProductID = Product.ProductID
				WHERE Orders.SellerID=?`
	}

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []*models.Orders{}
	for rows.Next() {
		order := &models.Orders{}
		err = rows.Scan(
			&order.OrderID,
			&order.UserID,
			&order.Product.Name,
			&order.Qty,
			&order.SellerID,
			&order.Status,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// UpdateStatus edits the status column value
// for the row which has the specified orderID
// column value.
func (m *OrderModel) UpdateStatus(orderid, status int) error {
	var stmt string
	if status == 1 {
		stmt = `UPDATE Product, Orders 
	SET Orders.Status = ?, Product.Inventory = Product.Inventory - Orders.Qty, Product.UnitSold = Product.UnitSold + Orders.Qty
	WHERE Orders.ProductID = Product.ProductID AND Orders.OrderID = ?`
	} else {
		stmt = `UPDATE Orders SET Status = ? Where OrderID = ?`
	}
	_, err := m.DB.Exec(stmt, status, orderid)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves an order which has the specified
// OrderID column value.
func (m *OrderModel) Get(orderID int) (*models.Orders, error) {

	stmt := `SELECT
	Orders.OrderID, Orders.UserID, Product.Name,
	Orders.Qty, Orders.SellerID, Orders.Status
	FROM Orders 
	LEFT JOIN Product ON Orders.ProductID = Product.ProductID
	WHERE Orders.OrderID=?`

	row := m.DB.QueryRow(stmt, orderID)

	order := &models.Orders{}
	err := row.Scan(
		&order.OrderID,
		&order.UserID,
		&order.Product.Name,
		&order.Qty,
		&order.SellerID,
		&order.Status,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return order, nil
}

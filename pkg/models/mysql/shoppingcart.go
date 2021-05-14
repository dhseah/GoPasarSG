package mysql

import (
	"ProjectGoLive/pkg/models"
	"database/sql"
)

// Define a CartModel type which wraps a sql.DB connection pool.
type CartModel struct {
	DB *sql.DB
}

// InsertItem inserts a new row into the table.
func (m *CartModel) InsertItem(userID string, productID int) error {
	stmt :=
		`INSERT INTO shoppingCart (UserID, ProductID, Qty) 
		VALUES (?,?,1)
		ON DUPLICATE KEY UPDATE Qty = Qty+1;`

	_, err := m.DB.Exec(stmt, userID, productID)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves every row in the table that
// has the specified UserID column value and
// returns the result.
func (m *CartModel) Get(userID string) ([]*models.CartItem, error) {
	stmt :=
		`SELECT
			Product.Name, Product.Inventory, ShoppingCart.Qty,
			Product.Price, Product.DiscountID, Product.SellerID,
			Product.ProductID
		FROM ShoppingCart
		LEFT JOIN Product ON ShoppingCart.ProductID = Product.ProductID
		WHERE  ShoppingCart.UserID = ?;`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shoppingCart := []*models.CartItem{}
	for rows.Next() {
		cartItem := &models.CartItem{}
		err = rows.Scan(
			&cartItem.Product.Name,
			&cartItem.Product.Inventory,
			&cartItem.Qty,
			&cartItem.Product.Price,
			&cartItem.Product.DiscountID,
			&cartItem.Product.SellerID,
			&cartItem.Product.ProductID,
		)
		if err != nil {
			return nil, err
		}
		shoppingCart = append(shoppingCart, cartItem)
	}

	return shoppingCart, nil
}

// Update updates the Qty column value in the
// table for the row that has the specified
// UserID & ProductID column values.
func (m *CartModel) Update(quantity, productid int, userID string) error {
	stmt := `UPDATE ShoppingCart 
			 SET Qty= ? 
			 WHERE (UserID = ? AND ProductID = ?)`

	res, err := m.DB.Exec(stmt, quantity, userID, productid)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return models.ErrNoRowsAffected
	}

	return nil
}

// DeleteItem deletes a row in the table that
// has the specified UserID & ProductID
// column value.
func (m *CartModel) DeleteItem(userID string, productID int) error {
	stmt := `Delete FROM shoppingCart 
			 WHERE (UserID = ? AND ProductID = ?)`

	_, err := m.DB.Exec(stmt, userID, productID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAll deletes every row in the table
// that has the specified UserID column value.
func (m *CartModel) DeleteAll(userID string) error {
	stmt := `Delete FROM shoppingCart
			 WHERE UserID = ?`

	_, err := m.DB.Exec(stmt, userID)
	if err != nil {
		return err
	}

	return nil
}

// CheckOut checks if the quantity column value
// for every row that has the specified UserID
// column value is less than the iventory column
// value in the product table via the ProductID
// foreign key.
func (m *CartModel) CheckOut(userid string) ([]*models.CartItem, error) {
	stmt :=
		`SELECT
			shoppingcart.UserID,
			shoppingcart.ProductID,
			product.Name,
			product.DiscountID,
			product.Inventory,
			product.Price,
			shoppingcart.Qty,
			CASE WHEN shoppingcart.Qty < product.Inventory
				THEN '0'
				ELSE '1'
			END
    		As Valid
		FROM shoppingcart
		LEFT JOIN product ON
			shoppingcart.ProductID = product.ProductID
		WHERE shoppingcart.UserID = ?;`

	rows, err := m.DB.Query(stmt, userid)
	if err != nil {
		return nil, err
	}

	var res []*models.CartItem
	for rows.Next() {
		i := &models.CartItem{}
		err := rows.Scan(
			&i.UserID,
			&i.Product.ProductID,
			&i.Product.Name,
			&i.Product.DiscountID,
			&i.Product.Inventory,
			&i.Product.Price,
			&i.Qty,
			&i.Invalid)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}

	return res, nil
}

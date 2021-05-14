package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"ProjectGoLive/pkg/models"
)

// Define a ProductModel type which wraps a sql.DB connection pool.
type ProductModel struct {
	DB *sql.DB
}

// Create inserts a new product into the database
// and return its ID.
func (m *ProductModel) Create(name, desc, keyword, sellerID string, price float64, inventory, catID, discID int) (int, error) {
	stmt := `INSERT INTO Product 
			 (Name, Description, Keyword, CategoryID, Price, DiscountID, Inventory, SellerID)
    		 VALUES(?, ?, ?, ?, ?, ?,?, ?)`

	result, err := m.DB.Exec(stmt, name, desc, keyword, catID, price, discID, inventory, sellerID)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Update edits the column values for a row which
// has the specified ProductID column value.
func (m *ProductModel) Update(name, desc, keyword string, price float64, inventory, ID, catID, discID int) error {
	stmt := `UPDATE Product 
			 SET Name=?, Description=?, Keyword=?, Price=?, CategoryID=?, DiscountID=?, inventory=?
			 WHERE ProductID = ?`

	_, err := m.DB.Exec(stmt, name, desc, keyword, price, catID, discID, inventory, ID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Get retrieves a row which has the specified
// ProductID column value.
func (m *ProductModel) Get(id int) (*models.Product, error) {
	stmt :=
		`SELECT 
			ProductID, Name, Description, Keyword, CategoryID, Price, DiscountID,
			Inventory, Rating, RatingNum, UnitSold, SellerID, Created, Modified 
		FROM product
		WHERE ProductID = ?`

	p := &models.Product{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(
		&p.ProductID, &p.Name, &p.Desc, &p.Keyword, &p.CategoryID,
		&p.Price, &p.DiscountID, &p.Inventory, &p.Rating,
		&p.RatingNum, &p.UnitSold, &p.SellerID, &p.Created,
		&p.Modified,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return p, nil
}

// GetAll retrieves every row in the product table.
func (m *ProductModel) GetAll() ([]*models.Product, error) {
	stmt :=
		`SELECT 
			ProductID, Name, Description, CategoryID, Price, DiscountID,
			Inventory, Rating, RatingNum, UnitSold, SellerID, Created, Modified 
		FROM product`

	products := []*models.Product{}

	rows, _ := m.DB.Query(stmt)

	for rows.Next() {
		p := &models.Product{}
		err := rows.Scan(
			&p.ProductID, &p.Name, &p.Desc, &p.CategoryID,
			&p.Price, &p.DiscountID, &p.Inventory, &p.Rating,
			&p.RatingNum, &p.UnitSold, &p.SellerID, &p.Created,
			&p.Modified,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// Delete deletes a row which has the specified
// ProductID column value.
func (m *ProductModel) Delete(id int) error {
	stmt := `DELETE FROM Product WHERE ProductID=?`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// GetSellerProducts retrieves every row which
// has the specified SellerID column value.
func (m *ProductModel) GetSellerProducts(sellerID string) ([]*models.Product, error) {
	stmt := `SELECT 
				 ProductID, Name, Description, Price, CategoryID, DiscountID,
				 Inventory, Created, SellerID, Rating, RatingNum, UnitSold
			FROM Product
			WHERE SellerID=?`

	rows, err := m.DB.Query(stmt, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*models.Product{}

	for rows.Next() {
		product := &models.Product{}
		err = rows.Scan(
			&product.ProductID, &product.Name, &product.Desc, &product.Price,
			&product.CategoryID, &product.DiscountID, &product.Inventory,
			&product.Created, &product.SellerID, &product.Rating,
			&product.RatingNum, &product.UnitSold,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

// GetSearchProducts retrieves the ProductID, Name
// Description, and Keyword column values from every
// row in the table.
func (m *ProductModel) GetSearchProducts() ([]*models.Product, error) {
	results, err := m.DB.Query("Select ProductID, Name, Description, Keyword FROM Product")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	products := []*models.Product{}
	for results.Next() {
		product := &models.Product{}
		err = results.Scan(&product.ProductID, &product.Name, &product.Keyword, &product.Desc)
		if err != nil {
			panic(err.Error())
		}
		products = append(products, product)

	}
	return products, nil
}

// GetSearchResults retrieves all the rows which
// have the specified ProductID column values.
// The list of ProductID(s) is returned from the
// search function.
func (m *ProductModel) GetSearchResults(rankedIndex []int) ([]*models.Product, error) {
	if len(rankedIndex) == 0 {
		return nil, nil
	}

	var pIDs string
	for _, ID := range rankedIndex {
		i := strconv.Itoa(ID)
		pIDs = pIDs + "," + i
	}
	pIDs = "(" + strings.TrimLeft(pIDs, ",") + ")"
	stmt := `SELECT ProductID, Name, Description, Price, CategoryID, DiscountID, Inventory, Created, SellerID, Rating, RatingNum, UnitSold
	FROM Product WHERE ProductID IN ` + pIDs

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*models.Product{}
	for rows.Next() {
		product := &models.Product{}
		err = rows.Scan(
			&product.ProductID,
			&product.Name,
			&product.Desc,
			&product.Price,
			&product.CategoryID,
			&product.DiscountID,
			&product.Inventory,
			&product.Created,
			&product.SellerID,
			&product.Rating,
			&product.RatingNum,
			&product.UnitSold,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

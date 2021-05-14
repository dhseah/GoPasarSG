package mysql

import (
	"database/sql"
	"errors"

	"ProjectGoLive/pkg/models"

	"github.com/go-sql-driver/mysql"
)

type UserModel struct {
	DB *sql.DB
}

// Create inserts a new row into the user table.
// Create returns an error if the statement could
// not be executed.
func (m *UserModel) Create(c *models.User) error {
	stmt := `INSERT INTO 
	   user (UserID, Password, FirstName, LastName, Phone, Email, Address, Seller) 
	  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := m.DB.Exec(stmt, c.UserID, c.Password, c.FirstName, c.LastName, c.PhoneNumber, c.Email, c.Address, c.Seller)
	if err != nil {
		if me := err.(*mysql.MySQLError); me.Number == 1062 {
			return models.ErrDuplicateEntry
		}
		return err

	}

	return nil
}

// Get searches the user table database for a row
// that has UserID column value = userID.
// Get returns a pointer to a User populated with
// retrieved values if the search was successful.
// It returns nil and an error otherwise.
func (m *UserModel) Get(userID string) (*models.User, error) {
	stmt := `SELECT * 
			FROM user
			WHERE UserID = ?`

	u := &models.User{}
	row := m.DB.QueryRow(stmt, userID)
	err := row.Scan(
		&u.UserID,
		&u.Password,
		&u.FirstName,
		&u.LastName,
		&u.PhoneNumber,
		&u.Email,
		&u.Address,
		&u.Seller,
		&u.Verified,
		&u.Created,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return u, nil
}

// UpdateVerified updates the Verified column value
// for a row that has UserID column value = userID.
// UpdateVerified returns an error if the statement
// could not be executed.
func (m *UserModel) UpdateVerified(userID string) error {
	stmt := `UPDATE User SET Verified=true WHERE UserID = ?`

	_, err := m.DB.Exec(stmt, userID)
	if err != nil {
		return err
	}

	return nil
}

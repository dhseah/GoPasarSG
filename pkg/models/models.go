package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrNoRowsAffected = errors.New("models: no rows affected")
var ErrDuplicateEntry = errors.New("models: duplicate entry in username/email detected")

var Category = []string{"Frozen Food", "Staples", "Meat and Seafood", "Beverages", "Fruit and Vegetables"}
var SortBy = []string{"Popular", "Highly Rated", "Price (Asc.)", "Price (Desc.)"}

var Discount = []string{"No Discount", "5% discount", "10% discount", "15% discount", "20% discount", "25% discount"}
var DiscMultiplier = []float64{1, .95, .9, .85, .8, .75}

var Status = []string{"Pending", "Accepted", "Cancelled"}

var MapOTP = make(map[string]string)

type Product struct {
	ProductID  int
	Name       string
	Desc       string
	CategoryID int
	Keyword    string
	Price      float64
	DiscountID int
	Inventory  int
	Rating     float64
	RatingNum  int
	UnitSold   int
	SellerID   string
	Created    time.Time
	Modified   time.Time
	Score      int
}

type User struct {
	UserID      string
	Password    string
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
	Address     string
	Seller      bool
	Verified    bool
	Created     time.Time
}

type Orders struct {
	OrderID int
	UserID  string
	Product struct {
		ProductID int
		Name      string
	}
	Qty      int
	SellerID string
	Status   int
	Created  time.Time
	Modified time.Time
}

type CartItem struct {
	UserID  string
	Product struct {
		ProductID  string
		Name       string
		DiscountID int
		Inventory  int
		Price      float64
		SellerID   string
	}
	Qty      int
	Invalid  bool
	Modified time.Time
}

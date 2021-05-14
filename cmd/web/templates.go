package main

import (
	"html/template"
	"math"
	"path/filepath"
	"time"

	"ProjectGoLive/pkg/forms"
	"ProjectGoLive/pkg/models"
)

// templateData holds all the data to be templated.
type templateData struct {
	CurrentYear int
	Flash       string
	Form        *forms.Form
	Error       string

	User     *models.User
	IsSeller bool

	Product    *models.Product
	Products   []*models.Product
	Categories []string
	Discounts  []string
	SortBy     []string

	ShoppingCart []*models.CartItem

	Orders []*models.Orders
	Status []string
}

// newTemplateCache parses all the template files in the
// specified directory and stores it in a cache for
// execution later. newTemplateCache should be called when
// the application is started.
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop through the pages one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl') from the full file path
		// and assign it to the name variable.
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'layout' templates to the
		// template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'partial' templates to the
		// template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page
		// (like 'home.page.tmpl') as the key.
		cache[name] = ts
	}

	// Return the map.
	return cache, nil
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,

	"getDisc":       getDiscount,
	"getCat":        getCategory,
	"getStatus":     getStatus,
	"getFinalPrice": getFinalPrice,
	"getCartTotal":  getCartTotal,
}

// humanDate is a template function that returns
// a formatted string representation of a
// time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// getDiscount is a template function that returns
// the string representation of the dsicount id.
func getDiscount(DiscID int) string {
	for i, v := range models.Discount {
		if i == DiscID {
			return v
		}
	}
	return "No Discount"
}

// getCategory is a template function that returns
// the string representation of the category id.
func getCategory(CatID int) string {
	for i, v := range models.Category {
		if i == CatID {
			return v
		}
	}
	return "No Category"
}

// getStatus is a template function that returns
// the string representation of the status id.
func getStatus(Status int) string {
	for i, v := range models.Status {
		if i == Status {
			return v
		}
	}
	return "No Status"
}

// getFinalPrice is a template function that returns
// the price of a cart item after applying discount.
func getFinalPrice(item models.CartItem) float64 {
	finalPrice := float64(item.Qty) * item.Product.Price * models.DiscMultiplier[item.Product.DiscountID]
	return math.Ceil(finalPrice*100) / 100
}

// getCartTotal is a template function that returns
// the price of all the item in the client's cart.
func getCartTotal(ShoppingCart []*models.CartItem) float64 {
	var CartTotal float64 = 0
	var cart models.CartItem
	for _, v := range ShoppingCart {
		cart = *v
		CartTotal += getFinalPrice(cart)
	}

	return math.Ceil(CartTotal*100) / 100
}

package main

import (
	"ProjectGoLive/pkg/models"
	"errors"
	"net/http"
	"strconv"
)

// ShoppingCart retrieves the shopping cart tagged to
// the userid value in the client's session cookie and
// writes the result to the http response.
func (app *application) ShoppingCart(w http.ResponseWriter, r *http.Request) {
	// a seller does not have a shopping cart
	isSeller := app.isSeller(r)
	if isSeller {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}
	// retrieves userid from the session cookie
	userID := app.session.GetString(r, "userid")

	// query shoppingcart table in the database
	// for the list of item in the user's cart
	cart, err := app.cart.Get(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	// display the list of items
	app.render(w, r, "shoppingcart.page.tmpl", &templateData{
		User:         &models.User{UserID: userID},
		ShoppingCart: cart,
		Discounts:    models.Discount,
	})
}

// AddToCart inserts a new item into the client's
// shopping cart at the database then redirect the
// client back to his/her shopping cart.
func (app *application) AddToCart(w http.ResponseWriter, r *http.Request) {
	// a seller does not have a shopping cart
	isSeller := app.isSeller(r)
	if isSeller {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}

	// retrieve userid from session cookie
	userid := app.session.GetString(r, "userid")

	// retrieve ProductID from url
	// the ProducID should be valid
	productID, err := strconv.Atoi(r.URL.Query().Get("productid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// perform the insert at the database
	err = app.cart.InsertItem(userid, productID)
	if err != nil {
		app.errorLog.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	app.session.Put(r, "flash", "Product successfully added to cart.")

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

// DeleteFromCart deletes an item from the client's
// shoppingcart table in the database then redirects
// the client back to his/her shopping cart.
func (app *application) DeleteFromCart(w http.ResponseWriter, r *http.Request) {
	// a seller does not have a shopping cart
	isSeller := app.isSeller(r)
	if isSeller {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}
	// retrieve userid from session cookie
	userid := app.session.GetString(r, "userid")

	// retrieve ProductID from url
	// the ProducID should be valid
	productid, err := strconv.Atoi(r.URL.Query().Get("productid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// perform the delete at the database
	err = app.cart.DeleteItem(userid, productid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	app.session.Put(r, "flash", "Cart item successfully deleted.")

	http.Redirect(w, r, "/shoppingcart/", http.StatusSeeOther)
}

// UpdateItemQty updates the quantity information
// of an item in the client's shopping cart then
// redirects the client back to his/her shopping cart.
func (app *application) UpdateItemQty(w http.ResponseWriter, r *http.Request) {
	// a seller does not have a shopping cart
	isSeller := app.isSeller(r)
	if isSeller {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}
	// retrieves userid from the session cookie
	userid := app.session.GetString(r, "userid")

	// retrieve productid & quantity from URL
	// both should be valid
	productid, err1 := strconv.Atoi(r.URL.Query().Get("productid"))
	quantity, err2 := strconv.Atoi(r.URL.Query().Get("quantity"))
	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// perform the update at the database
	err := app.cart.Update(quantity, productid, userid)
	if err != nil {
		// if an error did not occur when communicating
		// with the database, but no row was affected
		if errors.Is(err, models.ErrNoRowsAffected) {
			http.Redirect(w, r, "/shoppingcart", http.StatusSeeOther)
			return
		}
		// if an error occured when communicating with
		// the database
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	app.session.Put(r, "flash", "Cart item successfully updated.")

	http.Redirect(w, r, "/shoppingcart/", http.StatusSeeOther)
}

// CheckOut performs a check on the client's shopping
// cart to determine if every item in the cart is legal
// to checkout i.e. quantity less than product inventory.
func (app *application) CheckOut(w http.ResponseWriter, r *http.Request) {
	// retrieve userid from session cookie
	userid := app.session.GetString(r, "userid")

	// check that every item in the cart has legal qty
	// retrieve information from database
	shoppingcart, err := app.cart.CheckOut(userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	// loop through the results to check if every item
	// passes the check
	var pass bool = true
	for _, item := range shoppingcart {
		if item.Invalid {
			pass = false
		}
	}

	// if any item fails the check, it is flagged
	if !pass {
		app.render(w, r, "shoppingcart.page.tmpl", &templateData{
			ShoppingCart: shoppingcart,
		})
	} else {
		// else proceed to create an order for every item
		for _, v := range shoppingcart {
			v.UserID = userid
			err := app.orders.Create(v)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				app.render(w, r, "error.page.tmpl", &templateData{
					Error: http.StatusText(http.StatusInternalServerError),
				})
				return
			}
		}

		// then delete the user's shopping cart
		err = app.cart.DeleteAll(userid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			app.render(w, r, "error.page.tmpl", &templateData{
				Error: http.StatusText(http.StatusInternalServerError),
			})
			return
		}

		app.render(w, r, "success.page.tmpl", &templateData{
			User: &models.User{UserID: userid},
		})
	}
}

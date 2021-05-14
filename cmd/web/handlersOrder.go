package main

import (
	"errors"
	"net/http"
	"strconv"

	"ProjectGoLive/pkg/models"
)

// Orders retrieves all the Orders tagged to the client's
// userid from the database and writes the result to the
// http response.
func (app *application) Orders(w http.ResponseWriter, r *http.Request) {
	// retrieve userid from session cookie
	// & determine if the client is a seller
	userID := app.session.GetString(r, "userid")
	isSeller := app.isSeller(r)

	// retrieve the list of orders tied to userid
	// from the database
	orders, err := app.orders.GetAll(userID, isSeller)

	// reverse order for display
	for i, j := 0, len(orders)-1; i < j; i, j = i+1, j-1 {
		orders[i], orders[j] = orders[j], orders[i]
	}

	if err != nil {
		app.errorLog.Println("Error at Orders..", err)
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	// render different pages with different options
	// depending if client is a seller
	app.render(w, r, "order.page.tmpl", &templateData{
		User:     &models.User{UserID: userID, Seller: isSeller},
		Orders:   orders,
		Status:   models.Status,
		IsSeller: isSeller,
	})
}

// OrderUpdateStatus performs the update at the database
// and redirects the client to the orders page if the
// operation was successful.
func (app *application) OrderUpdateStatus(w http.ResponseWriter, r *http.Request) {
	// only a seller can affect order status
	isSeller := app.isSeller(r)
	if !isSeller {
		app.errorLog.Println(ErrInvalidQueryParams)
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}
	sellerID := app.session.GetString(r, "userid")

	// orderid & status parameters should be legal format
	orderid, err := strconv.Atoi(r.URL.Query().Get("orderid"))
	if err != nil {
		app.errorLog.Println(ErrInvalidQueryParams)
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	order, err := app.orders.Get(orderid)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			w.WriteHeader(http.StatusBadRequest)
			app.render(w, r, "error.page.tmpl", &templateData{
				Error: http.StatusText(http.StatusBadRequest),
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	// verify if the order belongs to the client accessing it
	if sellerID != order.SellerID {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}
	// order status can be affected only if it is pending
	if order.Status != 0 {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}

	// it is safe to assume this basic check is sufficient
	// at this point because of how the sever mux is set up
	status, err := strconv.Atoi(r.URL.Query().Get("status"))
	if err != nil || status >= len(models.Status) || status <= 0 {
		app.errorLog.Println(ErrInvalidQueryParams)
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// perform the update at the database
	err = app.orders.UpdateStatus(orderid, status)
	if err != nil {
		app.errorLog.Println(ErrMySQL, err)
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	http.Redirect(w, r, "/orders", http.StatusSeeOther)
}

// SellerHome retrieves all the products tagged to
// the client's SellerID from the database and writes
// the result to the http response. SellerHome can
// only be accessed by a registered Seller.
func (app *application) SellerHome(w http.ResponseWriter, r *http.Request) {
	// only a seller has a seller home page
	isSeller := app.isSeller(r)
	if !isSeller {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	sellerid := app.session.GetString(r, "userid")

	// retrieve product information from the database
	products, err := app.products.GetSellerProducts(sellerid)
	if err != nil {
		app.errorLog.Println("Error at UserProducts..", err)
		return
	}

	app.render(w, r, "sellerhome.page.tmpl", &templateData{
		User:     &models.User{UserID: sellerid, Seller: isSeller},
		Products: products,
	})
}

// SellerPage retrieves all the products tagged to the
// SellerID provided in the query parameter from the
// database and writes the result to the http response.
// SellerPage can be accessed any visitor to the platform.
func (app *application) SellerPage(w http.ResponseWriter, r *http.Request) {
	// sellerid parameter should be valid
	sellerID := r.URL.Query().Get("sellerid")
	isSeller := app.isSeller(r)
	if sellerID == "" {
		app.errorLog.Println(ErrInvalidQueryParams)
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// retrieve products tagged to the sellerid from the database
	products, err := app.products.GetSellerProducts(sellerID)
	if err != nil {
		app.errorLog.Println("Error at UserProducts..", err)
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	userID := app.session.GetString(r, "userid")
	td := &templateData{}
	if userID != "" {
		td.User = &models.User{UserID: userID, Seller: isSeller}
	}

	td.Products = products
	app.render(w, r, "sellerpage.page.tmpl", td)
}

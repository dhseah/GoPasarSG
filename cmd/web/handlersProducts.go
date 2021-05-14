package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"ProjectGoLive/pkg/forms"
	"ProjectGoLive/pkg/models"
	"ProjectGoLive/pkg/search"
	"ProjectGoLive/pkg/sort"
)

// Home retrieves a list of products from the database,
// filter and sort the results, then writes the result
// to the http response.
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	// retrieve UserID from session
	userID := app.session.GetString(r, "userid")
	isSeller := app.isSeller(r)

	// if UserID exist, pass it to template
	td := &templateData{}
	if userID != "" {
		td.User = &models.User{
			UserID: userID,
			Seller: isSeller,
		}
	}

	// retrieve CategoryID from URL
	category := r.URL.Query().Get("category")
	catID := -1
	for i, v := range models.Category {
		if category == v {
			catID = i
		}
	}

	// retrieve sortby from URL
	// if sortby exist, sort retrieved data according to sortby
	// else sort retrieved data according to default sortby
	sortby := r.URL.Query().Get("sortby")
	sortID := 0
	for i, v := range models.SortBy {
		if sortby == v {
			sortID = i
		}
	}

	// retrieve product information from the database
	products, err := app.products.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	// if CategoryID exist, filter retrieved data according
	// to catagory, else do not filter
	if catID != -1 {
		temp := []*models.Product{}
		for _, v := range products {
			if v.CategoryID == catID {
				temp = append(temp, v)
			}
		}
		products = temp
	}

	// perform sorting on the list with the selected sort logic
	is := sort.NewIntroSort(products, sortID)
	is.IntroSort()

	// prepare the templateData
	td.Products = products
	td.Categories = models.Category
	td.SortBy = models.SortBy

	app.render(w, r, "home.page.tmpl", td)
}

// ProductSearchResults looks up the inverted indexes for
// to compile a list of Products which include the search
// terms in their name, description, and keywords, retrieves
// them from the database then writes the result to the http
// response.
func (app *application) ProductSearchResults(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetString(r, "userid")
	isSeller := app.isSeller(r)

	// retrieve the search term from the URL
	text := r.URL.Query().Get("text")

	//
	intArray, IDScore := app.indSlice.Search(text)

	// retrieve product information from the database
	products, err := app.products.GetSearchResults(intArray)
	if err != nil {
		app.errorLog.Println("Error at ProductSearchResults..", err)
		return
	}

	// sorts the list according to their relevance score
	products = search.RankedProducts(products, IDScore)

	app.render(w, r, "searchresult.page.tmpl", &templateData{
		Products: products,
		User:     &models.User{UserID: userID, Seller: isSeller},
	})
}

// * CRUD operations for Product *

// ProductCreate inserts a new product into the database
// then redirects the client to the product page of the
// newly created product.
func (app *application) ProductCreate(w http.ResponseWriter, r *http.Request) {
	// only a verified seller can create a product
	isSeller := app.isSeller(r)
	verified := app.isVerified(r)
	if !isSeller || !verified {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}
	sellerID := app.session.GetString(r, "userid")

	// parse the submitted form
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// validate the submitted form
	form := forms.New(r.PostForm)
	form.Required("name", "description", "keyword", "price", "category", "discount", "inventory")
	form.MaxLength("name", 150)
	form.MaxLength("description", 500)
	form.MaxLength("price", 15) // should this be removed?
	form.MatchesPattern("price", forms.PriceRX)
	form.MaxLength("inventory", 10) // should this be removed?
	if !form.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "productcreate.page.tmpl", &templateData{
			Form:       form,
			Categories: models.Category,
			Discounts:  models.Discount,
		})
		return
	}

	// convert parameter values into appropritate format
	// for storage in the database
	priceFloat, _ := strconv.ParseFloat(form.Get("price"), 64)
	inventoryIDInt, _ := strconv.Atoi(form.Get("inventory"))

	var catID, discID int
	for i, v := range models.Category {
		if v == form.Get("category") {
			catID = i
		}
	}
	for i, v := range models.Discount {
		if v == form.Get("discount") {
			discID = i
		}
	}

	// perform the insertion at the database
	id, err := app.products.Create(form.Get("name"), form.Get("description"), form.Get("keyword"), sellerID, priceFloat, inventoryIDInt, catID, discID)
	if err != nil {
		app.errorLog.Println(ErrMySQL, err)
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	app.session.Put(r, "flash", "Product successfully added!")

	http.Redirect(w, r, fmt.Sprintf("/product?productid=%v", id), http.StatusSeeOther)
}

// ProductCreateForm writes a html form for creating
// a Product to the http response.
func (app *application) ProductCreateForm(w http.ResponseWriter, r *http.Request) {
	// only a verified seller can create a Product
	isSeller := app.isSeller(r)
	verified := app.isVerified(r)
	if !isSeller || !verified {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}
	sellerID := app.session.GetString(r, "userid")

	td := &templateData{
		Form:       forms.New(nil),
		User:       &models.User{UserID: sellerID},
		Categories: models.Category,
		Discounts:  models.Discount,
	}
	app.render(w, r, "productcreate.page.tmpl", td)
}

// ProductRead retrieves product information from
// the database, then writes the results to the
// http response.
func (app *application) ProductRead(w http.ResponseWriter, r *http.Request) {
	isSeller := app.isSeller(r)
	userID := app.session.GetString(r, "userid")

	// productid parameter value should be valid
	id, err := strconv.Atoi(r.URL.Query().Get("productid"))
	if err != nil || id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// retrieve the product information from the database
	product, err := app.products.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			w.WriteHeader(http.StatusNotFound)
			app.render(w, r, "error.page.tmpl", &templateData{
				Error: http.StatusText(http.StatusNotFound),
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	// different options will be presented depending
	// if the client is a seller
	app.render(w, r, "product.page.tmpl", &templateData{
		User:     &models.User{UserID: userID, Seller: isSeller},
		Product:  product,
		IsSeller: isSeller,
	})
}

// ProductUpdate updates the product information for
// a product at the database,  then redirects the client
// to the product page of the edited product.
func (app *application) ProductUpdate(w http.ResponseWriter, r *http.Request) {
	// only a verified seller can update product information
	isSeller := app.isSeller(r)
	verified := app.isVerified(r)
	if !isSeller || !verified {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}
	sellerID := app.session.GetString(r, "userid")

	// productid parameter value should be valid
	id, err := strconv.Atoi(r.URL.Query().Get("productid"))
	if err != nil {
		app.errorLog.Println(ErrInvalidQueryParams)
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// retrieve product from database to repopulate form
	// with original data
	p, err := app.products.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// verify if the product belongs to the client
	if sellerID != p.SellerID {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.errorLog.Println(ErrInvalidForm)
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "description", "keyword", "price", "category", "discount", "inventory")
	form.MaxLength("name", 150)
	form.MaxLength("description", 500)
	form.MaxLength("price", 15)
	form.MatchesPattern("price", forms.PriceRX)
	form.MaxLength("inventory", 10)
	if !form.Valid() {
		app.render(w, r, "update.page.tmpl", &templateData{
			Form:       form,
			User:       &models.User{UserID: sellerID},
			Categories: models.Category,
			Discounts:  models.Discount,
			Product:    p,
		})
		return
	}

	// convert parameter values to their appropriate
	// format for storage in the database
	priceFloat, _ := strconv.ParseFloat(form.Get("price"), 64)
	var catID, discID int
	for i, v := range models.Category {
		if v == form.Get("category") {
			catID = i
		}
	}
	for i, v := range models.Discount {
		if v == form.Get("discount") {
			discID = i
		}
	}
	inventory, _ := strconv.Atoi(form.Get("inventory"))

	// perform the update at the database
	err = app.products.Update(form.Get("name"), form.Get("description"), form.Get("keyword"), priceFloat, inventory, id, catID, discID)
	if err != nil {
		app.errorLog.Println(ErrMySQL, err)
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	app.session.Put(r, "flash", "Product successfully updated!")

	http.Redirect(w, r, fmt.Sprintf("/product?productid=%v", id), http.StatusSeeOther)
}

// ProductUpdateForm writes a html form for updating
// product information for a product in the database
// to the http response.
func (app *application) ProductUpdateForm(w http.ResponseWriter, r *http.Request) {
	// only a verified seller can update product information
	isSeller := app.isSeller(r)
	verified := app.isVerified(r)
	if !isSeller || !verified {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}
	sellerID := app.session.GetString(r, "userid")

	id, err := strconv.Atoi(r.URL.Query().Get("productid"))
	if err != nil {
		app.errorLog.Println(ErrInvalidQueryParams)
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// retrieve existing information from the database
	product, err := app.products.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			w.WriteHeader(http.StatusNotFound)
			app.render(w, r, "error.page.tmpl", &templateData{
				Error: http.StatusText(http.StatusNotFound),
			})
			return
		}
		app.errorLog.Println(ErrMySQL, err)
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	// verify if the product belongs to the client
	if sellerID != product.SellerID {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}

	app.render(w, r, "update.page.tmpl", &templateData{
		Form:       forms.New(nil),
		User:       &models.User{UserID: sellerID, Seller: isSeller},
		Categories: models.Category,
		Discounts:  models.Discount,
		Product:    product,
	})
}

// ProductDelete deletes a product from the database
// then redirects the client to the sellerhome page.
func (app *application) ProductDelete(w http.ResponseWriter, r *http.Request) {
	// only a verified seller can delete a product
	isSeller := app.isSeller(r)
	verified := app.isVerified(r)
	if !isSeller || !verified {
		app.infoLog.Println(ErrInvalidCredentials)
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}

	sellerID := app.session.GetString(r, "userid")

	// productid parameter value should be valid
	id, err := strconv.Atoi(r.URL.Query().Get("productid"))
	if err != nil || id < 1 {
		app.errorLog.Println(ErrInvalidQueryParams)
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	p, err := app.products.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// verify if the product belongs to the client
	if sellerID != p.SellerID {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}

	// perform the delete at the database
	err = app.products.Delete(id)
	if err != nil {
		app.errorLog.Println(ErrMySQL, err)
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	app.session.Put(r, "flash", "Product Successfully deleted.")

	http.Redirect(w, r, "/sellerhome", http.StatusSeeOther)
}

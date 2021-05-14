package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	// standard middleware chain
	stdstack := alice.New(app.recoverPanic, app.logRequest, secureHeaders, app.session.Enable, app.authenticate)
	// standard middleware chain plus user authentication
	authpipe := stdstack.Append(app.requireAuthentication)

	// BUYER
	r.Handle("/home", stdstack.ThenFunc(app.Home))

	// SELLER
	r.Handle("/seller", stdstack.ThenFunc(app.SellerPage)).Methods("GET").Queries("sellerid", "{sellerid}")
	r.Handle("/sellerhome", authpipe.ThenFunc(app.SellerHome))

	// PRODUCT
	r.Handle("/product/create", authpipe.ThenFunc(app.ProductCreateForm)).Methods("GET")
	r.Handle("/product/create", authpipe.ThenFunc(app.ProductCreate)).Methods("POST")
	r.Handle("/product", stdstack.ThenFunc(app.ProductRead)).Methods("GET").Queries("productid", "{productid}")
	r.Handle("/product/update", authpipe.ThenFunc(app.ProductUpdateForm)).Methods("GET").Queries("productid", "{productid}")
	r.Handle("/product/update", authpipe.ThenFunc(app.ProductUpdate)).Methods("POST").Queries("productid", "{productid}")
	r.Handle("/product/delete", authpipe.ThenFunc(app.ProductDelete)).Methods("GET").Queries("productid", "{productid}")
	r.Handle("/product/search", stdstack.ThenFunc(app.ProductSearchResults)).Methods("GET").Queries("text", "{text}")

	// LOG-IN, LOG-OUT
	r.Handle("/login", stdstack.ThenFunc(app.LogInForm)).Methods("GET")
	r.Handle("/login", stdstack.ThenFunc(app.LogIn)).Methods("POST")
	r.Handle("/logout", stdstack.ThenFunc(app.LogOut))
	r.Handle("/signup", stdstack.ThenFunc(app.SignUpForm)).Methods("GET")
	r.Handle("/signup", stdstack.ThenFunc(app.SignUp)).Methods("POST")

	// VERIFICATION
	r.Handle("/verifyuser", stdstack.ThenFunc(app.VerifyUserForm)).Methods("GET")
	r.Handle("/verifyuser", stdstack.ThenFunc(app.VerifyUser)).Methods("POST")
	r.Handle("/verifyuser/resendotp", stdstack.ThenFunc(app.resendOTP)).Methods("POST")

	// SHOPPING CART
	r.Handle("/shoppingcart/", authpipe.ThenFunc(app.ShoppingCart)).Methods("GET")
	r.Handle("/shoppingcart/add", authpipe.ThenFunc(app.AddToCart)).Methods("GET").Queries("productid", "{productid}")
	r.Handle("/shoppingcart/delete", authpipe.ThenFunc(app.DeleteFromCart)).Methods("POST").Queries("productid", "{productid}")
	r.Handle("/shoppingcart/update", authpipe.ThenFunc(app.UpdateItemQty)).Methods("GET").Queries("productid", "{productid}", "quantity", "{quantity}")

	// CHECKOUT
	r.Handle("/checkout", authpipe.ThenFunc(app.CheckOut)).Methods("POST")

	// ORDERS
	r.Handle("/orders", authpipe.ThenFunc(app.Orders)).Methods("GET")
	r.Handle("/order", authpipe.ThenFunc(app.OrderUpdateStatus)).Methods("POST").Queries("orderid", "{orderid}", "status", "{status}")

	// FAVICON
	r.Handle("/favicon.ico", http.NotFoundHandler())

	// FILESERVER for style sheets etc.
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/"))))

	return r
}

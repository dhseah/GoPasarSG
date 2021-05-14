package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// serverError prints the function stack to the error log
// and closes the connection with the client.
// serverError is only called when an unexpected error that
// was not caught by the application occur.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// addDefaultData inserts template data shared across handlers
// into the template. addDefaultData inserts the data into an
// existing templateData struct if it exist. Otherwisem it
// initialises a new templateData struct.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	// Add the flash message to the template data, if one exists.
	td.Flash = app.session.PopString(r, "flash")
	return td
}

// render executes the template with the specified name &
// provided template data and writes it to the http response.
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Execute the template set, including the dynamic data
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the buffer to the http.ResponseWriter
	buf.WriteTo(w)
}

// hashPassword hashes a cleartext password using
// the bcrypt password-hasing function
func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash)
}

// isAuthenticated checks the request's context
// for a value mapped to contextKeyAuthUser to
// determine if the client is logged in.
func (app *application) isAuthenticated(r *http.Request) bool {
	_, ok := r.Context().Value(contextKeyAuthUser).(AuthUser)
	return ok
}

// isAuthenticated checks the request's context
// for a value mapped to contextKeyAuthUser to
// determine if the client is a registered seller.
func (app *application) isSeller(r *http.Request) bool {
	authUser, ok := r.Context().Value(contextKeyAuthUser).(AuthUser)
	if !ok {
		return false
	}
	return authUser.Seller
}

// isAuthenticated checks the request's context
// for a value mapped to contextKeyAuthUser to
// determine if the client is a verified user.
func (app *application) isVerified(r *http.Request) bool {
	authUser, ok := r.Context().Value(contextKeyAuthUser).(AuthUser)
	if !ok {
		return false
	}
	return authUser.Verified
}

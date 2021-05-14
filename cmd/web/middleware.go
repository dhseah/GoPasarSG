package main

import (
	"ProjectGoLive/pkg/models"
	"context"
	"errors"
	"fmt"
	"net/http"
)

// recoverPanic recovers from error unwinding up
// the function stack. recoverPanic is placed at
// the top of the middleware chain so it catches
// unexpected errors.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// logRequest logs the ip address, protocol ver.
// method & requestURL of the incoming request
// before calling the next handler in the chain.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// secureHeaders applies HTTP security settings
// against XSS attacks.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

// requireAuthentication checks if the client is
// a verified user who is logged in.
// requireAuthentication redirects the client if
// he/she is not logged in and/or not verified.
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the client is not authenticated, redirect
		if loggedin := app.isAuthenticated(r); !loggedin {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// If the client is not a verified user, redirect
		if verified := app.isVerified(r); !verified {
			http.Redirect(w, r, "/verifyuser", http.StatusSeeOther)
			return
		}

		/* This section should be reviewed */
		// Else, set the "Cache-Control: no-store" header so pages
		// which require authentication are not stored in cache
		w.Header().Add("Cache-Control", "no-store")

		// and call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// authenticate checks the if client making the request
// has a valid userID value in his/her request coookie
// and passes this information down the middleware chain
// via the request's context.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if a userID value does not exists in the
		// session cookie, call the next handler
		exists := app.session.Exists(r, "userid")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		// else, check that the userID value exist in
		// the database. If no matching record is found
		// remove the invalid userID value from the
		// session cookie and call the next handler.
		user, err := app.users.Get(app.session.GetString(r, "userid"))
		if errors.Is(err, models.ErrNoRecord) {
			app.session.Remove(r, "authUser")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.errorLog.Println(ErrMySQL, err)
			w.WriteHeader(http.StatusInternalServerError)
			app.render(w, r, "error.page.tmpl", &templateData{
				Error: http.StatusText(http.StatusInternalServerError),
			})
			return
		}

		// else, the request is coming from an authenticated user
		// Create a new copy of the request, with an AuthUser value
		// added to the request context to indicate this, and call
		// the next handler with this new copy of the request
		au := AuthUser{user.UserID, user.Seller, user.Verified}
		ctx := context.WithValue(r.Context(), contextKeyAuthUser, au)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

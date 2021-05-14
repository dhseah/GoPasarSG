package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"ProjectGoLive/pkg/forms"
	"ProjectGoLive/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

// LogIn verifies if the form values submitted by the client
// points to a user in the database. LogIn creates a session
// for the client if the credentials are accurate. Otherwise,
// it writes a error message to the http response. LogIn
// redirects the client if he/she is already logged in.
func (app *application) LogIn(w http.ResponseWriter, r *http.Request) {
	// if user is already loggedin, redirect to home
	loggedin := app.isAuthenticated(r)
	verified := app.isVerified(r)
	isSeller := app.isSeller(r)
	if loggedin && verified {
		if isSeller {
			http.Redirect(w, r, "/sellerhome", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if loggedin && !verified {
		http.Redirect(w, r, "/verifyuser", http.StatusSeeOther)
		return
	}

	// if client submitted a form to log in
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{Error: http.StatusText(http.StatusBadRequest)})
		return
	}

	form := forms.New(r.PostForm)
	username := form.Get("userid")
	password := form.Get("password")

	// verify if the information submitted is valid
	var errUser, errPW error
	user, errUser := app.users.Get(username)
	if errUser != nil {
		errPW = bcrypt.CompareHashAndPassword([]byte{}, []byte(password))
	} else {
		errPW = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	}
	// if the submitted username returns no match in the database
	// or if the password is incorrect, reject the log-in attempt
	if errUser != nil || errPW != nil {
		w.WriteHeader(http.StatusUnauthorized)
		form.Errors.Add("unauthorized", "Invalid UserID and/or Password provided")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	}

	// else the log-in was successful
	// create a session for the user & redirect
	app.session.Put(r, "userid", username)

	if user.Seller {
		http.Redirect(w, r, "/sellerhome", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

// LogInForm writes a html form for loggin in to the http response.
// LogInForm redirects the client if he/she is already logged in.
func (app *application) LogInForm(w http.ResponseWriter, r *http.Request) {
	// if user is already loggedin, redirect to home
	loggedin := app.isAuthenticated(r)
	verified := app.isVerified(r)
	if loggedin && verified {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if loggedin && !verified {
		http.Redirect(w, r, "/verifyuser", http.StatusSeeOther)
		return
	}

	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

// Logout logs out an authenticated user by expiring the
// session tagged to the user.
func (app *application) LogOut(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userid")
	c := &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Now(),
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// SignUp inserts a new user into the database, then
// creates a session for the client and redirects him
// /her to the user verification page. SignUp redirects
// the client to the home page if he/she is already
// logged in.
func (app *application) SignUp(w http.ResponseWriter, r *http.Request) {
	// if user is already logged in, redirect to home
	if loggedin := app.isAuthenticated(r); loggedin {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	form := forms.New(r.PostForm)
	form.Required("userid", "password", "firstname", "lastname", "phone", "email", "address")
	form.MaxLength("userid", 30)
	form.MaxLength("password", 64)
	form.MinLength("password", 8)
	form.MaxLength("firstname", 20)
	form.MaxLength("lastname", 20)
	form.MaxLength("phone", 16)
	form.MinLength("phone", 8)
	form.MaxLength("email", 150)
	form.MatchesPattern("email", forms.EmailRX)
	form.MaxLength("address", 150)
	form.PermittedValues("signupas", "user", "seller")
	if !form.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	newUser := &models.User{
		UserID:      r.PostFormValue("userid"),
		Password:    hashPassword(r.PostFormValue("password")),
		FirstName:   r.PostFormValue("firstname"),
		LastName:    r.PostFormValue("lastname"),
		PhoneNumber: r.PostFormValue("phone"),
		Email:       r.PostFormValue("email"),
		Address:     r.PostFormValue("address"),
		Seller:      (r.PostFormValue("signupas") == "seller"),
	}

	// perform the insert into the database
	err = app.users.Create(newUser)
	if err != nil {
		if !errors.Is(err, models.ErrDuplicateEntry) {
			w.WriteHeader(http.StatusInternalServerError)
			app.render(w, r, "error.page.tmpl", &templateData{
				Error: http.StatusText(http.StatusInternalServerError),
			})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: "UserID/Email already exist",
		})
		return
	}

	app.session.Put(r, "userid", newUser.UserID)

	// spawn a goroutine to send an OTP to
	// verify email of user
	go app.sendOTP(r, newUser.Email)

	http.Redirect(w, r, "/verifyuser", http.StatusSeeOther)
}

// SignUpForm writes the html form for signing up to the http response.
// SignUpForm redirects the client if he/she is already logged in.
func (app *application) SignUpForm(w http.ResponseWriter, r *http.Request) {
	// if user is already loggedin & verified, redirect to home
	if loggedin := app.isAuthenticated(r); loggedin {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	fmt.Println("continue execution")
	app.render(w, r, "signup.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

// * USER VERIFICATION *

// VerifyUserForm writes the html form for submitting
// an OTP to the http response.
func (app *application) VerifyUserForm(w http.ResponseWriter, r *http.Request) {
	if isVerified := app.isVerified(r); isVerified {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	// If IsVerified is false, renders a form to user to enter OTP
	// Also a button for user to send a new OTP if lost first OTP email.
	app.render(w, r, "verifyuser.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

// VerifyUser updates the user's status to verified
// if the OTP submitted by the user matches the OTP
// mapped to the user's userID and
func (app *application) VerifyUser(w http.ResponseWriter, r *http.Request) {
	// redirect the client if he/she is already verified
	if isVerified := app.isVerified(r); isVerified {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	userid := app.session.GetString(r, "userid")

	//parse form to obtain form OTP value
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// check if the OTP is a 6-digit string
	form := forms.New(r.PostForm)
	form.Required("OTP")
	form.MaxLength("OTP", 6)
	form.MinLength("OTP", 6)
	if !form.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "verifyuser.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	// check if the OTP exist in the OTP map
	OTP := form.Get("OTP")
	if models.MapOTP[userid] == OTP {
		err = app.users.UpdateVerified(userid)
		if err != nil {
			app.errorLog.Println(ErrMySQL, err)
			w.WriteHeader(http.StatusInternalServerError)
			app.render(w, r, "error.page.tmpl", &templateData{
				Error: http.StatusText(http.StatusInternalServerError),
			})
			return
		}
		delete(models.MapOTP, userid)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusUnauthorized),
		})
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// resendOTP sends a new OTP to the client if the
// conditions are met.
func (app *application) resendOTP(w http.ResponseWriter, r *http.Request) {
	userid := app.session.GetString(r, "userid")

	// if there is still an OTP tagged to the client's
	// usedid, the client is reminded to check his/her email
	if _, exist := models.MapOTP[userid]; exist {
		w.WriteHeader(http.StatusBadRequest)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: "Please check your email for the OTP.",
		})
		return
	}

	// verify if the request comes from a registered user
	user, err := app.users.Get(userid)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			w.WriteHeader(http.StatusUnauthorized)
			app.render(w, r, "error.page.tmpl", &templateData{
				Error: http.StatusText(http.StatusUnauthorized),
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		app.render(w, r, "error.page.tmpl", &templateData{
			Error: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	// spawn a goroutine to send the email so
	// redirect can happen instantly
	go app.sendOTP(r, user.Email)

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

// sendOTP sends a one-time password to the user
// via smtp to verify the email provided by the
// user during sign-up.
func (app *application) sendOTP(r *http.Request, userEmail string) {

	userid := app.session.GetString(r, "userid")
	OTP := app.getOTPToken()

	models.MapOTP[userid] = OTP

	go func() {
		time.Sleep(time.Minute * 3)
		delete(models.MapOTP, userid)
	}()

	from := "gopasarsg@gmail.com"
	password := "GoPasar!@#123"

	// recipient's email address
	to := []string{userEmail}

	// smtp server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("To:" + userEmail + "\r\n" +
		"From: gopasarsg@gmail.com\r\n" +
		"Subject: Thank you for registering with GoPasarSG!\r\n" +
		"\r\n" + "Thank you for registering with GoPasar SG.\r\n" +
		"Your verification code is " + OTP + ".\r\n" + "\r\n" +
		"This OTP will expire in 3 minutes. \r\n" +
		"Please note that this is a self-generating email. Do not reply to this email.\r\n" +
		"Regards,\r\nGoPasar SG Team\r\n")

	// smtp authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}

	app.infoLog.Println("OTP-verification email sent!")
}

// getOTPToken generates a 6-digit one-time password
// credit https://github.com/tilaklodha/google-authenticator
func (app *application) getOTPToken() string {
	// secret has to be REMOVED
	secret := "goPasarSingapore"
	var interval int64 = time.Now().Unix() / 5
	//Converts secret to base32 Encoding. Base32 encoding desires a 32-character
	//subset of the twenty-six letters A–Z and ten digits 0–9
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		app.errorLog.Println(err)
		return ""
	}
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(interval))

	//Signing the value using HMAC-SHA1 Algorithm
	hash := hmac.New(sha1.New, key)
	hash.Write(bs)
	h := hash.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA) and we need 4 bytes.
	o := (h[19] & 15)

	var header uint32
	//Get 32 bit chunk from hash starting at the o
	r := bytes.NewReader(h[o : o+4])
	err = binary.Read(r, binary.BigEndian, &header)

	if err != nil {
		app.errorLog.Println(err)
		return ""
	}

	//Ignore most significant bits as per RFC 4226.
	//Takes division from one million to generate a remainder less than < 7 digits
	h12 := (int(header) & 0x7fffffff) % 1000000

	//Converts number as a string
	otp := strconv.Itoa(int(h12))
	return prefix0(otp)
}

// prefix0 appends padding zeros to the front
// of the OTP if it is less than 6-digits long.
func prefix0(otp string) string {
	if len(otp) == 6 {
		return otp
	}
	for i := (6 - len(otp)); i > 0; i-- {
		otp = "0" + otp
	}
	return otp
}

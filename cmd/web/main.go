package main

import (
	"ProjectGoLive/pkg/models/mysql"
	"ProjectGoLive/pkg/search"
	"crypto/tls"
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/joho/godotenv"
)

type contextKey string

const contextKeyAuthUser = contextKey("authUser")

type AuthUser struct {
	UserID   string
	Seller   bool
	Verified bool
}

type application struct {
	// loggers
	infoLog  *log.Logger
	alertLog *log.Logger
	errorLog *log.Logger
	panicLog *log.Logger
	fatalLog *log.Logger

	// template engine
	templateCache map[string]*template.Template

	// session manager
	session  *sessions.Session
	indSlice *search.IndexSlice //reverse index to hold word index of

	// database connection
	users    *mysql.UserModel
	products *mysql.ProductModel
	cart     *mysql.CartModel
	orders   *mysql.OrderModel
}

var (
	dbUsername string
	dbPassword string
	dbName     string
	secretKey  string
	port       string
	host       string
	infoLog    *log.Logger
	alertLog   *log.Logger
	errorLog   *log.Logger
	panicLog   *log.Logger
	fatalLog   *log.Logger
)

//getting all the information from .env file
func init() {

	secretKey = goDotEnvVariable("secretkey")
	dbUsername = goDotEnvVariable("dbUsername")
	dbPassword = goDotEnvVariable("dbPassword")
	dbName = goDotEnvVariable("dbName")
	port = goDotEnvVariable("port")
	host = goDotEnvVariable("host")

}

func main() {

	infoLogFile, err := os.OpenFile("./log/infolog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open info log file:", err)
	}

	alertLogFile, err := os.OpenFile("./log/alertlog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open alert log file:", err)
	}

	errorLogFile, err := os.OpenFile("./log/errorlog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open error log file:", err)
	}

	infoLog = log.New(io.MultiWriter(infoLogFile, os.Stdout), "INFO\t", log.Ldate|log.Ltime)
	alertLog = log.New(io.MultiWriter(alertLogFile, os.Stdout), "ALERT\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(io.MultiWriter(errorLogFile, os.Stderr), "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	panicLog = log.New(io.MultiWriter(errorLogFile, os.Stderr), "PANIC\t", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLog = log.New(io.MultiWriter(errorLogFile, os.Stderr), "FATAL\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dbUsername + ":" + dbPassword + "@/" + dbName + "?parseTime=true")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(secretKey))
	session.Lifetime = 12 * time.Hour

	app := &application{
		infoLog:       infoLog,
		alertLog:      alertLog,
		errorLog:      errorLog,
		panicLog:      panicLog,
		fatalLog:      fatalLog,
		templateCache: templateCache,
		session:       session,
		users:         &mysql.UserModel{DB: db},
		products:      &mysql.ProductModel{DB: db},
		cart:          &mysql.CartModel{DB: db},
		orders:        &mysql.OrderModel{DB: db},
	}

	app.indSlice = app.makeSearchIndexMap()

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         host + ":" + port,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go app.backgroundCleaner()
	go app.backgroundHelper()

	// update to HTTPS eventually
	infoLog.Printf("Starting server on %s", port)

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

//backgroundCleaner is a go routine to perform constant
//background clean up for shopping cart and unverified user.
func (app *application) backgroundCleaner() {

	for range time.Tick(time.Hour * 24) {

		err := app.cart.ShoppingCartCleanUp()
		if err != nil {
			app.errorLog.Println(err)
		}
		err = app.users.VerifiedUserCleanUp()
		if err != nil {
			app.errorLog.Println(err)
		}

	}

}

//backgroundHelper is a go routine to perform constant background update for search bar inverted index.
func (app *application) backgroundHelper() *search.IndexSlice {

	for range time.Tick(time.Second * 30) {

		app.indSlice = app.makeSearchIndexMap()
		infoLog.Printf("Inverted map is refreshed")

	}

	return nil
}

//function to create inverted index map for search words
func (app *application) makeSearchIndexMap() *search.IndexSlice {
	indSlice := search.IndexSlice{map[string][]int{}, map[string][]int{}, map[string][]int{}}
	productForSearch, err := app.products.GetSearchProducts()
	if err != nil {
		app.errorLog.Println(err)
		return nil
	}
	indSlice.Add(productForSearch)
	return &indSlice
}

// use godot package to load/read the .env file and return the value of the key
func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)

}

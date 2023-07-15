package main

import (
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/v2"
	"github.com/dapetoo/snippetbox/pkg/models/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *scs.SessionManager
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	//This will be used when running the app from the command line
	addr := flag.String("addr", ":4000", "HTTP network address")

	dsn := flag.String("dsn", "peter:password@/snippetbox?parseTime=true", "MySQL data source name")
	//secret := flag.String("secret", "", "Secret Key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	//Initialize a new template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new session manager and configure the session lifetime.
	session := scs.New()
	session.Lifetime = 12 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	//Initialize a new instance of application containing the dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		session:  session,
		snippets: &mysql.SnippetModel{
			DB: db,
		},
		templateCache: templateCache,
	}

	infoLog.Println("Database Connected Successfully....")

	//Setting DB
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	infoLog.Printf("Starting the webserver on port %v", *addr)
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	//err = srv.ListenAndServe()
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

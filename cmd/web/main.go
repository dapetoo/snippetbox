package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	//This will be used when running the app from the command line
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	//Logging to a file
	//f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer f.Close()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Initialize a new instance of application containing the dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	//FileServer to serve static files
	fileServer := http.FileServer(http.Dir("./ui/static"))

	//Register the File server with mux.Handle
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	infoLog.Printf("Starting the webserver on port %v", *addr)
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

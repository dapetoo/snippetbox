package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	//This will be used when running the app from the command line
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	//
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Llongfile)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

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

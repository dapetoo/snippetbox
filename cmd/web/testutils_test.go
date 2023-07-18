package main

import (
	"io"
	"log"
	"net/http/httptest"
	"testing"
)

// Create a new test application to return instance of the application struct
func newTestApplication(t *testing.T) *application {
	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}
}

// Define a custom testServer type which will anonymously embeds a httptest.Server instance
type testServer struct {
	*httptest.Server
}

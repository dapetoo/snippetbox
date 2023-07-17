package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	//Initialize a new httptest.ResponseRecorder
	rr := httptest.NewRecorder()

	//Initialize a new dummy http.Request
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	//Call the ping handler function, passing in the ResponseRecorder and the Request
	ping(rr, r)

	//Call the Result() on the ResponseRecorder to get the http.Response generated
	rs := rr.Result()

	//Examine the response to check that the status code written by the Ping Handler was 200
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	//Check the Response body written by the ping handler equals "OK"
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q got %q", "OK", body)
	}
}

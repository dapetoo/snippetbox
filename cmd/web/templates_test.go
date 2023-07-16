package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	//Initialize a new time.time object and pass it to the humanDate function
	tm := time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC)
	hd := humanDate(tm)

	//Check that the output match expected
	if hd != "17 Dec 2020 at 10:00" {
		t.Errorf("want %q; got %q", "17 Dec 2020 at 10:00", hd)
	}
}

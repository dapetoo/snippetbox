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

func TestHumanDate2(t *testing.T) {
	//Create a slice of anonymous structs containing the test case name input to HumanDate()
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Dec 2020 at 10:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Dec 2020 at 09:00",
		},
	}

	//Loop over the test cases
	for _, tt := range tests {
		//t.Run() to run a subtests for each case. The first parameter is the name of the test and the second parameter
		//is the anonymous function containing the actual test
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			if hd != tt.want {
				t.Errorf("want %q; got %q", tt.want, hd)
			}
		})
	}
}

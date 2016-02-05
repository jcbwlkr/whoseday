package main

import (
	"testing"
	"time"
)

func TestPicForDay(t *testing.T) {
	tests := []struct {
		day time.Time
		pic string
	}{
		{time.Date(2016, time.January, 8, 20, 0, 0, 0, time.Local), "kell.jpeg"},
		{time.Date(2016, time.January, 9, 20, 0, 0, 0, time.Local), "carter.jpeg"},
	}

	for i, test := range tests {
		actual := picForDay(test.day)
		if actual != test.pic {
			t.Errorf("Test %d: picForDay(%v) = %v, expected %v", i, test.day, actual, test.pic)
		}
	}
}

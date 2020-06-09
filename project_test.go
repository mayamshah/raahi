package main

import (
	"testing"
)

func TestStraightLine(t *testing.T) {
	tests := []struct {
		address string
		distance string
		error_fix float64
	}{
		{`12543 Palmtag Drive, Saratoga, CA`, `1`, .8},

	}
		for _, test := range tests {
		success, percent_error := execute(test.address, test.distance, straight_line, test.error_fix) 
		if !success || percent_error >= 10.0 {
			t.Errorf("Address : %q\n Distance: %q\n Route_Option: straight_line\n ErrorFix: %f\n Success: %t\n Percent Error: %f\n", test.address, test.distance, test.error_fix, success, percent_error)
		}
	}
} 

func TestSquareRoute(t *testing.T) {
	tests := []struct {
		address string
		distance string
		error_fix float64
	}{
		{`12543 Palmtag Drive, Saratoga, CA`, `1`, .8},

	}
		for _, test := range tests {
		success, percent_error := execute(test.address, test.distance, square_route, test.error_fix) 
		if !success || percent_error >= 10.0 {
			t.Errorf("Address : %q\n Distance: %q\n Route_Option: square_route \n ErrorFix: %f\n Success: %t\n Percent Error: %f\n", test.address, test.distance, test.error_fix, success, percent_error)
		}
	}
} 
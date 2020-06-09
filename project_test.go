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

////////////// DO NOT RUN THESE TESTS - THEY HAVE ALREADY BEEN RUN /////////////////
////////////// THE OUTPUT IS LOCATED IN THE REPO //////////////////////////////////
		
		// //suburban
		// {`12543 Palmtag Drive, Saratoga, CA`, `1`, 1},
		// {`327 Westwood Ln, Stockton, CA 95207`,`1`, 1},
		// {`1511 Brook Valley Ct, Dallas, TX 75232`,`1`, 1},
		// {`627 N Jackson St, Arlington, VA 22201`,`1`, 1},

		// //city
		// {`1050 Fell St, San Francisco, CA 94117`, `1`, 1},
		// {`425 5th Ave, New York, NY 10016`, `1`, 1},
		// {`1037 Chestnut St, Philadelphia, PA 1910`, `1`, 1},
		// {`4551 SW 5th St, Coral Gables, FL 33134`,`1`, 1},


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

////////////// DO NOT RUN THESE TESTS - THEY HAVE ALREADY BEEN RUN /////////////////
////////////// THE OUTPUT IS LOCATED IN THE REPO //////////////////////////////////

		// //suburban
		// {`12543 Palmtag Drive, Saratoga, CA`, `1`, 1},
		// {`327 Westwood Ln, Stockton, CA 95207`,`1`, 1},
		// {`1511 Brook Valley Ct, Dallas, TX 75232`,`1`, 1},
		// {`627 N Jackson St, Arlington, VA 22201`,`1`, 1},

		// //city
		// {`1050 Fell St, San Francisco, CA 94117`, `1`, 1},
		// {`425 5th Ave, New York, NY 10016`, `1`, 1},
		// {`1037 Chestnut St, Philadelphia, PA 1910`, `1`, 1},
		// {`4551 SW 5th St, Coral Gables, FL 33134`,`1`, 1},

	}
		for _, test := range tests {
		success, percent_error := execute(test.address, test.distance, square_route, test.error_fix) 
		if !success || percent_error >= 10.0 {
			t.Errorf("Address : %q\n Distance: %q\n Route_Option: square_route \n ErrorFix: %f\n Success: %t\n Percent Error: %f\n", test.address, test.distance, test.error_fix, success, percent_error)
		}
	}
} 
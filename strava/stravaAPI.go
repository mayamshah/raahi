package strava 

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"github.com/strava/go.strava"
	"sync"
	"github.com/twpayne/go-polyline"
)

func polylineToPath(polyline strava.Polyline) []float64 {
	coords := polyline.Decode()
	path := []float64{}
	for _, coord := range coords {
		path = append(path, coord[0])
		path = append(path, coord[1])
	}

	return path
}

func pathToPolyline(path []float64) string {
	var coords [][]float64
	for i := 0; i < (len(path) / 2); i+=2 {
		temp := []float64{path[i], path[i + 1]}
		//fmt.Println(temp)
		coords = append(coords, temp)
	}
	
	poly := fmt.Sprintf("%s", polyline.EncodeCoords(coords))

	return poly
}

func ExecuteStravaRequest(input string, distance_string string, radius_string string, error_fix float64, accessToken string) *FullResponse {
	//check to see if radius is a proper number
	radius, err := strconv.ParseFloat(radius_string, 64)
	if err != nil {
		return getErrorResponse("Not a valid radius")
	}

	//form url from address
	url := address_to_api_call(input)

	//get response from google api server
	response, error := api_request(url)
	if (error != ``) {
		return getErrorResponse(error)
	}

	//check to see if address exists
	if !check_responseGeocode(response) {
		return getErrorResponse("Address doesn't exist")
	}

	//get the latitude and longitude
	lat, lng := extract_coordinates(response)
	origin := NewPoint(lat, lng)

	//get the points for teh request
	top_right, error := get_point(origin, radius/2, 45, false)
	if (error != ``) {
		return getErrorResponse(error)
	}
	bottom_left, error := get_point(origin, radius/2, 45+180, false)
	if (error != ``) {
		return getErrorResponse(error)
	}
	client := strava.NewClient(accessToken)
	SegmentCall := strava.NewSegmentsService(client).Explore(bottom_left.lat, bottom_left.lng, top_right.lat, top_right.lng)
	SegmentCall.ActivityType("running")
	SegmentCall.MinimumCategory(1)
	SegmentCall.MaximumCategory(100)

	responses, err := SegmentCall.Do()
	if err != nil {
		return getErrorResponse(fmt.Sprintf(`%s`, err))
	}

	//fmt.Println(len(responses))

	if len(responses) <= 0 {
		return getErrorResponse(`No routes found`)
	}

	distance, err := strconv.ParseFloat(distance_string, 64)
	if err != nil {
		return getErrorResponse(`Not a valid distance`)
	}
	const mtoMi float64 = 0.00062137

	var wg sync.WaitGroup
	wg.Add(len(responses))

	var result []ResponseNew
	for i := 0; i < len(responses); i++ {
		go func (i int) {
			defer wg.Done()
			temp := new(ResponseNew)
			temp.Org = []float64{responses[i].StartLocation[0], responses[i].StartLocation[1]}
			temp.Dest = []float64{responses[i].EndLocation[0], responses[i].EndLocation[1]}
			oLat := strconv.FormatFloat(temp.Org[0], 'f', 6, 64)
			oLng := strconv.FormatFloat(temp.Org[1], 'f', 6, 64)
			dLat := strconv.FormatFloat(temp.Dest[0], 'f', 6, 64)
			dLng := strconv.FormatFloat(temp.Dest[1], 'f', 6, 64)
			path := polylineToPath(responses[i].Polyline)
			//too many waypoints, limiting to 23 would lose a lot of info about segment
			if (len(path) < (23 * 2 * 2)) {
				//limits to 23 waypoints
				if (len(path) > (23 * 2)) {
					var lessPath []float64
					removeNum := len(path) / 2 - 23
					increm := len(path) / 2 / removeNum
					if increm == 1 { fmt.Println("should not get here") } 
					j := 0
					for i := 0; i < len(path); i+=2 {
						j += 1
						if (j + 1) % increm != 0 {
							lessPath = append(lessPath, path[i], path[i + 1])
						}
					}
					path = lessPath
				}
				temp.Path = path
				tempPolyline := pathToPolyline(path)
				dirURL := "https://maps.googleapis.com/maps/api/directions/json?origin=" + oLat + "," + oLng + "&destination=" + dLat + "," + dLng + MODE + "&waypoints=via:enc:" + tempPolyline + ":" + KEY
				dist, turnLocs, err := distanceHelp(dirURL)
				if err == `` {
					// temp.Distance = responses[i].Distance*mtoMi
					temp.Distance = dist*mtoMi //we're returning a modified route we should return the appropriate distance
					temp.Directions = turnLocs
					result = append(result, *temp)
				}
			}
		} (i)
		
	}
	wg.Wait()

	if len(result) <= 0 {
		return getErrorResponse(`No displayable routes found`)
	}

	//sorts difference desired distance from path distance from least to greatest
	sort.SliceStable(result, func(i, j int) bool {
		return math.Abs(result[i].Distance-distance) < math.Abs(result[j].Distance-distance)
	})

	//return max 4 strava responses
	if len(result) > 4 {
		result = result[0:4]
	}

	return newFullResponse(result)

}

func ExecuteStrava(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var req Request
	_ = json.NewDecoder(r.Body).Decode(&req)

	StravaResponse := ExecuteStravaRequest(req.Address, req.Distance, "10", 1.0, `9ce41b35c7a159db331d3b0ae624ea8ced8ce595`)
	if err := json.NewEncoder(w).Encode(StravaResponse); err != nil {
		json.NewEncoder(w).Encode(getErrorResponse(fmt.Sprintf(`%s`, err)))
	}
}
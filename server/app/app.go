package app

import (
	// "bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"

	// "os"
	"sort"
	"strconv"
	"strings"
	"github.com/strava/go.strava"
)

const KEY = "&key=AIzaSyB32cCcL4gD_WIYPP6dAVSprY_QYE3arsk"
const GEOCODE_URL = "https://maps.googleapis.com/maps/api/geocode/json?address="
const DISTANCE_URL = "https://maps.googleapis.com/maps/api/distancematrix/json?units=imperial&"
const NEAREST_RODE_URL = "https://roads.googleapis.com/v1/nearestRoads?points="
const MODE = "&mode=walking"
const EQUATOR_LENGTH = 69.172
const NINTERSECT_URL = "http://api.geonames.org/findNearestIntersectionJSON?lat="
const GEOKEY = "&username=gulab"

type Request struct {
	Address  string `json:"address"`
	Distance string `json:"distance"`
}

type TokenResposne struct {
	AccessToken string `json:"accessToken"`
}

type Response struct {
	Path         [][]float64
	Distance     []float64
	PercentError float64
	Error        string
}

type FullResponse struct {
	Results []ResponseNew
	Error string
}

type ResponseNew struct {
	Org 		[]float64
	Dest		[]float64
	Path		[]float64
	Distance 	float64
	Directions 	[]LocOfTurn
}

type StravaResponse struct {
	Path  []float64
	Start []float64
	End   []float64
	Error string
}

type GeocodeGeometry struct {
	Location      map[string]interface{} `json:"location"`
	Location_type string                 `json:"location_type"`
	Viewport      interface{}            `json:"viewport"`
}

type GeocodeResults struct {
	Access_points      interface{}     `json:"access_points"`
	Address_components interface{}     `json:"address_components"`
	Formatted_address  string          `json:"formatted_address"`
	Geometry           GeocodeGeometry `json:"geometry"`
	Place_id           string          `json:"place_id"`
	Plus_code          interface{}     `json:"plus_code"`
	Types              interface{}     `json:"types"`
}

type GeocodeResp struct {
	Results []GeocodeResults `json:"results"`
	Status  string           `json:"status"`
}

type ValText struct {
	Value int    `json:"value"`
	Text  string `json:"text"`
}

type DistanceElems struct {
	Status   string  `json:"status"`
	Duration ValText `json:"duration"`
	Distance ValText `json:"distance"`
}

type DistanceRows struct {
	Elements []DistanceElems `json:"elements"`
}

type DistanceResp struct {
	DestAdds []string       `json:"destination_addresses"`
	OrgAdds  []string       `json:"origin_addresses"`
	Rows     []DistanceRows `json:"rows"`
	Status   string         `json:"status"`
}

type NearestPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type NearestLocation struct {
	Location      NearestPoint `json:"location"`
	OriginalIndex int          `json:"originalIndex"`
	PlaceId       string       `json:"placeId"`
}
type NearestResp struct {
	SnappedPoints []NearestLocation `json:"snappedPoints"`
}

type NearIntersectResp struct {
	Credits      string            `json:"credits"`
	Intersection map[string]string `json:"intersection"`
}

type Point struct {
	lat float64
	lng float64
}

type DistAndPath struct {
	path     []Point
	distance float64
	desired  float64
	turns	 []LocOfTurn
}

type Pt struct {
	Lat 		float64		`json:"lat"`
	Lng 		float64		`json:"lng"`
}

type Steps struct {
	Maneuver	string 		`json:"maneuver"`
	LocStep		Pt			`json:"start_location"`
	Html_Instructions string `json:"html_instructions"`
}

type Legs struct {
	Distance 	ValText		`json:"distance"`
	Steps 		[]Steps 	`json:"steps"`

}

type DirRoutes struct {
	Bounds		interface{}	`json:"bounds"`
	Copyright 	string		`json:"copyrights"`
	Legs 		[]Legs		`json:"legs"`
}

type DirResp struct {
	Geowpts 	interface{} `json:"geocoded_waypoints"`
	Rt 			[]DirRoutes `json:"routes"`
	Status		string 		`json:"status"`
}

type LocOfTurn struct {
	Turn 		string
	Instructions string
	Loc			[]float64
}

type make_route func(point Point, distance float64, offset float64) []Point

//creates a new point
func NewPoint(lat float64, lng float64) Point {
	this := new(Point)
	this.lat = lat
	this.lng = lng
	return *this
}

//forms an address for an api call
func address_to_api_call(address string) string {
	//properly form the address for the api url
	arr := strings.Split(address, " ")
	var address_url string
	for _, s := range arr {
		address_url += "+" + s
	}
	address_url = strings.TrimPrefix(address_url, "+")

	//create the api url and get the json response
	url := GEOCODE_URL + address_url + KEY
	return url
}

//makes a GET request at the given URL and checks for error
func api_request(url string) ([]byte, string) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, `Get request failed`
	}
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, `IOutil Failed`
	}

	return response, ``
}

//Checks a GeoCode Response body for error
func check_responseGeocode(response []byte) bool {
	var resp_body GeocodeResp
	json.Unmarshal(response, &resp_body)
	status := resp_body.Status
	if status == "OK" {
		return true
	}
	return false

}

//checks a Distance Reponse body for error
func check_responseDistance(response []byte) bool {
	var resp_body DistanceResp
	json.Unmarshal(response, &resp_body)
	status := resp_body.Status
	if status == "OK" {
		return true
	}
	return false

}

//Gets coordinates from a GeoCode Response Body
func extract_coordinates(response []byte) (float64, float64) {
	var resp_body GeocodeResp
	json.Unmarshal(response, &resp_body)
	coordinates := resp_body.Results[0].Geometry.Location

	return coordinates["lat"].(float64), coordinates["lng"].(float64)
}

func nearestIntersectionPoint(point Point) (Point, string) {

	url := NINTERSECT_URL + strconv.FormatFloat(point.lat, 'f', 6, 64) + "&lng=" + strconv.FormatFloat(point.lng, 'f', 6, 64) + GEOKEY

	response, err := api_request(url)

	if (err != ``) {
		return *new(Point), err
	}

	var resp_body NearIntersectResp
	json.Unmarshal(response, &resp_body)
	var resLat float64
	var resLng float64

	if s1, err := strconv.ParseFloat(resp_body.Intersection["lat"], 64); err == nil {
		resLat = s1
	} else {
		return *new(Point), `No nearest intersection point found`
	}

	if s2, err := strconv.ParseFloat(resp_body.Intersection["lng"], 64); err == nil {
		resLng = s2
	} else {
		return *new(Point), `No nearest intersection point found`
	}

	return NewPoint(resLat, resLng), ``


}

// given an origin, distance and angle, finds the corresponding point
func get_point(point Point, distance float64, angle float64, runNearestIntersection bool) (Point, string){
	//angle is degrees
	radians := angle * math.Pi / 180
	distance_lat := 1 / (69 / distance)
	one_degree_lng := math.Cos(point.lat*math.Pi/180) * EQUATOR_LENGTH
	distance_lng := 1 / (one_degree_lng / distance)

	lat := point.lat + math.Cos(radians)*distance_lat
	lng := point.lng + math.Sin(radians)*distance_lng
	
	if (!runNearestIntersection) {
		return NewPoint(lat, lng), ``
	}

    return nearestIntersectionPoint(NewPoint(lat, lng))
}

func points_to_angle(p1 Point, p2 Point) float64 {
	slope := (p2.lng - p1.lng) / (p1.lat - p2.lat)
	radians := math.Atan(slope)
	degrees := radians * 180 / math.Pi
	return degrees

}

func get_offset(point Point) float64 {
	url := NEAREST_RODE_URL + strconv.FormatFloat(point.lat, 'f', 6, 64) + "," + strconv.FormatFloat(point.lng, 'f', 6, 64) + KEY
	response, err := api_request(url)
	if (err != ``) {
		// this is okay because we can still try to make a request if offset doesnt exist
		return 0.0
	}
	var resp_body NearestResp
	json.Unmarshal(response, &resp_body)
	new_point := NewPoint(resp_body.SnappedPoints[0].Location.Latitude, resp_body.SnappedPoints[0].Location.Longitude)
	angle := points_to_angle(point, new_point)
	fmt.Println(angle)
	return angle
}

//Given an origin, desired distance, number of different routes it wants and a function that
//determines the shape of the route, returns a list of possible routes
func create_routes(point Point, distance float64, num float64, make_route make_route) [][]Point {
	angle_increase := 360 / num
	//seems to me like it works better without offset
	//offset := get_offset(point)
	offset := 0.0
	var routes [][]Point

	for i := 0.0; i < num; i = i + 1.0 {
		route := make_route(point, distance, offset)
		if (len(route) > 0) {
			routes = append(routes, route)
		}
		offset += angle_increase
		if (offset > 360) {
			offset = offset - 360
		}
	}
	return routes
}

//creates a straight line route
var straight_line make_route = func(point Point, distance float64, offset float64) []Point {

	p0, err := get_point(point, distance/2, offset, true)

	if (err != ``) {
		return []Point{}
	}

	return []Point{p0}
}

//creates a square route
var square_route make_route = func(point Point, distance float64, offset float64) []Point {

	side_length := (distance / 4)

	p0, err_0 := get_point(point, side_length, offset+135.0, true)
	p1, err_1 := get_point(point, distance/4*math.Sqrt(2), offset+90.0, true)
	p2, err_2 := get_point(point, side_length, offset+45.0, true)

	if (err_0 != `` || err_1 != `` || err_2 != ``) {
		return []Point{}
	}

	return []Point{p0, p1, p2}
}

var equilateral_triangle make_route = func(point Point, distance float64, offset float64) []Point {
	side_length := (distance / 3.0)

	p0, err_0 := get_point(point, side_length, offset+60.0,true)
	p1, err_1 := get_point(point, side_length, offset+120.0,true)

	if (err_0 != `` || err_1 != ``) {
		return []Point{}
	}

	return []Point{p0, p1}
}

//isosceles
var right_triangle make_route = func(point Point, distance float64, offset float64) []Point {
	side_length := distance / (2.0 + math.Sqrt(2.0))

	p0, err_0 := get_point(point, side_length, offset,true)
	p1, err_1 := get_point(point, side_length, offset+90.0,true)

	if (err_0 != `` || err_1 != ``) {
		return []Point{}
	}

	return []Point{p0, p1}
}

//isosceles
var right_triangleOther make_route = func(point Point, distance float64, offset float64) []Point {
	side_length := distance / (2.0 + math.Sqrt(2.0))
	p0, err_0 := get_point(point, side_length, offset+90.0,true)
	
	if (err_0 != ``) {
		return []Point{}
	}

	p1, err_1 := get_point(p0, side_length, offset,true)

	if (err_1 != ``) {
		return []Point{}
	}

	return []Point{p0, p1}
}

//finds distance between 2 given cooridnates
func distance(oLat float64, oLng float64, dLat float64, dLng float64) (int, string) {
	orgLat := strconv.FormatFloat(oLat, 'f', 6, 64)
	orgLng := strconv.FormatFloat(oLng, 'f', 6, 64)
	dstLat := strconv.FormatFloat(dLat, 'f', 6, 64)
	dstLng := strconv.FormatFloat(dLng, 'f', 6, 64)
	urlCall := DISTANCE_URL + "origins=" + orgLat + "," + orgLng + "&destinations=" + dstLat + "," + dstLng + MODE + KEY
	
	response, err := api_request(urlCall)
	if (err != ``) {
		return 0, err
	}

	if !check_responseDistance(response) {
		return 0, "Response Distance Status not okay"
	}

	var resp_body DistanceResp
	json.Unmarshal(response, &resp_body)

	return resp_body.Rows[0].Elements[0].Distance.Value, ``

}

func distanceHelp(dirURL string) (float64, []LocOfTurn, string) {
	response, err := api_request(dirURL)
	var result []LocOfTurn
	if (err != ``) {
		return 0.0, result, err
	}

	var resp_body DirResp
	json.Unmarshal(response, &resp_body)

	if (resp_body.Status == "OK") {

		for _, v := range resp_body.Rt[0].Legs[0].Steps {
			turnLocs := new(LocOfTurn)
			turnLocs.Turn = v.Maneuver
			turnLocs.Instructions = v.Html_Instructions
			temp := []float64{v.LocStep.Lat, v.LocStep.Lng}
			turnLocs.Loc = temp
			result = append(result, *turnLocs)
		}
		return float64(resp_body.Rt[0].Legs[0].Distance.Value), result, ""
	}
	fmt.Println("Status not okay")
	return 0.0, result, resp_body.Status
}

func get_distance(pathSlice [][]Point, org Point, desired float64) []DistAndPath {
	var allDists []DistAndPath
	oLat := strconv.FormatFloat(org.lat, 'f', 6, 64)
	oLng := strconv.FormatFloat(org.lng, 'f', 6, 64)
	for _, v := range pathSlice {
		tempUrl := "https://maps.googleapis.com/maps/api/directions/json?origin=" + oLat + "," + oLng + "&destination=" + oLat + "," + oLng + MODE + "&waypoints="
		for _, x := range v {
			cLat := strconv.FormatFloat(x.lat, 'f', 6, 64)
			cLng := strconv.FormatFloat(x.lng, 'f', 6, 64)
			tempUrl = tempUrl + "via:" + cLat + "," + cLng + "|"
		}
		tempUrl = strings.TrimSuffix(tempUrl, "|")
		url := tempUrl + KEY
		dist, turnLocs, err := distanceHelp(url)
		fmt.Println(turnLocs)
		fmt.Println(dist * 0.00062137, "for above turns")
		if (err == ``) {
			// convert output in meters to miles
			const mtoMi float64 = 0.00062137
			distMi := dist * mtoMi
			//only includes paths which are at least the desired length and at most desired length + 1 mi
			if ((distMi > desired) && (distMi < (desired + 1))) {
				temp := new(DistAndPath)
				temp.distance = distMi
				temp.path = v
				temp.desired = desired
				temp.turns = turnLocs
				allDists = append(allDists, *temp)
			}
		}
	}
	return allDists
}

//given paths, origin, and desired distance, returns slice of all distance lengths and paths
func getDistance(pathSlice [][]Point, org Point, desired float64) []DistAndPath {
	var allDists []DistAndPath
	for _, v := range pathSlice {

		totalDist := 0
		prevLat := org.lat
		prevLng := org.lng
		overall_error := ``
		for _, x := range v {
			curLat := x.lat
			curLng := x.lng
			some, err  := distance(prevLat, prevLng, curLat, curLng)
			if (err == ``) {
				totalDist += some 
			} else {
				overall_error = err
				break
			}
			prevLat = curLat
			prevLng = curLng
		}

		some, err := distance(prevLat, prevLng, org.lat, org.lng)
		totalDist += some

		if (overall_error == `` && err == ``) {
			// convert output in meters to miles
			const mtoMi float64 = 0.00062137
			//only includes paths which are at least the desired length
			if float64(totalDist)*mtoMi > desired {
				// if true {
				temp := new(DistAndPath)
				temp.distance = float64(totalDist) * mtoMi
				temp.path = v
				temp.desired = desired
				allDists = append(allDists, *temp)
			}
		}

	}

	return allDists
}

func getErrorResponse(error string) *FullResponse {
	this := new(FullResponse)
	this.Results = nil
	this.Error = error
	return this
}

func newFullResponse(results []ResponseNew) *FullResponse {
	this := new(FullResponse)
	this.Results = results
	this.Error = ``
	return this
}

//given an address, distance and route, finds a path
func execute_request(input string, distance_string string, error_fix float64) *FullResponse {
	var routeOrder []make_route
	routeOrder = append(routeOrder, square_route, equilateral_triangle, right_triangle, right_triangleOther, straight_line)
	//convert distance to float64
	//check to see if distance is a proper number
	distance, err := strconv.ParseFloat(distance_string, 64)
	if err != nil {
		return getErrorResponse("Not a valid distance")
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
		return getErrorResponse("Address does not exist")
	}

	//get the latitude and longitude
	lat, lng := extract_coordinates(response)
	origin := NewPoint(lat, lng)
	fmt.Println(origin)

	//get the possible routes
	routes := create_routes(origin, distance*(error_fix), 8.0, routeOrder[0])

	desired := distance
	pathDetails := get_distance(routes, origin, float64(desired))

	attempts := 0
	for len(pathDetails) < 8 {
		attempts += 1
		if attempts > 4 {
			fmt.Println(`At least we tried`)
			break
		}
		fmt.Println(len(pathDetails), "returnable routes")
		fmt.Println("attempt", attempts)
		routes = create_routes(origin, distance*(error_fix), 8.0, routeOrder[attempts])
		pathDetails = append(pathDetails, get_distance(routes, origin, float64(desired))...)

	}
	//checks to make sure there are paths found
	if len(pathDetails) == 0 {
		return getErrorResponse("No paths found")
	}
	fmt.Println(len(pathDetails))

	//sorts difference desired distance from path distance from least to greatest
	sort.SliceStable(pathDetails, func(i, j int) bool {
		return math.Abs(pathDetails[i].distance-pathDetails[i].desired) < math.Abs(pathDetails[j].distance-pathDetails[j].desired)
	})

	//percent_error := (pathDetails[0].distance - pathDetails[0].desired) / pathDetails[0].desired * 100
	//Outputs best path with distance and percent error

	if len(pathDetails) > 8 {
		//only want best 8 if there are more then 8
		pathDetails = pathDetails[0:8]
	}
	
	//@Agam what does this do?
	//it creates a 2d array of the appropriate size
	// resPaths := make([][]float64, len(pathDetails))
	// var resDists []float64

	// var i int
	// for _, elem := range pathDetails {
	// 	fmt.Println(elem)
	// 	resPaths[i] = append(resPaths[i], origin.lat, origin.lng)
	
	// 	for _, pt := range pathDetails[i].path {
	// 		curLat := pt.lat
	// 		curLng := pt.lng
	// 		resPaths[i] = append(resPaths[i], curLat, curLng)
	// 	}
	// 	resDists = append(resDists, pathDetails[i].distance)
	// 	i++
	// }


	var result []ResponseNew
	for i := 0; i < len(pathDetails); i++ {
		temp := new(ResponseNew)
		temp.Org = append(temp.Org, origin.lat, origin.lng)
		temp.Dest = append(temp.Dest, origin.lat, origin.lng)
		for _, pt := range pathDetails[i].path {
			temp.Path = append(temp.Path, pt.lat, pt.lng)
		}
		temp.Distance = pathDetails[i].distance
		temp.Directions = pathDetails[i].turns
		result = append(result, *temp)
	}
	fmt.Println(result)


	return newFullResponse(result)
}



// Thinking all we have to do is pass result of execute request into json.NewEncoder thing but not too sure
func Execute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var req Request
	_ = json.NewDecoder(r.Body).Decode(&req)
	//path, distance, percent_error, err := execute_request(req.Address, req.Distance, 1.0)
	//json.NewEncoder(w).Encode(newResponse(path, distance, percent_error, err))
	json.NewEncoder(w).Encode(execute_request(req.Address, req.Distance, 1.0))
}

func polylineToPath(polyline strava.Polyline) []float64 {
	coords := polyline.Decode()
	path := []float64{}
	for _, coord := range coords {
		path = append(path, coord[0])
		path = append(path, coord[1])
	}

	return path
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
	fmt.Println(top_right)
	fmt.Println(bottom_left)
	client := strava.NewClient(accessToken)
	SegmentCall := strava.NewSegmentsService(client).Explore(bottom_left.lat, bottom_left.lng, top_right.lat, top_right.lng)
	SegmentCall.ActivityType("running")
	SegmentCall.MinimumCategory(1)
	SegmentCall.MaximumCategory(100)

	responses, err := SegmentCall.Do()
	if err != nil {
		return getErrorResponse(fmt.Sprintf(`%s`, err))
	}

	fmt.Println(len(responses))

	if len(responses) <= 0 {
		return getErrorResponse(`No routes found`)
	}

	distance, err := strconv.ParseFloat(distance_string, 64)
	if err != nil {
		return getErrorResponse(`Not a valid distance`)
	}
	const mtoMi float64 = 0.00062137

	//sorts difference desired distance from path distance from least to greatest
	sort.SliceStable(responses, func(i, j int) bool {
		return math.Abs(responses[i].Distance*mtoMi-distance) < math.Abs(responses[j].Distance*mtoMi-distance)
	})

	//return max 4 strava responses
	if len(responses) > 4 {
		responses = responses[0:4]
	}
	// best_distance_index := 0
	// best_distance_difference := math.Abs(distance - responses[0].Distance)

	// for i, resp := range responses {
	// 	if math.Abs(distance-resp.Distance) < best_distance_difference {
	// 		best_distance_index = i
	// 		best_distance_difference = math.Abs(distance - resp.Distance)
	// 	}
	// }

	//start := responses[0].StartLocation
	//end := responses[0].EndLocation
	path := polylineToPath(responses[0].Polyline)
	fmt.Println(path)
	if (len(path) > (23 * 2)) {
		lessPath := make([]float64, 23 * 2)
		incr := len(path) / 2 / 23
		for i := 0; i < 46; i+=2 {
			lessPath[i] = path[i * incr]
			lessPath[i + 1] = path[i * incr + 1]
		}
		path = lessPath
	}


	resPaths := make([][]float64, len(responses))
	resStarts := make([][]float64, len(responses))
	resEnds := make([][]float64, len(responses))
	var resDists []float64
	for i := 0; i < len(responses); i++ {
		resPaths[i] = polylineToPath(responses[i].Polyline)
		resStarts[i] = []float64{responses[i].StartLocation[0], responses[i].StartLocation[1]}
		resEnds[i] = []float64{responses[i].EndLocation[0], responses[i].EndLocation[1]}
		resDists = append(resDists, responses[i].Distance*mtoMi)
	}
	fmt.Println(resPaths)
	fmt.Println(resStarts)
	fmt.Println(resEnds)
	fmt.Println(resDists)

	var result []ResponseNew
	for i := 0; i < len(responses); i++ {
		temp := new(ResponseNew)
		temp.Org = []float64{responses[i].StartLocation[0], responses[i].StartLocation[1]}
		temp.Dest = []float64{responses[i].EndLocation[0], responses[i].EndLocation[1]}
		temp.Path = polylineToPath(responses[i].Polyline)
		temp.Distance = responses[i].Distance*mtoMi
		temp.Directions = nil
		result = append(result, *temp)
	}
	
	return newFullResponse(result)

}

// This might need to be changed too see line 670
func ExecuteStrava(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var req Request
	_ = json.NewDecoder(r.Body).Decode(&req)

	// get token from python client 
	resp, err := api_request("http://localhost:5000/token")

	if (err != ``) {
		json.NewEncoder(w).Encode(getErrorResponse(err))
		return
	}

	var resp_body TokenResposne
	json.Unmarshal(resp, &resp_body)

	StravaResponse := ExecuteStravaRequest(req.Address, req.Distance, "10", 1.0, resp_body.AccessToken)
	json.NewEncoder(w).Encode(StravaResponse)
}

func Tester(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var req Request
	_ = json.NewDecoder(r.Body).Decode(&req)

	temp := new(ResponseNew)
	temp.Org = []float64{37.2862852, -122.0081542}
	temp.Dest = []float64{37.2862852, -122.0081542}
	temp.Path = []float64{37.28319,-122.01359, 37.28837, -122.01194}
	temp.Distance = 1.39870387


	stop_1 := new(LocOfTurn)
	stop_1.Instructions = "Head <b>south</b> on <b>Palmtag Dr</b> toward <b>Bellwood Dr</b>"
	stop_1.Loc = []float64{37.2862401,  -122.0078964}
	stop_1.Turn =  ""

	stop_2 := new(LocOfTurn)
	stop_2.Instructions = "Turn <b>right</b> onto <b>Bellwood Dr</b>"
	stop_2.Loc = []float64{37.2856419, -122.0080037}
	stop_2.Turn =  "turn-right"

	stop_3 := new(LocOfTurn)
	stop_3.Instructions = "Turn <b>left</b> onto <b>Titus Ave</b>"
	stop_3.Loc = []float64{37.2857637, -122.0100177}
	stop_3.Turn =  "turn-left"

	temp.Directions = []LocOfTurn{*stop_1, *stop_2, *stop_3}



	temp2 := new(ResponseNew)
	temp2.Org = []float64{37.2862852, -122.0081542}
	temp2.Dest = []float64{37.2862852, -122.0081542}
	temp2.Path = []float64{37.2866, -122.01331, 37.28257, -122.01174, 37.28244, -122.00863}
	temp2.Distance = 1.39870387

	temp2.Directions = []LocOfTurn{*stop_1, *stop_2, *stop_3}

	temp3 := new(ResponseNew)
	temp3.Org = []float64{37.2862852, -122.0081542}
	temp3.Dest = []float64{37.2862852, -122.0081542}
	temp3.Path = []float64{37.28352, -122.01188, 37.28181, -122.00754, 37.2837, -122.00532}
	temp3.Distance = 1.39870387

	temp3.Directions = []LocOfTurn{*stop_1, *stop_2, *stop_3}

	response := newFullResponse([]ResponseNew{*temp, *temp2, *temp3})
	fmt.Println(response)

	json.NewEncoder(w).Encode(response)

}

// func main() {

// 	//setup the scanner
// 	scanner := bufio.NewScanner(os.Stdin)

// 	//ask for address
// 	fmt.Printf("Enter your address in the following format: Street Address, City, State\n")

// 	//read user input
// 	var input string
// 	for scanner.Scan() {
// 		input = scanner.Text()
// 		if strings.Contains(input, "") {
// 			break
// 		}
// 	}

// 	//ask for distance
// 	// fmt.Printf("Enter desired distance in miles\n")

// 	//read user input
// 	// var distance string
// 	// for scanner.Scan() {
// 	// 	distance = scanner.Text()
// 	// 	if strings.Contains(distance, "") {
// 	// 		break
// 	// 	}
// 	// }
// 	distance := "1"
// 	execute_request(input, distance, 1.0)
// 	fmt.Println("strava response below")
// 	//ExecuteStravaRequest(input, distance, "10", 1.0, `d11fe6a257303cf24de1181586f335f86d24db88`)
// }

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

const KEY = "&key=AIzaSyB32cCcL4gD_WIYPP6dAVSprY_QYE3arsk"
const GEOCODE_URL = "https://maps.googleapis.com/maps/api/geocode/json?address="
const DISTANCE_URL = "https://maps.googleapis.com/maps/api/distancematrix/json?units=imperial&"
const NEAREST_RODE_URL = "https://roads.googleapis.com/v1/nearestRoads?points="
const MODE = "&mode=walking"
const EQUATOR_LENGTH = 69.172

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

type valText struct {
	Value int    `json:"value"`
	Text  string `json:"text"`
}

type DistanceElems struct {
	Status   string  `json:"status"`
	Duration valText `json:"duration"`
	Distance valText `json:"distance"`
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
	Latitude float64 	`json:"latitude"`
	Longitude float64  	`json:"longitude"`
}

type NearestLocation struct {
	Location NearestPoint  `json:"location"`
	OriginalIndex int      `json:"originalIndex"`
	PlaceId string         `json:"placeId"`
}
type NearestResp struct {
	SnappedPoints []NearestLocation `json:"snappedPoints"`
}


type Point struct {
	lat float64
	lng float64
}

type DistAndPath struct {
	path     []Point
	distance float64
	desired  float64
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
func api_request(url string) []byte {

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return response
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

// given an origin, distance and angle, finds the corresponding point
func get_point(point Point, distance float64, angle float64) Point {
	//angle is degrees
	radians := angle * math.Pi / 180
	distance_lat := 1 / (69 / distance)
	one_degree_lng := math.Cos(point.lat*math.Pi/180) * EQUATOR_LENGTH
	distance_lng := 1 / (one_degree_lng / distance)

	return NewPoint(point.lat+math.Cos(radians)*distance_lat, point.lng+math.Sin(radians)*distance_lng)
}

func points_to_angle (p1 Point, p2 Point) float64 {
	slope := (p2.lng - p1.lng) / (p1.lat - p2.lat)
	radians := math.Atan(slope)
	degrees := radians * 180 / math.Pi
	return degrees

}
func get_offset(point Point) float64 {
	url := NEAREST_RODE_URL +  strconv.FormatFloat(point.lat, 'f', 6, 64) + "," + strconv.FormatFloat(point.lng, 'f', 6, 64) + KEY
	response := api_request(url)
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
	offset := get_offset(point)
	var routes [][]Point

	for offset < 360 {
		route := make_route(point, distance, offset)
		routes = append(routes, route)
		offset += angle_increase
	}
	return routes
}

//creates a straight line route
var straight_line make_route = func(point Point, distance float64, offset float64) []Point {
	return []Point{get_point(point, distance/2, offset)}
}

//creates a square route
var square_route make_route = func(point Point, distance float64, offset float64) []Point {

	side_length := (distance / 4)

	p0 := get_point(point, side_length, offset+135.0)
	p1 := get_point(point, distance/4*math.Sqrt(2), offset+90.0)
	p2 := get_point(point, side_length, offset+45.0)

	return []Point{p0, p1, p2}
}

var equilateral_triangle make_route = func(point Point, distance float64, offset float64) []Point {
	side_length := (distance / 3.0)

	p0 := get_point(point, side_length, offset+60.0)
	p1 := get_point(point, side_length, offset+120.0)

	return []Point{p0, p1}
}

//isosceles
var right_triangle make_route = func(point Point, distance float64, offset float64) []Point {
	side_length := distance / (2.0 + math.Sqrt(2.0))
	p0 := get_point(point, side_length, offset)
	p1 := get_point(point, side_length, offset+90.0)

	return []Point{p0, p1}
}

//isosceles
var right_triangleOther make_route = func(point Point, distance float64, offset float64) []Point {
	side_length := distance / (2.0 + math.Sqrt(2.0))
	p0 := get_point(point, side_length, offset+90.0)
	p1 := get_point(p0, side_length, offset)

	return []Point{p0, p1}
}

//finds distance between 2 given cooridnates
func distance(oLat float64, oLng float64, dLat float64, dLng float64) int {
	orgLat := strconv.FormatFloat(oLat, 'f', 6, 64)
	orgLng := strconv.FormatFloat(oLng, 'f', 6, 64)
	dstLat := strconv.FormatFloat(dLat, 'f', 6, 64)
	dstLng := strconv.FormatFloat(dLng, 'f', 6, 64)
	urlCall := DISTANCE_URL + "origins=" + orgLat + "," + orgLng + "&destinations=" + dstLat + "," + dstLng + MODE + KEY
	response := api_request(urlCall)

	if !check_responseDistance(response) {
		fmt.Println("Status not okay")
	}

	var resp_body DistanceResp
	json.Unmarshal(response, &resp_body)

	return resp_body.Rows[0].Elements[0].Distance.Value

}

//given paths, origin, and desired distance, returns slice of all distance lengths and paths
func getDistance(pathSlice [][]Point, org Point, desired float64) []DistAndPath {
	var allDists []DistAndPath
	for _, v := range pathSlice {
		someDist := 0
		prevLat := org.lat
		prevLng := org.lng
		for _, x := range v {
			curLat := x.lat
			curLng := x.lng
			someDist += distance(prevLat, prevLng, curLat, curLng)
			prevLat = curLat
			prevLng = curLng
		}
		someDist += distance(prevLat, prevLng, org.lat, org.lng)
		// convert output in meters to miles
		const mtoMi float64 = 0.00062137
		//only includes paths which are at least the desired length
		if float64(someDist)*mtoMi > desired {
			// if true {
			temp := new(DistAndPath)
			temp.distance = float64(someDist) * mtoMi
			temp.path = v
			temp.desired = desired
			allDists = append(allDists, *temp)
		}

	}
	//fmt.Println(allDists)
	return allDists
}

//given an address, distance and route, finds a path
func execute(input string, distance_string string, route_function make_route, error_fix float64) (bool, float64) {

	//convert distance to float64
	//check to see if distance is a proper number
	distance, err := strconv.ParseFloat(distance_string, 64)
	if err != nil {
		fmt.Println("Not a valid distance")
		return false, 0
	}

	//form url from address
	url := address_to_api_call(input)

	//get response from google api server
	response := api_request(url)

	//check to see if address exists
	if !check_responseGeocode(response) {
		fmt.Println("Address doesn't exist")
		return false, 0
	}

	//get the latitude and longitude
	lat, lng := extract_coordinates(response)
	origin := NewPoint(lat, lng)

	//get the possible routes
	routes := create_routes(origin, distance*(error_fix), 8.0, route_function)
	//desired distance hard coded as 1 mile
	desired := 1
	pathDetails := getDistance(routes, origin, float64(desired))

	//checks to make sure there are paths found
	if len(pathDetails) == 0 {
		fmt.Println("No paths found")
		return false, 0
	}

	//sorts difference desired distance from path distance from least to greatest
	sort.SliceStable(pathDetails, func(i, j int) bool {
		return math.Abs(pathDetails[i].distance-pathDetails[i].desired) < math.Abs(pathDetails[j].distance-pathDetails[j].desired)
	})
	fmt.Println(pathDetails)

	percent_error := (pathDetails[0].distance - pathDetails[0].desired) / pathDetails[0].desired * 100
	//Outputs best path with distance and percent error
	fmt.Println("Best Path:", pathDetails[0].path, "\nDistance:", pathDetails[0].distance, "\nPercent Error:", percent_error, "%")

	return true, percent_error
}

func main() {

	//setup the scanner
	scanner := bufio.NewScanner(os.Stdin)

	//ask for address
	fmt.Printf("Enter your address in the following format: Street Address, City, State\n")

	//read user input
	var input string
	for scanner.Scan() {
		input = scanner.Text()
		if strings.Contains(input, "") {
			break
		}
	}

	//ask for distance
	fmt.Printf("Enter desired distance in miles\n")

	//read user input
	var distance string
	for scanner.Scan() {
		distance = scanner.Text()
		if strings.Contains(distance, "") {
			break
		}
	}

	execute(input, distance, straight_line, 1.0)

}

// func main() {
// 	origin := NewPoint(37.2864076,-122.0081492)
// 	get_offset(origin)
// }

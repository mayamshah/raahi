package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const KEY = "&key=AIzaSyB32cCcL4gD_WIYPP6dAVSprY_QYE3arsk"
const URL = "https://maps.googleapis.com/maps/api/geocode/json?address="
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

type Elems struct {
	Status   string  `json:"status"`
	Duration valText `json:"duration"`
	Distance valText `json:"distance"`
}

type Between struct {
	Elements []Elems `json:"elements"`
}

type DistanceResp struct {
	DestAdds []string  `json:"destination_addresses"`
	OrgAdds  []string  `json:"origin_addresses"`
	Rows     []Between `json:"rows"`
	Status   string    `json:"status"`
}

type Point struct {
	lat float64
	lng float64
}

func distance(oLat float64, oLng float64, dLat float64, dLng float64) int {
	orgLat := strconv.FormatFloat(oLat, 'f', 6, 64)
	orgLng := strconv.FormatFloat(oLng, 'f', 6, 64)
	dstLat := strconv.FormatFloat(dLat, 'f', 6, 64)
	dstLng := strconv.FormatFloat(dLng, 'f', 6, 64)
	const mode = "walking"
	const key = "&key=AIzaSyB32cCcL4gD_WIYPP6dAVSprY_QYE3arsk"
	const url = "https://maps.googleapis.com/maps/api/distancematrix/json?units=imperial&"
	urlCall := url + "origins=" + orgLat + "," + orgLng + "&destinations=" + dstLat + "," + dstLng + "&mode=" + mode + key
	response := api_request(urlCall)

	var resp_body DistanceResp
	json.Unmarshal(response, &resp_body)

	return resp_body.Rows[0].Elements[0].Distance.Value

}

func getDistance(pathSlice [][]Point, org Point) [8]float64 {
	var allDists [8]float64
	for i, v := range pathSlice {
		someDist := 0
		prevLat := org.lat
		prevLng := org.lng
		for _, x := range v {
			curLat := x.lat
			curLng := x.lng
			someDist += distance(prevLat, prevLng, curLat, curLng)
			print(someDist)
			prevLat = curLat
			prevLng = curLng
		}
		someDist += distance(prevLat, prevLng, org.lat, org.lng)
		const mtoMi float64 = 0.00062137
		allDists[i] = float64(someDist) * mtoMi
	}
	fmt.Println(allDists)
	return allDists
}

func address_to_api_call(address string) string {
	//properly form the address for the api url
	arr := strings.Split(address, " ")
	var address_url string
	for _, s := range arr {
		address_url += "+" + s
	}
	address_url = strings.TrimPrefix(address_url, "+")

	//create the api url and get the json response
	url := URL + address_url + KEY
	return url
}

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

func check_response(response []byte) bool {
	var resp_body GeocodeResp
	json.Unmarshal(response, &resp_body)
	status := resp_body.Status
	if status == "OK" {
		return true
	}
	return false

}

func extract_coordinates(response []byte) (float64, float64) {
	var resp_body GeocodeResp
	json.Unmarshal(response, &resp_body)
	coordinates := resp_body.Results[0].Geometry.Location

	return coordinates["lat"].(float64), coordinates["lng"].(float64)
}

func possible_routes_straight_line(lat float64, lng float64, distance float64) [][]Point {
	distance_lat := 1 / (69 / distance)
	one_degree_lng := math.Cos(lat*math.Pi/180) * EQUATOR_LENGTH
	distance_lng := 1 / (one_degree_lng / distance)
	root2_lat := math.Sqrt(2) * distance_lat
	root2_lng := math.Sqrt(2) * distance_lng

	p0 := Point{lat: lat + distance_lat, lng: lng}
	p1 := Point{lat: lat + root2_lat, lng: lng + root2_lng}
	p2 := Point{lat: lat, lng: lng + distance_lng}
	p3 := Point{lat: lat - root2_lat, lng: lng + root2_lng}
	p4 := Point{lat: lat - distance_lat, lng: lng}
	p5 := Point{lat: lat - root2_lat, lng: lng - root2_lng}
	p6 := Point{lat: lat, lng: lng - distance_lng}
	p7 := Point{lat: lat + root2_lat, lng: lng - root2_lng}

	return [][]Point{{p0}, {p1}, {p2}, {p3}, {p4}, {p5}, {p6}, {p7}}

}

func main() {

	//setup the scanner
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter your address in the following format: Street Address, City, State\n")

	//read user input
	var input string
	for scanner.Scan() {
		input = scanner.Text()
		if strings.Contains(input, "") {
			break
		}
	}

	//form url from address
	url := address_to_api_call(input)

	//get response from google api server
	response := api_request(url)

	//check to see if address exists
	if !check_response(response) {
		fmt.Println("Address doesn't exist")
		return
	}

	//get the latitude and longitude
	lat, lng := extract_coordinates(response)
	fmt.Println(lat, lng)
	//get the possible routes
	routes := possible_routes_straight_line(lat, lng, 0.5)

	fmt.Println(routes)
	var temp Point = Point{lat: lat, lng: lng}
	getDistance(routes, temp)
}

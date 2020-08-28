package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const KEY = "&key=AIzaSyB32cCcL4gD_WIYPP6dAVSprY_QYE3arsk"
const GEOCODE_URL = "https://maps.googleapis.com/maps/api/geocode/json?address="
const DISTANCE_URL = "https://maps.googleapis.com/maps/api/distancematrix/json?units=imperial&"
const NEAREST_RODE_URL = "https://roads.googleapis.com/v1/nearestRoads?points="
const TRAILS_URL = "https://www.trailrunproject.com/data/get-trails?lat="
const TRPKEY = "&key=200872374-da245a0b782e30b71d825d92027d8012"
const MODE = "&mode=walking"
const EQUATOR_LENGTH = 69.172
const NINTERSECT_URL = "http://api.geonames.org/findNearestIntersectionJSON?lat="
const GEOKEY = "&username=gulab"
const M_TO_MI = 0.00062137 //meters to miles ratio

type Request struct {
	Address  			string 					`json:"address"`
	Distance 			string 					`json:"distance"`
}

type TokenResposne struct {
	AccessToken 		string 					`json:"accessToken"`
}

type Response struct {
	Path         		[][]float64
	Distance     		[]float64
	PercentError 		float64
	Error        		string
}

type FullResponse struct {
	Results 			[]ResponseNew
	Error 				string
}

type TrailResponse struct {
	Results 			[]TrailInfo 
	Origin  			[]float64
	Error 				string
}

type ResponseNew struct {
	Org 				[]float64
	Dest				[]float64
	Path				[]float64
	Distance 			float64
	Directions 			[]LocOfTurn
}

type GeocodeGeometry struct {
	Location      		map[string]interface{}	`json:"location"`
	LocationType 		string                 	`json:"location_type"`
	Viewport      		interface{}            	`json:"viewport"`
}

type GeocodeResults struct {
	AccessPoints      	interface{}     		`json:"access_points"`
	AddressComponents 	interface{}     		`json:"address_components"`
	FormattedAddress  	string          		`json:"formatted_address"`
	Geometry           	GeocodeGeometry 		`json:"geometry"`
	PlaceID           	string          		`json:"place_id"`
	PlusCode         	interface{}     		`json:"plus_code"`
	Types              	interface{}     		`json:"types"`
}

type GeocodeResp struct {
	Results 			[]GeocodeResults 		`json:"results"`
	Status  			string           		`json:"status"`
}

type ValText struct {
	Value 				int    					`json:"value"`
	Text  				string 					`json:"text"`
}

type DistanceElems struct {
	Status   			string  				`json:"status"`
	Duration 			ValText 				`json:"duration"`
	Distance 			ValText 				`json:"distance"`
}

type DistanceRows struct {
	Elements 			[]DistanceElems 		`json:"elements"`
}

type DistanceResp struct {
	DestAdds 			[]string       			`json:"destination_addresses"`
	OrgAdds  			[]string       			`json:"origin_addresses"`
	Rows     			[]DistanceRows 			`json:"rows"`
	Status   			string         			`json:"status"`
}

type NearestPoint struct {
	Latitude  			float64 				`json:"latitude"`
	Longitude 			float64 				`json:"longitude"`
}

type NearestLocation struct {
	Location      		NearestPoint 			`json:"location"`
	OriginalIndex 		int          			`json:"originalIndex"`
	PlaceID       		string       			`json:"placeId"`
}
type NearestResp struct {
	SnappedPoints 		[]NearestLocation 		`json:"snappedPoints"`
}

type NearIntersectResp struct {
	Credits      		string            		`json:"credits"`
	Intersection 		map[string]string 		`json:"intersection"`
}

type Point struct {
	lat 				float64
	lng 				float64
}

type DistAndPath struct {
	path     			[]Point
	distance 			float64
	desired  			float64
	turns	 			[]LocOfTurn
}

type Pt struct {
	Lat 				float64					`json:"lat"`
	Lng 				float64					`json:"lng"`
}

type Steps struct {
	Maneuver			string 					`json:"maneuver"`
	LocStep				Pt						`json:"start_location"`
	LocStepEnd			Pt						`json:"end_location"`
	HtmlInstructions 	string 					`json:"html_instructions"`
	Dist 				ValText  				`json:"distance"`
}

type Legs struct {
	Distance 			ValText					`json:"distance"`
	Steps 				[]Steps 				`json:"steps"`
}

type DirRoutes struct {
	Bounds				interface{}				`json:"bounds"`
	Copyright 			string					`json:"copyrights"`
	Legs 				[]Legs					`json:"legs"`
}

type DirResp struct {
	Geowpts 			interface{} 			`json:"geocoded_waypoints"`
	Rt 					[]DirRoutes 			`json:"routes"`
	Status				string 					`json:"status"`
}

type LocOfTurn struct {
	Turn 				string
	Instructions 		string
	Loc					[]float64
	EndLoc 				[]float64
}

type TrailInfo struct {
	Name				string
	Summary				string
	Location			string
	Length				float64
	DistFromOrg			float64
	Coords 				[]float64
}

type Trails struct {
	Name				string					`json:"name"`
	Summary				string					`json:"summary"`
	Location			string					`json:"location"`
	Distance			float64					`json:"length"`
	Lat					float64					`json:"latitude"`
	Lon					float64					`json:"longitude"`
}

type TrailsResp struct {
	Trails 				[]Trails				`json:"trails"`
	Success				int						`json:"success"`
}

type makeRoute func(point Point, distance float64, offset float64) []Point

//creates a new point
func NewPoint(lat float64, lng float64) Point {
	this := new(Point)
	this.lat = lat
	this.lng = lng
	return *this
}

//forms an address for an api call
func addressToAPICall(address string) string {
	//properly form the address for the api url
	arr := strings.Split(address, " ")
	var addressURL string
	for _, s := range arr {
		addressURL += "+" + s
	}
	addressURL = strings.TrimPrefix(addressURL, "+")

	//create the api url and get the json response
	url := GEOCODE_URL + addressURL + KEY
	return url
}

//makes a GET request at the given URL and checks for error
func apiRequest(url string) ([]byte, string) {

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
func checkResponseGeocode(response []byte) bool {
	var respBody GeocodeResp
	json.Unmarshal(response, &respBody)
	status := respBody.Status
	if status == "OK" {
		return true
	}
	return false

}

//checks a Distance Reponse body for error
func checkResponseDistance(response []byte) bool {
	var respBody DistanceResp
	json.Unmarshal(response, &respBody)
	status := respBody.Status
	if status == "OK" {
		return true
	}
	return false

}

//Gets coordinates from a GeoCode Response Body
func extractCoordinates(response []byte) (float64, float64) {
	var respBody GeocodeResp
	json.Unmarshal(response, &respBody)
	coordinates := respBody.Results[0].Geometry.Location

	return coordinates["lat"].(float64), coordinates["lng"].(float64)
}

func nearestIntersectionPoint(point Point) (Point, string) {

	url := NINTERSECT_URL + strconv.FormatFloat(point.lat, 'f', 6, 64) + "&lng=" + strconv.FormatFloat(point.lng, 'f', 6, 64) + GEOKEY

	response, err := apiRequest(url)

	if (err != ``) {
		return *new(Point), err
	}

	var respBody NearIntersectResp
	json.Unmarshal(response, &respBody)
	var resLat float64
	var resLng float64

	if s1, err := strconv.ParseFloat(respBody.Intersection["lat"], 64); err == nil {
		resLat = s1
	} else {
		return *new(Point), `No nearest intersection point found`
	}

	if s2, err := strconv.ParseFloat(respBody.Intersection["lng"], 64); err == nil {
		resLng = s2
	} else {
		return *new(Point), `No nearest intersection point found`
	}

	return NewPoint(resLat, resLng), ``


}

func getTrails(org Point) ([]TrailInfo, string) {
	url := TRAILS_URL + strconv.FormatFloat(org.lat, 'f', 6, 64) + "&lon=" + strconv.FormatFloat(org.lng, 'f', 6, 64) + "&maxDistance=10&maxResults=8" + TRPKEY
	response, err := apiRequest(url)
	var trailResult []TrailInfo
	if (err != ``) {
		return trailResult, err
	}
	var respBody TrailsResp
	json.Unmarshal(response, &respBody)

	if respBody.Success == 0 {
		return trailResult, `No trails found`
	}
	//fmt.Println(respBody)

	for _, trail := range respBody.Trails {
		temp := new(TrailInfo)
		temp.Name = trail.Name
		temp.Location = trail.Location
		temp.Length = trail.Distance
		
		if strings.Contains(trail.Summary, "summary") {
			temp.Summary = ``
		} else {
			temp.Summary = trail.Summary
		}
		// fmt.Println(trail.Summary)
		distFromOrg, err :=  distance(org.lat, org.lng, trail.Lat, trail.Lon)
		if err != `` {
			return *new([]TrailInfo), `distance error`
		}
		temp.DistFromOrg = float64(distFromOrg) * M_TO_MI
		temp.Coords = []float64{trail.Lat, trail.Lon}
		trailResult = append(trailResult, *temp)
	}
	if len(trailResult) == 0 {
		return trailResult, `no trails found`
	}
	sort.SliceStable(trailResult, func(i, j int) bool {
		return trailResult[i].DistFromOrg < trailResult[j].DistFromOrg
	})
	fmt.Println(trailResult)
	return trailResult, ``
}

// given an origin, distance and angle, finds the corresponding point
func getPoint(point Point, distance float64, angle float64, runNearestIntersection bool) (Point, string){
	//angle is degrees
	radians := angle * math.Pi / 180
	distanceLat := 1 / (69 / distance)
	oneDegreeLng := math.Cos(point.lat*math.Pi/180) * EQUATOR_LENGTH
	distanceLng := 1 / (oneDegreeLng / distance)

	lat := point.lat + math.Cos(radians)*distanceLat
	lng := point.lng + math.Sin(radians)*distanceLng
	
	if (!runNearestIntersection) {
		return NewPoint(lat, lng), ``
	}

    return nearestIntersectionPoint(NewPoint(lat, lng))
}

//Given an origin, desired distance, number of different routes it wants and a function that
//determines the shape of the route, returns a list of possible routes
func createRoutes(point Point, distance float64, numb float64, routeFunction makeRoute) [][]Point {
	fmt.Println("Starting to create routes")
	angleIncrease := 360 / numb
	num := int(numb)
	var routes [][]Point
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func(i int) {
			defer wg.Done()
			route := routeFunction(point, distance, float64(i) * angleIncrease)
			if (len(route) > 0) {
				routes = append(routes, route)
			}
		} (i)
	}
	wg.Wait()
	fmt.Println("Finished creating routes")
	return routes
}

//creates a straight line route
var straightLine makeRoute = func(point Point, distance float64, offset float64) []Point {

	p0, err := getPoint(point, distance/2, offset, true)

	if (err != ``) {
		return []Point{}
	}

	return []Point{p0}
}

//creates a square route
var squareRoute makeRoute = func(point Point, distance float64, offset float64) []Point {

	sideLength := (distance / 4)

	p0, err0 := getPoint(point, sideLength, offset+135.0, true)
	p1, err1 := getPoint(point, distance/4*math.Sqrt(2), offset+90.0, true)
	p2, err2 := getPoint(point, sideLength, offset+45.0, true)

	if (err0 != `` || err1 != `` || err2 != ``) {
		return []Point{}
	}

	return []Point{p0, p1, p2}
}

var equilateralTriangle makeRoute = func(point Point, distance float64, offset float64) []Point {
	sideLength := (distance / 3.0)

	p0, err0 := getPoint(point, sideLength, offset+60.0,true)
	p1, err1 := getPoint(point, sideLength, offset+120.0,true)

	if (err0 != `` || err1 != ``) {
		return []Point{}
	}

	return []Point{p0, p1}
}

//isosceles
var rightTriangle makeRoute = func(point Point, distance float64, offset float64) []Point {
	sideLength := distance / (2.0 + math.Sqrt(2.0))

	p0, err0 := getPoint(point, sideLength, offset,true)
	p1, err1 := getPoint(point, sideLength, offset+90.0,true)

	if (err0 != `` || err1 != ``) {
		return []Point{}
	}

	return []Point{p0, p1}
}

//isosceles
var rightTriangleOther makeRoute = func(point Point, distance float64, offset float64) []Point {
	sideLength := distance / (2.0 + math.Sqrt(2.0))
	p0, err0 := getPoint(point, sideLength, offset+90.0,true)
	
	if (err0 != ``) {
		return []Point{}
	}

	p1, err1 := getPoint(p0, sideLength, offset,true)

	if (err1 != ``) {
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
	
	response, err := apiRequest(urlCall)
	if (err != ``) {
		return 0, err
	}

	if !checkResponseDistance(response) {
		return 0, "Response Distance Status not okay"
	}

	var respBody DistanceResp
	json.Unmarshal(response, &respBody)

	return respBody.Rows[0].Elements[0].Distance.Value, ``

}

func distanceHelp(dirURL string) (float64, []LocOfTurn, string) {
	response, err := apiRequest(dirURL)
	var result []LocOfTurn
	if (err != ``) {
		return 0.0, result, err
	}

	var respBody DirResp
	json.Unmarshal(response, &respBody)

	if (respBody.Status == "OK") {
		tempDist := 0
		for _, v := range respBody.Rt[0].Legs[0].Steps {
			turnLocs := new(LocOfTurn)
			turnLocs.Turn = v.Maneuver
			turnLocs.Instructions = v.HtmlInstructions
			temp := []float64{v.LocStep.Lat, v.LocStep.Lng}
			tempEnd := []float64{v.LocStepEnd.Lat, v.LocStepEnd.Lng}
			turnLocs.Loc = temp
			turnLocs.EndLoc = tempEnd
			result = append(result, *turnLocs)
			tempDist += v.Dist.Value
		}
		return float64(tempDist), result, ""
	}
	fmt.Println("Status not okay")
	return 0.0, result, respBody.Status
}

func getDistance(pathSlice [][]Point, org Point, desired float64) ([]DistAndPath, float64) {
	var allDists []DistAndPath
	oLat := strconv.FormatFloat(org.lat, 'f', 6, 64)
	oLng := strconv.FormatFloat(org.lng, 'f', 6, 64)
	sumDist := 0.0
	numPaths := len(pathSlice)
	if numPaths == 0 {
		return allDists, desired
	}
	fmt.Println("Starting to get distances")
	var wg sync.WaitGroup
	wg.Add(numPaths)

	for i := 0; i < numPaths; i++ {
		go func (i int) {
			defer wg.Done()
			tempUrl := "https://maps.googleapis.com/maps/api/directions/json?origin=" + oLat + "," + oLng + "&destination=" + oLat + "," + oLng + MODE + "&waypoints="
			for _, x := range pathSlice[i] {
				cLat := strconv.FormatFloat(x.lat, 'f', 6, 64)
				cLng := strconv.FormatFloat(x.lng, 'f', 6, 64)
				tempUrl = tempUrl + "via:" + cLat + "," + cLng + "|"
			}
			tempUrl = strings.TrimSuffix(tempUrl, "|")
			url := tempUrl + KEY
			dist, turnLocs, err := distanceHelp(url)
			if (err == ``) {
				// convert output in meters to miles
				distMi := dist * M_TO_MI
				sumDist += distMi
				//only includes paths which are at least the desired length and at most desired length + 1 mi
				if ((distMi > desired) && (distMi < (desired + 1))) {
					temp := new(DistAndPath)
					temp.distance = distMi
					temp.path = pathSlice[i]
					temp.desired = desired
					temp.turns = turnLocs
					allDists = append(allDists, *temp)
				}
			}
		} (i)
	}
	wg.Wait()
	fmt.Println("Finished getting distances")
	fmt.Println(sumDist, len(pathSlice))
	return allDists, sumDist / float64(numPaths)
}

func getErrorResponse(error string) *FullResponse {
	this := new(FullResponse)
	this.Results = nil
	this.Error = error
	return this
}

func getTrailErrorResponse(error string) *TrailResponse {
	this := new(TrailResponse)
	this.Results = nil
	this.Origin = nil
	this.Error = error
	return this
}

func newFullResponse(results []ResponseNew) *FullResponse {
	this := new(FullResponse)
	this.Results = results
	this.Error = ``
	return this
}

func getInitFix(org Point) float64 {
	p0, err0 := getPoint(org, .707107, 45.0, true)
	p1, err1 := getPoint(org, .707107, 135.0, true)
	p2, err2 := getPoint(org, .707107, 225.0, true)
	p3, err3 := getPoint(org, .707107, 315.0, true)

	if (err0 != `` || err1 != `` || err2 != `` || err3 != ``) {
		return 1.0
	}
	var ptRoute [][]Point
	ptRoute = append(ptRoute, []Point{p1, p2, p3})
	_, avg := getDistance(ptRoute, p0, 4.0)
	return 4.0 / avg
}

//given an address, distance and route, finds a path
func executeRequest(input string, distanceString string) *FullResponse {
	var routeOrder []makeRoute
	routeOrder = append(routeOrder, equilateralTriangle, equilateralTriangle, rightTriangle, rightTriangle, rightTriangleOther, rightTriangleOther, squareRoute, squareRoute, straightLine, straightLine)
	//convert distance to float64
	//check to see if distance is a proper number
	distance, err := strconv.ParseFloat(distanceString, 64)
	if err != nil {
		return getErrorResponse("Not a valid distance")
	}

	//form url from address
	url := addressToAPICall(input)

	//get response from google api server
	response, error := apiRequest(url)

	if (error != ``) {
		return getErrorResponse(error)
	}

	//check to see if address exists
	if !checkResponseGeocode(response) {
		return getErrorResponse("Address does not exist")
	}

	//get the latitude and longitude
	lat, lng := extractCoordinates(response)
	origin := NewPoint(lat, lng)
	fmt.Println(origin)

	initFix := getInitFix(origin)
	fmt.Println(initFix)
	//get the possible routes
	routes := createRoutes(origin, distance*(initFix), 8.0, routeOrder[0])

	desired := distance
	pathDetails, avgDist := getDistance(routes, origin, float64(desired))
	fmt.Println(avgDist)
	errorFix := desired / avgDist
	fmt.Println(errorFix)
	attempts := 0
	for len(pathDetails) < 8 {
		attempts += 1
		if attempts > 9 {
			fmt.Println(`At least we tried`)
			break
		}
		fmt.Println(len(pathDetails), "returnable routes")
		fmt.Println("attempt", attempts)
		if attempts % 2 == 0 {
			errorFix = 1.0
		}
		fmt.Println("with error fix of", errorFix)
		routes = createRoutes(origin, distance*(errorFix)*(initFix), 8.0, routeOrder[attempts])
		morePaths, avgDist := getDistance(routes, origin, float64(desired))
		fmt.Println(avgDist, desired)
		errorFix = desired / avgDist
		pathDetails = append(pathDetails, morePaths...)

	}
	//checks to make sure there are paths found
	if len(pathDetails) == 0 {
		return getErrorResponse("No paths found")
	}
	fmt.Println(len(pathDetails))
	
	var wg sync.WaitGroup
	wg.Add(len(pathDetails))
	var result []ResponseNew
	for i := 0; i < len(pathDetails); i++ {
		go func (i int) {
			defer wg.Done()
			temp := new(ResponseNew)
			temp.Org = append(temp.Org, origin.lat, origin.lng)
			temp.Dest = append(temp.Dest, origin.lat, origin.lng)
			for _, pt := range pathDetails[i].path {
				temp.Path = append(temp.Path, pt.lat, pt.lng)
			}
			temp.Distance = pathDetails[i].distance
			fmt.Println(temp.Distance)
			temp.Directions = pathDetails[i].turns
			result = append(result, *temp)
		} (i)
	}
	wg.Wait()

	//sorts difference desired distance from path distance from least to greatest
	sort.SliceStable(result, func(i, j int) bool {
		if len(result[i].Path) > 1 && len(result[j].Path) > 1 {
			return math.Abs(result[i].Distance-distance) < math.Abs(result[j].Distance-distance)
		}
		if len(result[i].Path) == 1 && len(result[j].Path) == 1 {
			return math.Abs(result[i].Distance-distance) < math.Abs(result[j].Distance-distance)
		}
		return len(result[i].Path) > len(result[i].Path)
	})

	if len(result) > 8 {
		//only want best 8 if there are more then 8
		result = result[0:8]
	}

	return newFullResponse(result)
}


func Execute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var req Request
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := json.NewEncoder(w).Encode(executeRequest(req.Address, req.Distance)); err != nil {
		json.NewEncoder(w).Encode(getErrorResponse(fmt.Sprintf(`%s`, err)))
	}
}

func ExecuteTrailRequest(address string) *TrailResponse {

	//form url from address
	url := addressToAPICall(address)

	//get response from google api server
	response, error := apiRequest(url)
	if (error != ``) {
		return getTrailErrorResponse(error)
	}

	//check to see if address exists
	if !checkResponseGeocode(response) {
		return getTrailErrorResponse("Address does not exist")
	}

	//get the latitude and longitude
	lat, lng := extractCoordinates(response)
	origin := NewPoint(lat, lng)

	trails, error := getTrails(origin)

	if (error != ``) {
		return getTrailErrorResponse(error)
	} else {
		this := new(TrailResponse)
		this.Results = trails
		this.Origin = []float64{lat, lng}
		this.Error = ``
		return this
	}

}

func ExecuteTrail(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var req Request
	_ = json.NewDecoder(r.Body).Decode(&req)


	TrailResponse := ExecuteTrailRequest(req.Address)

	if err := json.NewEncoder(w).Encode(TrailResponse); err != nil {
		json.NewEncoder(w).Encode(getTrailErrorResponse(fmt.Sprintf(`%s`, err)))
	}
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


	stop1 := new(LocOfTurn)
	stop1.Instructions = "Head <b>south</b> on <b>Palmtag Dr</b> toward <b>Bellwood Dr</b>"
	stop1.Loc = []float64{37.2862401,  -122.0078964}
	stop1.EndLoc = []float64{37.2862401,  -122.0078964}
	stop1.Turn =  ""

	stop2 := new(LocOfTurn)
	stop2.Instructions = "Turn <b>right</b> onto <b>Bellwood Dr</b>"
	stop2.Loc = []float64{37.2856419, -122.0080037}
	stop2.EndLoc = []float64{37.2862401,  -122.0078964}
	stop2.Turn =  "turn-right"

	stop3 := new(LocOfTurn)
	stop3.Instructions = "Turn <b>left</b> onto <b>Titus Ave</b>"
	stop3.Loc = []float64{37.2857637, -122.0100177}
	stop3.EndLoc = []float64{37.2862401,  -122.0078964}
	stop3.Turn =  "turn-left"

	temp.Directions = []LocOfTurn{*stop1, *stop2, *stop3}



	temp2 := new(ResponseNew)
	temp2.Org = []float64{37.2862852, -122.0081542}
	temp2.Dest = []float64{37.2862852, -122.0081542}
	temp2.Path = []float64{37.2866, -122.01331, 37.28257, -122.01174, 37.28244, -122.00863}
	temp2.Distance = 1.39870387

	temp2.Directions = []LocOfTurn{*stop1, *stop2, *stop3}

	temp3 := new(ResponseNew)
	temp3.Org = []float64{37.2862852, -122.0081542}
	temp3.Dest = []float64{37.2862852, -122.0081542}
	temp3.Path = []float64{37.28352, -122.01188, 37.28181, -122.00754, 37.2837, -122.00532}
	temp3.Distance = 1.39870387

	temp3.Directions = []LocOfTurn{*stop1, *stop2, *stop3}

	response := newFullResponse([]ResponseNew{*temp, *temp2, *temp3})
	fmt.Println(response)

	json.NewEncoder(w).Encode(response)

}
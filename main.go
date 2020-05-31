package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "math"
 )

const KEY = "&key=AIzaSyB32cCcL4gD_WIYPP6dAVSprY_QYE3arsk"
const URL = "https://maps.googleapis.com/maps/api/geocode/json?address="
const EQUATOR_LENGTH = 69.172

func address_to_api_call (address string) string {
    //properly form the address for the api url
    arr := strings.Split(address," ")
    var address_url string
    for _, s := range arr {
        address_url += "+" + s
    }
    address_url = strings.TrimPrefix(address_url, "+")

    //create the api url and get the json response
    url := URL + address_url + KEY
    return url
}

func api_request (url string) []byte {

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

func check_response (response []byte) bool {
    var resp_body interface{}
    json.Unmarshal(response,&resp_body)
    status := resp_body.(map[string]interface{})["status"]
    if (status == "OK") {return true}
    return false

}

func extract_coordinates (response []byte) (float64,float64) {
    var resp_body interface{}
    json.Unmarshal(response,&resp_body)
    results := resp_body.(map[string]interface{})["results"]
    geometry := results.([]interface{})[0].(map[string]interface{})["geometry"]
    location := geometry.(map[string]interface{})["location"]
    coordinates := location.(map[string]interface{})
    
    return coordinates["lat"].(float64), coordinates["lng"].(float64)
}

func possible_routes (lat float64, lng float64, distance float64) [8][2]float64{
    distance_lat := 1 / (69 / distance)
    one_degree_lng := math.Cos(lat * math.Pi/180) * EQUATOR_LENGTH
    distance_lng := 1 / (one_degree_lng / distance)
    root2_lat :=  math.Sqrt(2) * distance_lat
    root2_lng :=  math.Sqrt(2) * distance_lng

    p0 := [2]float64{lat + distance_lat, lng}
    p1 := [2]float64{lat + root2_lat, lng + root2_lng}
    p2 := [2]float64{lat, lng + distance_lng}
    p3 := [2]float64{lat - root2_lat, lng + root2_lng}
    p4 := [2]float64{lat - distance_lat, lng}
    p5 := [2]float64{lat - root2_lat, lng - root2_lng}
    p6 := [2]float64{lat, lng - distance_lng}
    p7 := [2]float64{lat + root2_lat, lng - root2_lng}

    return [8][2]float64{p0, p1, p2, p3, p4, p5, p6, p7}

}

func main() {

    //setup the scanner
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Printf("Enter your address in the following format: Street Address, City, State\n")

    //read user input
    var input string
    for scanner.Scan() {
        input = scanner.Text()
        if strings.Contains(input,"") { break }
    }

    //form url from address
    url := address_to_api_call(input)

    //get response from google api server
    response := api_request(url)

    //check to see if address exists
    if (!check_response(response)) {
        fmt.Println("Address doesn't exist")
        return
    }

    //get the latitude and longitude
    lat, lng := extract_coordinates(response)

    //get the possible routes
    routes := possible_routes(lat, lng, 0.5)

    fmt.Println(routes)

}
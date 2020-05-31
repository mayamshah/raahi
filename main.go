package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
    "net/http"
    "io/ioutil"
    "encoding/json"
 )

const KEY = "&key=AIzaSyB32cCcL4gD_WIYPP6dAVSprY_QYE3arsk"
const URL = "https://maps.googleapis.com/maps/api/geocode/json?address="

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

    fmt.Println(lat)
    fmt.Println(lng)


}
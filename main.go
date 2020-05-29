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

    //properly form the address for the api url
    arr := strings.Split(input," ")
    var address string
    for _, s := range arr {
        address += "+" + s
    }
    address = strings.TrimPrefix(address, "+")

    //create the api url and get the json response
    url := URL + address + KEY
    resp, err := http.Get(url)
    if err != nil {
        panic(err)
    }
    response, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close() 


    //get the latitude and longitude
    var f interface{}
    json.Unmarshal(response, &f)   
    m := f.(map[string]interface{})
    f = m["results"]
    n := f.([]interface{})
    o := n[0]
    g := o.(map[string]interface{})
    f = g["geometry"]
    j := f.(map[string]interface{})
    f = j["location"]
    i := f.(map[string]interface{})
    fmt.Println(i["lat"])
    fmt.Println(i["lng"])



}
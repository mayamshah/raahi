package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
 )

const KEY = "AIzaSyB32cCcL4gD_WIYPP6dAVSprY_QYE3arsk"

func main() {

    str := "https://maps.googleapis.com/maps/api/geocode/json?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA&key="
    str += KEY
    resp, err := http.Get(str)
    if err != nil {
    	}
    response, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s",response)
}
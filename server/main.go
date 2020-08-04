package main

import (
    "fmt"
    "log"
    "net/http"
    router "m20project.com/m/router"
)

func main() {

	//require python server to be working 
	// resp, err := http.Get("http://localhost:5000/token")
	// if (err != nil) {
	// 	fmt.Println("Please start the python strava token server")
	// 	panic(err)
	// }

	// if (resp.StatusCode != 200) {
	// 	fmt.Println("Make sure to visit localhost:5000/authorize before starting the backend server")
	// 	return
	// }


	r := router.Router()
    fmt.Println("Starting server on the port 8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
}


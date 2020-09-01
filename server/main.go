package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "bufio"
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

	fmt.Println("Enter in Google API key")
	reader := bufio.NewReader(os.Stdin)
	google_key, _ := reader.ReadString('\n')

	fmt.Println("Enter in Trail Run Project key")
	trails_key, _ := reader.ReadString('\n')

	r := router.Router(google_key, trails_key)
    fmt.Println("Starting server on the port 8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
}


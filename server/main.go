package main

import (
    "fmt"
    "log"
    "net/http"
    router "m20project.com/m/router"
    "flag"
    "errors"
)

func main() {
    var accessToken string
	flag.StringVar(&accessToken, "token", ``, "Access Token")
	flag.Parse()
	if (accessToken == ``) {
		panic(errors.New("no access token set"))
	}

	r := router.Router(accessToken)
    fmt.Println("Starting server on the port 8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
}
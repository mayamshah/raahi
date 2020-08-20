package router 

import (
	app "m20project.com/m/app"
	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/api/execute", app.Execute).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/executestrava", app.ExecuteStrava).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/executetrail", app.ExecuteTrail).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/tester", app.Tester).Methods("POST", "OPTIONS")

	return router
}
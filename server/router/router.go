package router 

import (
	app "m20project.com/m/app"
	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router(token string) *mux.Router {

	router := mux.NewRouter()

	tool := app.NewTool(token)

	router.HandleFunc("/api/execute", app.Execute).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/executestrava", tool.ExecuteStrava).Methods("POST", "OPTIONS")
	return router
}
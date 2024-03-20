package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// testHandler responds with a JSON object
func testHandler2(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"test": "ok"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/test", testHandler2).Methods("GET")
	http.ListenAndServe(":3000", router)
}

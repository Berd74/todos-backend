package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Check if the request method is GET
	if r.Method == "GET" {
		// Create a response structure
		response := map[string]string{
			"test": "ok",
		}

		// Convert the response structure to JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			// Handle the error in case JSON marshaling fails
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the JSON response
		w.Write(jsonResponse)
	} else {
		// Respond with an error if the method is not GET
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {

	// Register the handler function for the /test route
	http.HandleFunc("/test", testHandler)

	fmt.Printf("server started")

	// Start the HTTP server on port 3000 and listen for requests
	http.ListenAndServe(":3000", nil)

}

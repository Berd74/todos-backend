package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func SendRequest(httpMethod string, url string, body *map[string]any) *httptest.ResponseRecorder {
	var jsonBody []byte
	if body != nil {
		json, err := json.Marshal(body)
		if err != nil {
			panic("something went wrong with parsing json")
		}
		jsonBody = json
	}
	req, _ := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	Router.ServeHTTP(w, req)
	return w
}

func Map[T any, U any](s []T, f func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

func MustUnmarshal[T any](response *httptest.ResponseRecorder) T {
	var body T
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		panic("Error parsing JSON response")
	}
	return body
}

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"firebase.google.com/go/v4/auth"
	"fmt"
	guuid "github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"todoBackend/firebase"
)

type FirebaseAdminSettings struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
	APIKey                  string `json:"api_key"`
}

type TestUser struct {
	IdToken string
	UserId  string
}

func createTestUser() *TestUser {
	ctx := context.Background()
	id := guuid.New()
	email := id.String() + "@testing.user.com"
	password := "yourStrong!Password"
	params := (&auth.UserToCreate{}).
		Email(email).
		EmailVerified(true).
		Password(password)
	newUser, err := firebase.AuthClient.CreateUser(ctx, params)
	if err != nil {
		log.Fatal(err)
	}
	token, err := getIdToken(email, password)
	if err != nil {
		log.Fatal(err)
	}
	testUser := TestUser{token, newUser.UID}
	return &testUser
}

func getApiKey() string {
	// Example of reading JSON from a file and unmarshalling into the struct
	data, err := os.ReadFile("./firebase-adminsdk.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}

	var account FirebaseAdminSettings
	err = json.Unmarshal(data, &account)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return ""
	}

	return account.APIKey
}

func getIdToken(email, password string) (string, error) {

	apiKey := getApiKey()

	data := map[string]string{
		"email":             email,
		"password":          password,
		"returnSecureToken": "true",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=" + apiKey
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		IdToken string `json:"idToken"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result.IdToken, nil
}

//todo
func deleteTestUser(uid string) error {
	ctx := context.Background()
	return firebase.AuthClient.DeleteUser(ctx, uid)
}

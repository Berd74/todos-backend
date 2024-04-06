package tests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"testing"
	"todoBackend/database"
	"todoBackend/firebase"
	"todoBackend/model"
	"todoBackend/routes"
)

const token = `eyJhbGciOiJSUzI1NiIsImtpZCI6IjgwNzhkMGViNzdhMjdlNGUxMGMzMTFmZTcxZDgwM2I5MmY3NjYwZGYiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoibW9nb2dhbWluZyBnYW1pbmciLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EvQUNnOG9jSmd6Y1Q1N3Q3Z1FvVlpJSkFzWkNRXzdsaVZoekpzbWVhSlFLUkNpbmRVNXMtUVFRPXM5Ni1jIiwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL3RvZG9zLTY0NWY0IiwiYXVkIjoidG9kb3MtNjQ1ZjQiLCJhdXRoX3RpbWUiOjE3MTIzMTY2NjAsInVzZXJfaWQiOiJsUkJsNlBEY3QzWExBTmxKa2ZNRnJleXZ5WUUzIiwic3ViIjoibFJCbDZQRGN0M1hMQU5sSmtmTUZyZXl2eVlFMyIsImlhdCI6MTcxMjM2MjExMSwiZXhwIjoxNzEyMzY1NzExLCJlbWFpbCI6Im1vZ29nYW1pbmc3NUBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJnb29nbGUuY29tIjpbIjEwNjQ0ODY3MTMzNjk0NTEyMTI2MyJdLCJlbWFpbCI6WyJtb2dvZ2FtaW5nNzVAZ21haWwuY29tIl19LCJzaWduX2luX3Byb3ZpZGVyIjoiZ29vZ2xlLmNvbSJ9fQ.VddK1y0oi9unNr3yAB5zdMdaMiHxW63h3eP97mRIsGV_7Ev5kchGysQWeBN__CuYGPo5Pi1wgQSIbCYP8vp7RfDt68W3VtV25Rc4KE18Im3P7uaqDGmjvxpwb9795expIFQT5gVI5BvNibKKATB64tsjYj2EyfM49YZ-3CiLzxFhoNWoahr4Am1pnLe3-TtMSDAWj-C2VIPlzwmAHWOD0sOJQ7orh5c9Kvalq7AeHqpBlBDSTXLqQzFI5baOl0YUwJTQUvb-9pa1gTUwQtFDYnUU4fVphaIOcxPinihChv2pqdTaZZZNwOEbneclUi3_vUaMkzKnSBwnHPaeC1l7sw`
const userId = `rNKQUxPs-qvcd-hshC-E4M1-WPoyalu2todo`

var Router *gin.Engine

func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		fmt.Println("Error changing directory:", err)
		return
	}
	firebase.InitFirebase()
	database.InitDatabase()
	setupRouter()
	prepareUser()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestApi(t *testing.T) {

	t.Run("collection creation with different parameters", func(t *testing.T) {
		const name = "Collection 1"

		payload := map[string]any{"name": name}

		response := SendRequest(http.MethodPost, "/collection/", &payload)

		body := MustUnmarshal[struct {
			Data model.Collection `json:"data"`
		}](response)

		if body.Data.Name != name {
			t.Errorf("Wrong name: got %v want %v", body.Data.Name, name)
		}
	})

	t.Run("collection creation with different parameters", func(t *testing.T) {
		const name = "Collection 2"
		const desc = "Desc"

		payload := map[string]any{"name": name, "description": desc}

		response := SendRequest(http.MethodPost, "/collection/", &payload)

		body := MustUnmarshal[struct {
			Data model.Collection `json:"data"`
		}](response)

		if body.Data.Name != name {
			t.Errorf("Wrong name: got %v want %v", body.Data.Name, name)
		}
		if *body.Data.Description != desc {
			t.Errorf("Wrong name: got %v want %v", body.Data.Description, desc)
		}
	})

	t.Run("incorrect name", func(t *testing.T) {
		const name = "X"

		payload := map[string]any{"name": name}

		response := SendRequest(http.MethodPost, "/collection/", &payload)

		if response.Code != http.StatusBadRequest {
			t.Errorf("Wrong status: got %v want %v", response.Code, http.StatusBadRequest)
		}
	})

}

func setupRouter() {
	gin.SetMode(gin.TestMode)
	Router = gin.Default()
	routes.Todo(Router.Group("/todo"))
	routes.Collection(Router.Group("/collection"))
}

func prepareUser() {
	// get and remove all existing collections and tasks

	response := SendRequest(http.MethodGet, fmt.Sprintf("/collection/?userIds=%v", userId), nil)
	if s := response.Code; s != http.StatusOK {
		panic("setup error - make sure that token & userId are correct")
	}

	body := MustUnmarshal[struct {
		Data []model.Collection `json:"data"`
	}](response)

	collectionsIds := Map(body.Data, func(c model.Collection) string {
		return c.CollectionId
	})
	if len(collectionsIds) != 0 {
		response = SendRequest(http.MethodDelete, fmt.Sprintf("/collection/?ids=%v", strings.Join(collectionsIds, ",")), nil)
		if s := response.Code; s != http.StatusOK {
			panic("setup error - make sure that token & userId are correct")
		}
	}

}

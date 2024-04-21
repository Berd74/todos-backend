package tests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"testing"
	"todoBackend/database"
	"todoBackend/firebase"
	"todoBackend/model"
	"todoBackend/routes"
	"todoBackend/utils"
)

var Token *string
var UserId *string
var Router *gin.Engine

func TestMain(m *testing.M) {
	os.Setenv("SPANNER_EMULATOR_HOST", "localhost:9010")
	err := os.Chdir("..")
	if err != nil {
		fmt.Println("Error changing directory:", err)
		return
	}
	firebase.InitFirebase()
	database.InitDatabase()
	setupRouter()
	user := createTestUser()
	Token = &user.IdToken
	userIdString := utils.GoogleIdToUuid(user.UserId)
	UserId = &userIdString
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

}

func setupRouter() {
	gin.SetMode(gin.TestMode)
	Router = gin.Default()
	routes.Todo(Router.Group("/todo"))
	routes.Collection(Router.Group("/collection"))
}

//func prepareUser() {
//	// get and remove all existing collections and tasks
//
//	response := SendRequest(http.MethodGet, fmt.Sprintf("/collection/?userIds=%v", UserId), nil)
//	if s := response.Code; s != http.StatusOK {
//		panic("setup error - make sure that token & userId are correct")
//	}
//
//	body := MustUnmarshal[struct {
//		Data []model.Collection `json:"data"`
//	}](response)
//
//	collectionsIds := Map(body.Data, func(c model.Collection) string {
//		return c.CollectionId
//	})
//	if len(collectionsIds) != 0 {
//		response = SendRequest(http.MethodDelete, fmt.Sprintf("/collection/?ids=%v", strings.Join(collectionsIds, ",")), nil)
//		if s := response.Code; s != http.StatusOK {
//			panic("setup error - make sure that token & userId are correct")
//		}
//	}
//
//}

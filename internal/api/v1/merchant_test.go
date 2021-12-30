package v1

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abhinavmsra/go-api/internal/api"
	"github.com/abhinavmsra/go-api/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(cors.Default())

	return router
}

func getStorage() repository.Storage {
	connectionString := "postgres://api@db:5432/api_test?sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	storage := repository.NewStorage(db)
	if err = storage.RunMigrations(connectionString); err != nil {
		log.Println(err)
	}
	// truncate dataabase for fresh start
	db.Query(`TRUNCATE TABLE admins, merchants RESTART IDENTITY CASCADE;`)

	adminSecret, err := api.GenerateSecret()
	if err != nil {
		panic(err)
	}

	storage.CreateAdmin("god", adminSecret)

	return storage
}

func getAdmin(storage repository.Storage) *api.Admin {
	var admin = new(api.Admin)
	row := storage.FindAdminByName("god")

	if row.Err() != nil {
		panic("No Admin")
	}

	if err := row.Scan(&admin.Id, &admin.CreatedAt, &admin.UpdatedAt, &admin.Name, &admin.Secret); err != nil {
		panic("No Admin")
	}

	return admin
}

func TestIndexMerchant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	endpoint := "/api/v1/merchants"

	r := setupRouter()
	storage := getStorage()
	r.GET(endpoint, IndexMerchant(storage))

	admin := getAdmin(storage)

	data := make(map[string]int)
	data["Bearer XXXX"] = http.StatusForbidden
	data[fmt.Sprintf("Bearer %s", admin.Secret)] = http.StatusOK

	for token, statusCode := range data {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			t.Fatalf("Couldn't create request: %v\n", err)
		}

		req.Header.Add("Authorization", token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != statusCode {
			t.Fatalf("Expected to get status %d but instead got %d\n", statusCode, w.Code)
		}
	}
}

func TestCreateMerchant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	endpoint := "/api/v1/merchants"

	r := setupRouter()
	storage := getStorage()
	admin := getAdmin(storage)

	r.POST(endpoint, CreateMerchant(storage))

	merchant := api.Merchant{
		Name: "New Merchant",
	}

	var buffer bytes.Buffer
	if err := jsonapi.MarshalPayload(&buffer, &merchant); err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(buffer.Bytes()))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", admin.Secret))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusCreated, w.Code)
	}
}

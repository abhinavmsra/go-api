package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/abhinavmsra/go-api/internal/api"
	"github.com/abhinavmsra/go-api/internal/repository"
)

func getAdminByAuthHeader(storage repository.Storage, auth string) (*api.Admin, error) {
	var admin = new(api.Admin)

	parts := strings.Split(auth, " ")
	row := storage.FindAdminBySecret(parts[1])

	if row.Err() != nil {
		return admin, errors.New("Invalid Bearer Token")
	}

	if err := row.Scan(&admin.Id, &admin.CreatedAt, &admin.UpdatedAt, &admin.Name, &admin.Secret); err != nil {
		return admin, errors.New("Invalid Bearer Token")
	}

	return admin, nil
}

func IndexMerchant(storage repository.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		_, err := getAdminByAuthHeader(storage, c.Request.Header.Get("Authorization"))
		if err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusForbidden, "")
			return
		}

		var page int
		var perPage int

		if page, err = strconv.Atoi(c.Query("page")); err != nil {
			page = 1
		}

		if perPage, err = strconv.Atoi(c.Query("per_page")); err != nil {
			perPage = 25
		}

		rows, err := storage.IndexMerchant(perPage, (page-1)*perPage)
		if err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		merchants := []*api.Merchant{}

		for rows.Next() {
			var merchant = new(api.Merchant)
			if err := rows.Scan(&merchant.Id, &merchant.CreatedAt, &merchant.UpdatedAt, &merchant.Name, &merchant.Secret); err != nil {
				log.Printf("[ERROR]: %v", err)
				c.JSON(http.StatusInternalServerError, err)
				return
			}

			merchants = append(merchants, merchant)
		}

		var w strings.Builder
		if err = jsonapi.MarshalPayload(&w, merchants); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		c.String(http.StatusOK, w.String())
	}
}

func CreateMerchant(storage repository.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		_, err := getAdminByAuthHeader(storage, c.Request.Header.Get("Authorization"))
		if err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusForbidden, "")
			return
		}

		var merchant = new(api.Merchant)
		if err = jsonapi.UnmarshalPayload(c.Request.Body, merchant); err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		newSecret, err := api.GenerateSecret()
		if err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		row := storage.CreateMerchant(merchant.Name, newSecret)
		if row.Err() != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		row.Scan(&merchant.Id, &merchant.CreatedAt, &merchant.UpdatedAt, &merchant.Name, &merchant.Secret)

		var w strings.Builder
		if err = jsonapi.MarshalPayload(&w, merchant); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		c.String(http.StatusCreated, w.String())
	}
}

func ShowMerchant(storage repository.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		_, err := getAdminByAuthHeader(storage, c.Request.Header.Get("Authorization"))
		if err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusForbidden, "")
			return
		}

		var merchant = new(api.Merchant)

		c.Header("Content-Type", "application/json")
		id := c.Param("id")

		row := storage.ShowMerchant(id)
		if row.Err() != nil {
			log.Printf("[ERROR]: %v", row.Err())
			c.JSON(http.StatusNotFound, row.Err())
			return
		}

		if err = row.Scan(&merchant.Id, &merchant.CreatedAt, &merchant.UpdatedAt, &merchant.Name, &merchant.Secret); err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusNotFound, err)
			return
		}

		var w strings.Builder
		if err = jsonapi.MarshalPayload(&w, merchant); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		c.String(http.StatusOK, w.String())
	}
}

func UpdateMerchant(storage repository.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		merchant, err := GetMerchantByAuthHeader(storage, c.Request.Header.Get("Authorization"))
		id := strconv.Itoa(merchant.Id)
		if err != nil || id != c.Param("id") {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusForbidden, "")
			return
		}

		var updatedMerchant = new(api.Merchant)
		if err = jsonapi.UnmarshalPayload(c.Request.Body, updatedMerchant); err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		row := storage.UpdateMerchant(id, updatedMerchant.Name)
		if row.Err() != nil {
			log.Printf("[ERROR]: %v", row.Err())
			c.JSON(http.StatusBadRequest, row.Err())
			return
		}

		row = storage.ShowMerchant(id)
		if err = row.Scan(&merchant.Id, &merchant.CreatedAt, &merchant.UpdatedAt, &merchant.Name, &merchant.Secret); err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusNotFound, err)
			return
		}

		var w strings.Builder
		if err = jsonapi.MarshalPayload(&w, merchant); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		c.String(http.StatusOK, w.String())
	}
}

func DeleteMerchant(storage repository.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		merchant, err := GetMerchantByAuthHeader(storage, c.Request.Header.Get("Authorization"))
		id := strconv.Itoa(merchant.Id)
		if err != nil || id != c.Param("id") {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusForbidden, "")
			return
		}

		if row := storage.DeleteMerchant(id); row.Err() != nil {
			log.Printf("[ERROR]: %v", row.Err())
			c.JSON(http.StatusNotFound, row.Err())
			return
		}

		c.String(http.StatusNoContent, "")
	}
}

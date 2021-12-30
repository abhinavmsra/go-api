package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"log"
	"net/http"
	"strings"

	"github.com/abhinavmsra/go-api/internal/api"
	"github.com/abhinavmsra/go-api/internal/repository"
)

func IndexMerchant(storage repository.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		rows, err := storage.IndexMerchant()
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

		var merchant = new(api.Merchant)
		var err error

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
		var merchant = new(api.Merchant)
		var err error

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
		var merchant = new(api.Merchant)
		var err error

		c.Header("Content-Type", "application/json")
		id := c.Param("id")

		if err = jsonapi.UnmarshalPayload(c.Request.Body, merchant); err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		row := storage.UpdateMerchant(id, merchant.Name)
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

		if row := storage.DeleteMerchant(c.Param("id")); row.Err() != nil {
			log.Printf("[ERROR]: %v", row.Err())
			c.JSON(http.StatusNotFound, row.Err())
			return
		}

		c.String(http.StatusNoContent, "")
	}
}

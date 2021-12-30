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

func GetMerchantByAuthHeader(storage repository.Storage, auth string) (*api.Merchant, error) {
	var merchant = new(api.Merchant)

	parts := strings.Split(auth, " ")
	row := storage.FindMerchantBySecret(parts[1])

	if row.Err() != nil {
		return merchant, errors.New("Invalid Bearer Token")
	}

	if err := row.Scan(&merchant.Id, &merchant.CreatedAt, &merchant.UpdatedAt, &merchant.Name, &merchant.Secret); err != nil {
		return merchant, errors.New("Invalid Bearer Token")
	}

	return merchant, nil
}

func IndexMember(storage repository.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		merchant, err := GetMerchantByAuthHeader(storage, c.Request.Header.Get("Authorization"))
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

		rows, err := storage.IndexMember(merchant.Id, perPage, (page-1)*perPage)
		if err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		members := []*api.Member{}

		for rows.Next() {
			var member = new(api.Member)
			if err := rows.Scan(&member.Id, &member.CreatedAt, &member.UpdatedAt, &member.Name, &member.Secret, &member.Email, &member.MerchantId); err != nil {
				log.Printf("[ERROR]: %v", err)
				c.JSON(http.StatusInternalServerError, err)
				return
			}

			members = append(members, member)
		}

		var w strings.Builder
		if err = jsonapi.MarshalPayload(&w, members); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		c.String(http.StatusOK, w.String())
	}
}

func CreateMember(storage repository.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		merchant, err := GetMerchantByAuthHeader(storage, c.Request.Header.Get("Authorization"))
		if err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusForbidden, "")
			return
		}

		var member = new(api.Member)
		if err = jsonapi.UnmarshalPayload(c.Request.Body, member); err != nil {
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

		row := storage.CreateMember(member.Name, newSecret, member.Email, merchant.Id)
		if row.Err() != nil {
			log.Printf("[ERROR]: %v", row.Err())
			c.JSON(http.StatusBadRequest, row.Err())
			return
		}

		row.Scan(&member.Id, &member.CreatedAt, &member.UpdatedAt, &member.Name, &member.Secret, &member.Email, &member.MerchantId)

		var w strings.Builder
		if err = jsonapi.MarshalPayload(&w, member); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		c.String(http.StatusCreated, w.String())
	}
}

func ShowMember(storage repository.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		merchant, err := GetMerchantByAuthHeader(storage, c.Request.Header.Get("Authorization"))
		if err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusForbidden, "")
			return
		}

		var member = new(api.Member)
		id := c.Param("id")

		row := storage.ShowMember(id, merchant.Id)
		if row.Err() != nil {
			log.Printf("[ERROR]: %v", row.Err())
			c.JSON(http.StatusNotFound, row.Err())
			return
		}

		if err = row.Scan(&member.Id, &member.CreatedAt, &member.UpdatedAt, &member.Name, &member.Secret, &member.Email, &member.MerchantId); err != nil {
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

func UpdateMember(storage repository.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		merchant, err := GetMerchantByAuthHeader(storage, c.Request.Header.Get("Authorization"))
		if err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusForbidden, "")
			return
		}

		var member = new(api.Member)
		id := c.Param("id")

		if err = jsonapi.UnmarshalPayload(c.Request.Body, member); err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		row := storage.UpdateMember(id, merchant.Id, member.Name)
		if row.Err() != nil {
			log.Printf("[ERROR]: %v", row.Err())
			c.JSON(http.StatusBadRequest, row.Err())
			return
		}

		row = storage.ShowMember(id, merchant.Id)
		if err = row.Scan(&member.Id, &member.CreatedAt, &member.UpdatedAt, &member.Name, &member.Secret, &member.Email, &member.MerchantId); err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusNotFound, err)
			return
		}

		var w strings.Builder
		if err = jsonapi.MarshalPayload(&w, member); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		c.String(http.StatusOK, w.String())
	}
}

func DeleteMember(storage repository.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		merchant, err := GetMerchantByAuthHeader(storage, c.Request.Header.Get("Authorization"))
		if err != nil {
			log.Printf("[ERROR]: %v", err)
			c.JSON(http.StatusForbidden, "")
			return
		}

		if row := storage.DeleteMember(c.Param("id"), merchant.Id); row.Err() != nil {
			log.Printf("[ERROR]: %v", row.Err())
			c.JSON(http.StatusNotFound, row.Err())
			return
		}

		c.String(http.StatusNoContent, "")
	}
}

package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mathehluiz/plant-care-tracker/domain"
	"github.com/mathehluiz/plant-care-tracker/internal/errs"
	"github.com/mathehluiz/plant-care-tracker/pkg/validate"
)

func CreateCare(storer domain.CareStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := struct {
			PlantId   int64     `json:"plantId"`
			NextCare  time.Time `json:"nextCare"`
			Name      string    `json:"name"`
			Notes     string    `json:"notes"`
		}{}
		userId := c.GetString("auth:bearer:id")
		parsedUserId, err := strconv.ParseInt(userId, 10, 64)
		if err != nil {
			fmt.Println(err)
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Println(err)
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}
		validations := validate.Struct(req)
		if len(validations) > 0 {
			fmt.Println(validations)
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		care, err := domain.NewCare(req.PlantId, parsedUserId, req.NextCare, req.Name, req.Notes)
		if err != nil {
			fmt.Println(err)
			DefaultError(c, http.StatusBadRequest, err)
			return
		}

		id, err := storer.CreateCare(c.Request.Context(), care)
		if err != nil {
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func GetPlantCares(storer domain.CareStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}
		plantId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		cares, err := storer.GetPlantCares(c.Request.Context(), plantId)
		if err != nil {
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, cares)
	}
}

func GetCareByID(storer domain.CareStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}
		careId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		care, err := storer.GetCareByID(c.Request.Context(), careId)
		if err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
				return
			}
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, care)
	}
}

func UpdateCare(storer domain.CareStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := struct {
			PlantId   int64     `json:"plantId"`
			LastCare  time.Time `json:"lastCare"`
			NextCare  time.Time `json:"nextCare"`
			Name      string    `json:"name"`
			Notes     string    `json:"notes"`
		}{}
		id := c.Param("id")
		careId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
		}

		userID := c.GetString("auth:bearer:id")
		parsedUserID, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}
		validations := validate.Struct(req)
		if len(validations) > 0 {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		care, err := storer.GetCareByID(c.Request.Context(), careId)
		if err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
				return
			}
			DefaultError(c, http.StatusInternalServerError, err)
		}

		err = care.Update(req.PlantId, parsedUserID, req.LastCare, req.NextCare, req.Name, req.Notes)
		if err != nil {
			if errors.Is(err, errs.ErrInvalidCareName) {
				DefaultError(c, http.StatusBadRequest, err)
				return
			}
			if errors.Is(err, errs.ErrInvalidCareNotes) {
				DefaultError(c, http.StatusBadRequest, err)
				return
			}
		}

		err = storer.UpdateCare(c.Request.Context(), care)
		if err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
				return
			}
			DefaultError(c, http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully updated"})
	}
}

func DeleteCare(storer domain.CareStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}
		careId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		err = storer.DeleteCare(c.Request.Context(), careId)
		if err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
				return
			}
			DefaultError(c, http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
	}
}
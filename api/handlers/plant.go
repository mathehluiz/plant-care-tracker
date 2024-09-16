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

func CreatePlant(storer domain.PlantStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := struct {
			Name            string    `json:"name"`
			Location        string    `json:"location"`
			AcquisitionDate time.Time `json:"acquisitionDate"`
			CareFrequency   int       `json:"careFrequency"`
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
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		plant, err := domain.NewPlant(req.Name, req.Location, req.AcquisitionDate, req.CareFrequency, parsedUserId)
		if err != nil {
			DefaultError(c, http.StatusBadRequest, err)
			return
		}

		id, err := storer.CreatePlant(c.Request.Context(), plant)
		if err != nil {
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func GetPlantByID(storer domain.PlantStorer) gin.HandlerFunc {
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

		plant, err := storer.GetPlantByID(c.Request.Context(), plantId)
		if err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
				return
			}
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, plant)
	}
}

func GetPlantsByUserID(storer domain.PlantStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr := c.GetString("auth:bearer:id")
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			DefaultError(c, http.StatusBadRequest, err)
			return
		}

		plants, err := storer.GetPlantsByUserID(c.Request.Context(), userId)
		if err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
				return
			}
			DefaultError(c, http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, plants)
	}
}

func UpdatePlant(storer domain.PlantStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}
		parsedID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}
		req := struct {
			Name            string    `json:"name"`
			Location        string    `json:"location"`
			AcquisitionDate time.Time `json:"acquisitionDate"`
			CareFrequency   int       `json:"careFrequency"`
		}{}

		if err := c.ShouldBindJSON(&req); err != nil {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}
		validations := validate.Struct(req)
		if len(validations) > 0 {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		plant, err := storer.GetPlantByID(c.Request.Context(), parsedID)
		if err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
				return
			}
			DefaultError(c, http.StatusInternalServerError, err)
		}

		err = plant.Update(req.Name, req.Location, req.AcquisitionDate, req.CareFrequency)
		if err != nil {
			DefaultError(c, http.StatusBadRequest, err)
			return
		}

		err = storer.UpdatePlant(c.Request.Context(), plant)
		if err != nil {
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Plant updated successfully"})
	}
}

func DeletePlant(storer domain.PlantStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}
		parsedID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		err = storer.DeletePlant(c.Request.Context(), parsedID)
		if err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
				return
			}
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Plant deleted successfully"})
	}
}
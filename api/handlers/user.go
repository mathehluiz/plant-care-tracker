package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mathehluiz/plant-care-tracker/domain"
	"github.com/mathehluiz/plant-care-tracker/internal/cache"
	"github.com/mathehluiz/plant-care-tracker/internal/errs"
	"github.com/mathehluiz/plant-care-tracker/internal/random"
	"github.com/mathehluiz/plant-care-tracker/pkg/jwt"
	"github.com/mathehluiz/plant-care-tracker/pkg/mailer"
	"github.com/mathehluiz/plant-care-tracker/pkg/validate"
)

func Login(storer domain.UserStorer, cacher cache.ConnectionStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := struct {
			Email        string `json:"email"`
			Username     string `json:"username"`
			Password     string `json:"password"`
		}{}

		if err := c.ShouldBindJSON(&req); err != nil {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		var user *domain.User
		var err error

		if req.Email != "" {
			user, err = storer.GetUserByEmail(c, req.Email)
			if err != nil {
				if errors.Is(err, errs.ErrSelectNotMatch) {
					DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
					return
				}

				DefaultError(c, http.StatusInternalServerError, err)
				return
			}
		}

		if req.Username != "" && user == nil {
			user, err = storer.GetUserByUsername(c, req.Username)
			if err != nil {
				if errors.Is(err, errs.ErrSelectNotMatch) {
					DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
					return
				}

				DefaultError(c, http.StatusInternalServerError, err)
				return
			}
		}

		if user == nil {
			DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
			return
		}

		if err := user.VerifyPassword(req.Password); err != nil {
			DefaultError(c, http.StatusUnauthorized, errs.ErrInvalidPassword)
			return
		}

		token, err := jwt.GenerateToken(user.Id, user.Roles, user.Verified)
		if err != nil {
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func VerifyCode(storer domain.UserStorer, cacher cache.ConnectionStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := struct {
			Code string `json:"code" validate:"required"`
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

		userId, err := cacher.Get(c, req.Code)
		if err != nil {
			if errors.Is(err, cache.ErrNil) {
				DefaultError(c, http.StatusUnauthorized, errs.ErrInvalidCode)
				return
			}

			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		user, err := storer.GetUserByExternalId(c, userId)
		if err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
				return
			}

			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		_ = cacher.Delete(c, user.ExternalId)

		token, err := jwt.GenerateToken(user.ExternalId, user.Roles, user.Verified)
		if err != nil {
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func RefreshToken(storer domain.UserStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("auth:bearer:id")

		user, err := storer.GetUserByExternalId(c, userId)
		if err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
			}

			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		token, err := jwt.GenerateToken(user.ExternalId, user.Roles, user.Verified)
		if err != nil {
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func GetMe(storer domain.UserStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("auth:bearer:id")

		user, err := storer.GetUserByExternalId(c, userId)
		if err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
			}

			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func RegisterUser(storer domain.UserStorer, cacher cache.ConnectionStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := struct {
			Username string `json:"username" validate:"required,min=4,max=20"`
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required"`
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

		user, err := domain.NewUser(req.Username, req.Email, req.Password, []string{""})
		if err != nil {
			DefaultError(c, http.StatusBadRequest, err)
			return
		}

		externalId, err := storer.CreateUser(c, user)
		if err != nil {
			if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
				if strings.Contains(err.Error(), "users_username_key") {
					DefaultError(c, http.StatusBadRequest, errs.ErrUsernameAlreadyExists)
					return
				}

				DefaultError(c, http.StatusBadRequest, errs.ErrEmailAlreadyExists)
				return
			}

			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		token, err := jwt.GenerateToken(externalId, user.Roles, false)
		if err != nil {
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		code := random.GenetareRandomCode()

		if err := cacher.Set(c, time.Duration(15)*time.Minute, externalId, code); err != nil {
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}
		fmt.Println(code)

		go func() {
			err := mailer.SendConfirmationEmail(req.Email, code)
			if err != nil {
				log.Println("error sending email", err)
				return
			}
		}()

		c.JSON(http.StatusCreated, gin.H{"userId": externalId, "token": token})
	}
}

func VerifyEmail(storer domain.UserStorer, cacher cache.ConnectionStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := struct {
			Code string `json:"code" validate:"required"`
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

		userId := c.GetString("auth:bearer:id")

		verified := c.GetBool("auth:bearer:verified")
		if verified {
			DefaultError(c, http.StatusBadRequest, errs.ErrAlreadyVerified)
			return
		}

		user, err := storer.GetUserByExternalId(c, userId)
		if err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
				return
			}

			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		if user.Verified {
			DefaultError(c, http.StatusBadRequest, errs.ErrAlreadyVerified)
			return
		}

		code, err := cacher.Get(c, user.ExternalId)
		if err != nil {
			if errors.Is(err, cache.ErrNil) {
				code = random.GenetareRandomCode()

				if err := cacher.Set(c, time.Duration(15)*time.Minute, user.ExternalId, code); err != nil {
					DefaultError(c, http.StatusInternalServerError, err)
					return
				}

				go func() {
					err := mailer.SendConfirmationEmail(user.Email, code)
					if err != nil {
						log.Println("error sending email", err)
						return
					}
				}()

				DefaultError(c, http.StatusBadRequest, errs.ErrCodeExpired)
			}

			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		if code != req.Code {
			DefaultError(c, http.StatusUnauthorized, errs.ErrInvalidCode)
			return
		}

		if err := storer.VerifyUser(c, userId); err != nil {
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		_ = cacher.Delete(c, userId)

		token, err := jwt.GenerateToken(user.ExternalId, user.Roles, true)
		if err != nil {
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func SetActive(storer domain.UserStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := struct {
			Active bool `json:"active"`
		}{}

		if err := c.ShouldBindJSON(&req); err != nil {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		userId := c.GetString("auth:bearer:id")

		if err := storer.UpdateActiveUserStatus(c, userId, req.Active); err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
				return
			}

			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func DeleteUser(storer domain.UserStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("id")

		if userId == "" {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		if err := storer.DeleteUser(c, userId); err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
				return
			}

			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func ChangeRoles(storer domain.UserStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := struct {
			Id     string   `json:"id" validate:"required"`
			Roles  []string `json:"roles" validate:"required"`
			Method string   `json:"method" validate:"required,oneof=add remove"`
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

		var err error
		switch req.Method {
		case "add":
			err = storer.AddRolesToUser(c, req.Id, req.Roles)
		case "remove":
			err = storer.RemoveRolesFromUser(c, req.Id, req.Roles)
		default:
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		if err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				DefaultError(c, http.StatusNotFound, errs.ErrNotFound)
				return
			}

			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func ResetPassword(storer domain.UserStorer, cacher cache.ConnectionStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := struct {
			Email string `json:"email" validate:"required,email"`
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

		user, err := storer.GetUserByEmail(c, req.Email)
		if err != nil {
			if errors.Is(err, errs.ErrSelectNotMatch) {
				c.Status(200)
				return
			}

			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		code := uuid.NewString()

		for {
			if v, _ := cacher.Get(c, code); v != "" {
				code = uuid.NewString()
				continue
			}

			break
		}

		if err := cacher.Set(c, time.Duration(5)*time.Minute, code, user.ExternalId); err != nil {
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		err = mailer.SendPasswordResetEmail(user.Email, code)
		if err != nil {
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.Status(200)
	}
}

func ChangePassword(storer domain.UserStorer, cacher cache.ConnectionStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := struct {
			Password string `json:"password" validate:"required"`
		}{}

		if err := c.ShouldBindJSON(&req); err != nil {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		code := c.Param("id")
		if code == "" {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		userId, err := cacher.Get(c, code)
		if err != nil {
			if errors.Is(err, cache.ErrNil) {
				DefaultError(c, http.StatusUnauthorized, errs.ErrInvalidCode)
				return
			}

			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		var user domain.User
		if err := user.HashPass(req.Password); err != nil {
			DefaultError(c, http.StatusBadRequest, err)
			return
		}

		if err := storer.UpdatePassword(c, userId, user.Password); err != nil {
			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		_ = cacher.Delete(c, code)
	}
}

func CheckChangePasswordStatus(cacher cache.ConnectionStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("id")
		if code == "" {
			DefaultError(c, http.StatusBadRequest, errs.ErrInvalidBody)
			return
		}

		userId, err := cacher.Get(c, code)
		if err != nil {
			if errors.Is(err, cache.ErrNil) {
				DefaultError(c, http.StatusUnauthorized, errs.ErrInvalidCode)
				return
			}

			DefaultError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"userId": userId})
	}
}

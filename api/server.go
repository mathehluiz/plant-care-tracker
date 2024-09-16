package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mathehluiz/plant-care-tracker/api/handlers"
	"github.com/mathehluiz/plant-care-tracker/api/middlewares"
	"github.com/mathehluiz/plant-care-tracker/domain"
	"github.com/mathehluiz/plant-care-tracker/internal/cache"
)

type server struct {
	uStorer domain.UserStorer
	pStorer domain.PlantStorer
	cStorer domain.CareStorer
	cacher  cache.ConnectionStorer
}

func NewServer(uStorer domain.UserStorer, pStorer domain.PlantStorer, cStorer domain.CareStorer, cacher cache.ConnectionStorer) server {
	return server{
		uStorer: uStorer,
		pStorer: pStorer,
		cStorer: cStorer,
		cacher:  cacher,
	}
}

func (s server) Start() {
	r := gin.Default()
	s.setupRoles(r)

	log.Fatalln(r.Run(":8080"))
}

func (s server) setupRoles(rg *gin.Engine) {
	keys := []string{"123", "456"}

	bearerMiddleware := middlewares.AddMiddlewares(middlewares.ValidateRoles())
	apiKeyMiddleware := middlewares.AddMiddlewares(middlewares.ValidateAPIKey(keys))

	v1 := rg.Group("/api/v1")

	v1.POST("/login", handlers.Login(s.uStorer, s.cacher))
	v1.POST("/verify-code", handlers.VerifyCode(s.uStorer, s.cacher))
	v1.POST("/refresh-token", bearerMiddleware, handlers.RefreshToken(s.uStorer))
	v1.GET("/me", bearerMiddleware, handlers.GetMe(s.uStorer))

	v1.POST("/register", handlers.RegisterUser(s.uStorer, s.cacher))
	v1.POST("/verify-email", bearerMiddleware, handlers.VerifyEmail(s.uStorer, s.cacher))
	v1.POST("/reset-password", handlers.ResetPassword(s.uStorer, s.cacher))
	v1.POST("/reset-password/:id", handlers.ChangePassword(s.uStorer, s.cacher))
	v1.GET("/reset-password/:id", handlers.CheckChangePasswordStatus(s.cacher))

	v1.PATCH("/set-active", bearerMiddleware, handlers.SetActive(s.uStorer))

	v1.DELETE("/delete-user/:id", apiKeyMiddleware, handlers.DeleteUser(s.uStorer))
	v1.POST("/change-roles", apiKeyMiddleware, handlers.ChangeRoles(s.uStorer))

	v1.POST("/plants", bearerMiddleware, handlers.CreatePlant(s.pStorer))
	v1.GET("/plants/:id", bearerMiddleware, handlers.GetPlantByID(s.pStorer))
	v1.GET("/plants", bearerMiddleware, handlers.GetPlantsByUserID(s.pStorer))
	v1.PATCH("/plants/:id", bearerMiddleware, handlers.UpdatePlant(s.pStorer))
	v1.DELETE("/plants/:id", bearerMiddleware, handlers.DeletePlant(s.pStorer))

	v1.POST("/cares", bearerMiddleware, handlers.CreateCare(s.cStorer))
	v1.GET("/cares/:id", bearerMiddleware, handlers.GetCareByID(s.cStorer))
	v1.GET("/cares/plant/:id", bearerMiddleware, handlers.GetPlantCares(s.cStorer))
	v1.PATCH("/cares/:id", bearerMiddleware, handlers.UpdateCare(s.cStorer))
	v1.DELETE("/cares/:id", bearerMiddleware, handlers.DeleteCare(s.cStorer))
}

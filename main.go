package main

import (
	"context"
	"os"

	"github.com/mathehluiz/plant-care-tracker/api"
	"github.com/mathehluiz/plant-care-tracker/internal/cache"
	"github.com/mathehluiz/plant-care-tracker/internal/db"
	"github.com/mathehluiz/plant-care-tracker/internal/db/repositories"
	l "github.com/mathehluiz/plant-care-tracker/pkg/logger"
	"github.com/mathehluiz/plant-care-tracker/pkg/mailer"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	client, err := db.Start(ctx)
	if err != nil {
		l.Logger.Fatal("Cannot start db", zap.Error(err))
	}
	defer client.Close()

	l.Logger.Info("Connected with Database successfully 🚀")
	
	cacheClient, err := cache.Start(ctx)
	if err != nil {
		l.Logger.Fatal("Cannot start cache", zap.Error(err))
	}
	l.Logger.Info("Connected with Cache successfully 🚀")

	defer cacheClient.Close()

	userStorage := repositories.NewUserRepository(client)
	plantStorage := repositories.NewPlantRepository(client)
	careStorage := repositories.NewCareRepository(client)

	mailer.Init(os.Getenv("RESEND_API_KEY"))

	sv := api.NewServer(userStorage, plantStorage, careStorage, cacheClient)
	sv.Start()
}

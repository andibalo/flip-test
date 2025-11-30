package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andibalo/flip-test/internal/config"
	"github.com/andibalo/flip-test/internal/controller"
	v1 "github.com/andibalo/flip-test/internal/controller/v1"
	"github.com/andibalo/flip-test/internal/middleware"
	"github.com/andibalo/flip-test/internal/repository"
	"github.com/andibalo/flip-test/internal/service"
	"github.com/andibalo/flip-test/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.InitConfig()
	appLogger := logger.InitLogger(cfg)

	router := gin.Default()

	router.Use(middleware.LogPreReq(appLogger))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost", "http://localhost:3000", "http://localhost:5173"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(corsConfig))

	transactionRepo := repository.NewTransactionRepository()
	transactionService := service.NewTransactionService(transactionRepo, appLogger)
	transactionController := v1.NewTransactionController(transactionService)

	transactionController.AddRoutes(router)
	registerHandlers(router, &controller.HealthCheck{})

	server := &http.Server{
		Addr:    cfg.AppAddress(),
		Handler: router,
	}

	go func() {
		appLogger.Info(fmt.Sprintf("Server starting on port %s", cfg.AppAddress()))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			appLogger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	appLogger.Info("Server exited")
}

func registerHandlers(g *gin.Engine, handlers ...controller.Handler) {
	for _, handler := range handlers {
		handler.AddRoutes(g)
	}
}

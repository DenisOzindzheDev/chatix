package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DenisOzindzheDev/chatix/auth/internal/config"
	"github.com/DenisOzindzheDev/chatix/auth/internal/github"
	"github.com/DenisOzindzheDev/chatix/auth/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	// config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config read error: %v", err)
	}
	fmt.Println(cfg)
	github.InitOAuth(cfg)

	//router
	router := gin.Default()
	handler.RegisterRoutes(router)

	// server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}
	go func() {
		log.Printf("Server is listening on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %v", err)
		}
	}()

	//gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Print("Shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server exiting")
}

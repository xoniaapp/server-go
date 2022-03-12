package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting...")

	gin.SetMode(gin.ReleaseMode)
	err := godotenv.Load()

	if err != nil {
		log.Fatalln("FAILED TO LOAD ENVIRONMENT VARIABLES!")
	}

	ds, err := initDS()

	if err != nil {
		log.Fatalf("UNABLE TO INITIALIZE SERVER! ERR: %v\n", err)
	}

	router, err := inject(ds)

	if err != nil {
		log.Fatalf("FAILED TO INJECT DATA SOURCES! ERR: %v\n", err)
	}

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("FAILED TO INITIALIZE SERVER! ERR: %v\n", err)
		}
	}()

	log.Printf("SERVER STARTED ON PORT %v\n", srv.Addr)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ds.close(); err != nil {
		log.Fatalf("ERR OCCURED: %v\n", err)
	}

	log.Println("SHUTTING DOWN SERVER...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("SERVER FORCED TO SHUT DOWN: %v\n", err)
	}
}

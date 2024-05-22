package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/umeshkaushik-21122000/golang-lru-cache/api"
	"github.com/umeshkaushik-21122000/golang-lru-cache/app"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("running backend server for LRU cache")

	//redisClient := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "redis987456",
	//})
	//
	//cacheClient := cache.NewCacheImpl(redisClient)

	lruClient := app.Constructor(15)

	handler := api.NewHandler(lruClient)

	router := mux.NewRouter()

	router.HandleFunc("/get/{key}", handler.Get).Methods(http.MethodGet)
	router.HandleFunc("/set", handler.Set).Methods(http.MethodPost)
	router.HandleFunc("/del/{key}", handler.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/ws", handler.GetAll)

	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handlerr := c.Handler(router)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop() // Ensure ticker is stopped to avoid resource leak
		for {
			select {
			case <-ticker.C:
				//fmt.Println("Tick at", time.Now())
				lruClient.DeleteExpired()
			}
		}
	}()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerr))

}

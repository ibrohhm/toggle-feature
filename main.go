package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/toggle-feature/connection"
	"github.com/toggle-feature/handler"
	"github.com/toggle-feature/repository"
	"github.com/toggle-feature/service"
)

func main() {
	e := godotenv.Load()
	if e != nil {
		log.Println(e)
	}

	client, ctx, cancel := connection.NewMongodbConn()
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	router := httprouter.New()
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Headers", "*")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})

	mdwr_handler := handler.NewMiddleware()
	toggleFeatureRepo := repository.NewToggleFeatureRepository(client)
	toggleFeatureService := service.NewToggleFeatureService(toggleFeatureRepo)
	toggleFeatureHandler := handler.NewToggleFeatureHandler(toggleFeatureService)
	router.GET("/healthz", Healthz)

	router.GET("/toggle-features", mdwr_handler.Middleware(toggleFeatureHandler.Index))
	router.GET("/toggle-features/:id", mdwr_handler.Middleware(toggleFeatureHandler.Get))
	router.POST("/toggle-features", mdwr_handler.Middleware(toggleFeatureHandler.Insert))
	router.PATCH("/toggle-features/:id", mdwr_handler.Middleware(toggleFeatureHandler.Update))
	router.DELETE("/toggle-features/:id", mdwr_handler.Middleware(toggleFeatureHandler.Delete))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Println("listen at port " + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		panic(fmt.Sprintf("%s: %s", "Failed to listen and server", err))
	}
}

func Healthz(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode("ok")
}

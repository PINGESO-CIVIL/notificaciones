package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"notificaciones/config"
	"notificaciones/controllers"
)

func setupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/notificaciones", controllers.GetNotificaciones)
	mux.HandleFunc("POST /api/notificaciones", controllers.CreateNotificacion)
	mux.HandleFunc("OPTIONS /api/notificaciones", controllers.OptionsNotificaciones)
	mux.HandleFunc("GET /api/v1/notificaciones", controllers.GetNotificaciones)
	mux.HandleFunc("POST /api/v1/notificaciones", controllers.CreateNotificacion)
	mux.HandleFunc("OPTIONS /api/v1/notificaciones", controllers.OptionsNotificaciones)

	// Healthchecks según estándar
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if config.DB == nil {
			http.Error(w, "Base de datos no inicializada", http.StatusServiceUnavailable)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := config.DB.PingContext(ctx); err != nil {
			http.Error(w, "Base de datos no disponible", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})
	return mux
}

func main() {
	// Conectar a PostgreSQL
	config.ConnectDB()

	// Definir rutas
	mux := setupRouter()

	// Levantar servidor
	fmt.Println("Servidor corriendo en el puerto 8080 (PostgreSQL)")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

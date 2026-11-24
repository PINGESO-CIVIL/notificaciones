package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func ConnectDB() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5440/notificaciones_db?sslmode=disable"
	}

	var err error
	DB, err = sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("Error al abrir conexión con PostgreSQL: %v", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = DB.PingContext(ctx); err != nil {
		log.Printf("Aviso: No se pudo conectar a PostgreSQL aún (%v). Se reintentará en cada consulta.", err)
	} else {
		fmt.Println("Conectado exitosamente a PostgreSQL!")
		initSchema()
	}
}

func initSchema() {
	if DB == nil {
		return
	}

	query := `
	CREATE TABLE IF NOT EXISTS notificaciones (
		id VARCHAR(36) PRIMARY KEY,
		destinatario VARCHAR(255) NOT NULL,
		asunto VARCHAR(255) NOT NULL,
		cuerpo TEXT NOT NULL,
		estado VARCHAR(50) NOT NULL DEFAULT 'pendiente',
		fecha_creacion TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := DB.Exec(query); err != nil {
		log.Printf("Aviso al inicializar esquema: %v", err)
	}
}

package config

import "testing"

func TestConnectDB(t *testing.T) {
	// Llamar a ConnectDB ejecuta la inicialización
	ConnectDB()

	if DB == nil {
		t.Fatal("Se esperaba una instancia sql.DB inicializada")
	}

	initSchema()
}

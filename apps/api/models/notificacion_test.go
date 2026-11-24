package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSerializacionNotificacion(t *testing.T) {
	ahora := time.Now().UTC().Truncate(time.Second)
	notif := Notificacion{
		ID:            "550e8400-e29b-41d4-a716-446655440000",
		Destinatario:  "usuario@ejemplo.cl",
		Asunto:        "Confirmación de reserva",
		Cuerpo:        "Su hora fue agendada correctamente.",
		Estado:        "enviado_mock",
		FechaCreacion: ahora,
	}

	datos, err := json.Marshal(notif)
	if err != nil {
		t.Fatalf("error serializando Notificacion: %v", err)
	}

	var deserializada Notificacion
	if err := json.Unmarshal(datos, &deserializada); err != nil {
		t.Fatalf("error deserializando Notificacion: %v", err)
	}

	if deserializada.Destinatario != notif.Destinatario || deserializada.Asunto != notif.Asunto {
		t.Errorf("los campos deserializados no coinciden: %+v", deserializada)
	}
}

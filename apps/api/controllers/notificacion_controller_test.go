package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"notificaciones/config"
)

func TestOptionsNotificaciones(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/notificaciones", nil)
	rec := httptest.NewRecorder()

	OptionsNotificaciones(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Se esperaba 200 OK, se obtuvo %d", rec.Code)
	}
}

func TestCreateNotificacionValidaciones(t *testing.T) {
	casos := []struct {
		nombre string
		body   string
	}{
		{"JSON malformado", `{invalido}`},
		{"Campos vacíos", `{"destinatario":"","asunto":"","cuerpo":""}`},
		{"Correo inválido", `{"destinatario":"no-es-correo","asunto":"Hola","cuerpo":"Mensaje"}`},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/notificaciones", bytes.NewBufferString(tc.body))
			rec := httptest.NewRecorder()

			CreateNotificacion(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Caso %s: se esperaba 400 Bad Request, se obtuvo %d", tc.nombre, rec.Code)
			}
		})
	}
}

func TestFlujoNotificaciones(t *testing.T) {
	config.ConnectDB()

	payload := `{"destinatario":"cliente@usach.cl","asunto":"Aviso de prueba","cuerpo":"Contenido de prueba"}`
	reqPost := httptest.NewRequest(http.MethodPost, "/api/v1/notificaciones", bytes.NewBufferString(payload))
	recPost := httptest.NewRecorder()

	CreateNotificacion(recPost, reqPost)

	if recPost.Code != http.StatusCreated {
		t.Fatalf("Se esperaba 201 Created, se obtuvo %d (%s)", recPost.Code, recPost.Body.String())
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/notificaciones", nil)
	recGet := httptest.NewRecorder()

	GetNotificaciones(recGet, reqGet)

	if recGet.Code != http.StatusOK && recGet.Code != http.StatusServiceUnavailable {
		t.Fatalf("Se esperaba 200 o 503 en GET, se obtuvo %d", recGet.Code)
	}
}

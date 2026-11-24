package controllers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	"notificaciones/config"
	"notificaciones/models"
)

const (
	headerContentType = "Content-Type"
	mimeJSON          = "application/json"
	mimeProblemJSON   = "application/problem+json"
	timeoutConsulta   = 5 * time.Second
)

func configurarCORS(w http.ResponseWriter) {
	origenPermitido := os.Getenv("FRONTEND_URL")
	if origenPermitido == "" {
		origenPermitido = "http://localhost:5100"
	}
	w.Header().Set("Access-Control-Allow-Origin", origenPermitido)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Idempotency-Key, X-Request-Id")
}

func responderProblema(w http.ResponseWriter, r *http.Request, estado int, titulo, detalle string) {
	w.Header().Set(headerContentType, mimeProblemJSON)
	w.WriteHeader(estado)
	problema := models.ProblemaRFC9457{
		Type:     "about:blank",
		Title:    titulo,
		Status:   estado,
		Detail:   detalle,
		Instance: r.URL.Path,
	}
	_ = json.NewEncoder(w).Encode(problema)
}

func generarUUIDv4() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

// OptionsNotificaciones responde a solicitudes preflight CORS.
func OptionsNotificaciones(w http.ResponseWriter, _ *http.Request) {
	configurarCORS(w)
	w.WriteHeader(http.StatusOK)
}

// GetNotificaciones retorna el historial de notificaciones registradas.
func GetNotificaciones(w http.ResponseWriter, r *http.Request) {
	configurarCORS(w)

	if config.DB == nil {
		responderProblema(w, r, http.StatusServiceUnavailable, "Servicio no disponible", "La base de datos no está inicializada.")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutConsulta)
	defer cancel()

	filas, err := config.DB.QueryContext(ctx, `
		SELECT id, destinatario, asunto, cuerpo, estado, fecha_creacion
		FROM notificaciones
		ORDER BY fecha_creacion DESC
	`)
	if err != nil {
		responderProblema(w, r, http.StatusInternalServerError, "Error interno", "No se pudieron consultar las notificaciones.")
		return
	}
	defer filas.Close()

	lista := make([]models.Notificacion, 0)
	for filas.Next() {
		var n models.Notificacion
		if errScan := filas.Scan(&n.ID, &n.Destinatario, &n.Asunto, &n.Cuerpo, &n.Estado, &n.FechaCreacion); errScan != nil {
			responderProblema(w, r, http.StatusInternalServerError, "Error interno", "Error al procesar registros de notificaciones.")
			return
		}
		lista = append(lista, n)
	}

	if errRows := filas.Err(); errRows != nil {
		responderProblema(w, r, http.StatusInternalServerError, "Error interno", "Error al procesar registros de notificaciones.")
		return
	}

	w.Header().Set(headerContentType, mimeJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(lista)
}

// CreateNotificacion valida el payload de correo y registra el despacho.
func CreateNotificacion(w http.ResponseWriter, r *http.Request) {
	configurarCORS(w)

	var req models.NotificacionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responderProblema(w, r, http.StatusBadRequest, "Solicitud inválida", "El cuerpo JSON está malformado.")
		return
	}

	req.Destinatario = strings.TrimSpace(req.Destinatario)
	req.Asunto = strings.TrimSpace(req.Asunto)
	req.Cuerpo = strings.TrimSpace(req.Cuerpo)

	if req.Destinatario == "" || req.Asunto == "" || req.Cuerpo == "" {
		responderProblema(w, r, http.StatusBadRequest, "Campos obligatorios faltantes", "Los campos 'destinatario', 'asunto' y 'cuerpo' son obligatorios.")
		return
	}

	if _, err := mail.ParseAddress(req.Destinatario); err != nil {
		responderProblema(w, r, http.StatusBadRequest, "Destinatario inválido", "El formato del correo electrónico no es válido.")
		return
	}

	id, err := generarUUIDv4()
	if err != nil {
		responderProblema(w, r, http.StatusInternalServerError, "Error interno", "No se pudo generar el identificador de la notificación.")
		return
	}

	notif := models.Notificacion{
		ID:            id,
		Destinatario:  req.Destinatario,
		Asunto:        req.Asunto,
		Cuerpo:        req.Cuerpo,
		Estado:        "enviado_mock",
		FechaCreacion: time.Now().UTC(),
	}

	if config.DB != nil {
		ctx, cancel := context.WithTimeout(r.Context(), timeoutConsulta)
		defer cancel()

		_, errDB := config.DB.ExecContext(ctx, `
			INSERT INTO notificaciones (id, destinatario, asunto, cuerpo, estado, fecha_creacion)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, notif.ID, notif.Destinatario, notif.Asunto, notif.Cuerpo, notif.Estado, notif.FechaCreacion)
		if errDB != nil {
			responderProblema(w, r, http.StatusInternalServerError, "Error de persistencia", "No se pudo registrar la notificación en base de datos.")
			return
		}
	}

	w.Header().Set(headerContentType, mimeJSON)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(notif)
}

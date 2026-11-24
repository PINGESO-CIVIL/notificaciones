package models

import "time"

// NotificacionRequest representa el payload de entrada para despachar un correo.
type NotificacionRequest struct {
	Destinatario string `json:"destinatario"`
	Asunto       string `json:"asunto"`
	Cuerpo       string `json:"cuerpo"`
}

// Notificacion representa el registro persistido de una notificación enviada o simulada.
type Notificacion struct {
	ID            string    `json:"id"`
	Destinatario  string    `json:"destinatario"`
	Asunto        string    `json:"asunto"`
	Cuerpo        string    `json:"cuerpo"`
	Estado        string    `json:"estado"`
	FechaCreacion time.Time `json:"fecha_creacion"`
}

// ProblemaRFC9457 define el formato estándar de errores HTTP (application/problem+json).
type ProblemaRFC9457 struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
}

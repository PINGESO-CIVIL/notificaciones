# Gateway de Notificaciones — Grupo 10 (DAS)

Repositorio del microservicio **Paquete de Trabajo 10: Gateway de Notificaciones** desarrollado bajo los lineamientos técnicos del curso **DAS (Desarrollo de Microservicios)**.

### Enlaces de Evidencia
* **Pipeline de CI (GitHub Actions):** *(Pegar aquí la URL de la pestaña Actions de su repositorio en GitHub)*
* **Panel de SonarQube:** *(Pegar aquí la URL pública/túnel del dashboard con `id=notificaciones`)* (Clave del proyecto: `notificaciones`)

---

## 1. Descripción del Servicio (PT10)

Servicio genérico e independiente encargado de recibir payloads de notificación (`destinatario`, `asunto`, `cuerpo`) vía API REST y despachar correos electrónicos mediante protocolo SMTP.

* **Dependencias de otros grupos:** Ninguna (servicio base independiente).
* **Servicios consumidores:** Motor de Agendamiento (Grupo 9) y otros módulos que requieran notificaciones por correo.

---

## 2. Composición y Estructura del Repositorio

```text
Grupo10-PINGESO/
├── apps/
│   ├── api/                                   # Backend en Go
│   │   ├── config/
│   │   │   ├── db.go                          # Conexión PostgreSQL (pgx/v5) y esquema notificaciones
│   │   │   └── db_test.go                     # Pruebas de integración de base de datos
│   │   ├── controllers/
│   │   │   ├── notificacion_controller.go     # Handlers REST, validaciones y errores RFC 9457
│   │   │   └── notificacion_controller_test.go# Pruebas unitarias y de integración
│   │   ├── models/
│   │   │   ├── notificacion.go                # Structs de request, response y error RFC 9457
│   │   │   └── notificacion_test.go           # Pruebas de serialización JSON
│   │   ├── Dockerfile                         # Multi-stage build (alpine:3.23, USER 10001:10001, healthcheck)
│   │   ├── go.mod                             # Módulo Go (`notificaciones`, Go 1.26, pgx v5.10.0)
│   │   ├── go.sum                             # Checksums verificados
│   │   ├── main.go                            # Enrutador HTTP (/healthz, /readyz, /api/v1/notificaciones)
│   │   └── main_test.go                       # Pruebas de liveness, readiness y rutas
│   └── web/                                   # Frontend SPA (Vue.js / Vite)
│       ├── Dockerfile                         # Multi-stage build (node:22-alpine + nginx-unprivileged)
│       ├── nginx.conf                         # Servidor web en puerto interno 8080
│       └── package.json                       # Dependencias Vue 3 + Vite
├── contracts/
│   └── openapi.yaml                           # Contrato OpenAPI 3.1.0 (API-First)
├── docs/
│   └── modos-de-falla.md                      # Matriz de resiliencia y modos de falla
├── .github/
│   └── workflows/
│       └── ci.yml                             # Pipeline CI (calidad, seguridad e imagen)
├── .env.example                               # Plantilla de variables de entorno y SMTP
├── .trivyignore.yaml                          # Excepciones autorizadas de Trivy
├── docker-compose.yaml                        # Orquestación local (PostgreSQL 17, Mailpit, API y Web)
├── sonar-project.properties                   # Configuración SonarQube (`projectKey=notificaciones`)
└── README.md                                  # Documentación y evidencias del Hito 0
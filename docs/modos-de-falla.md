# Modos de Falla — Gateway de Notificaciones (Grupo 10)

Documento de resiliencia y modos de falla del microservicio **Gateway de Notificaciones (PT10)** según el estándar **DAS (Desarrollo de Microservicios)**.

---

## 1. Matriz de Modos de Falla

| ID | Modo de Falla | Causa Raíz Potencial | Mitigación / Qué la Impide | Código HTTP / Comportamiento | Cómo se Verifica |
|:---|:---|:---|:---|:---:|:---|
| **MF-01** | Base de datos PostgreSQL no disponible | Contenedor apagado, problemas de red o credenciales erróneas | Conexión con timeout explícito (5s), reintento por consulta y endpoint `/readyz` | `503 Service Unavailable` / `500 Internal Server Error` (RFC 9457 sin trazas) | `GET /readyz` devuelve 503; pruebas en `db_test.go` |

| **MF-02** | Payload JSON inválido o malformado | Cliente envía sintaxis JSON corrupta | Decodificación temprana en el handler antes de invocar lógica o BD | `400 Bad Request` (`application/problem+json`) | `POST /api/v1/notificaciones` con `{invalido}` en `notificacion_controller_test.go` |

| **MF-03** | Campos obligatorios faltantes (`destinatario`, `asunto`, `cuerpo`) | Servicio consumidor omite campos del contrato | Validación estricta de strings no vacíos en el controlador | `400 Bad Request` (`application/problem+json`) | `POST /api/v1/notificaciones` con campos vacíos en `notificacion_controller_test.go` |

| **MF-04** | Formato de correo electrónico inválido | Dirección en `destinatario` no cumple estándar RFC 5322 | Validación mediante `net/mail.ParseAddress` antes de persistir o enviar | `400 Bad Request` (`application/problem+json`) | `POST /api/v1/notificaciones` con `"no-es-correo"` en `notificacion_controller_test.go` |

| **MF-05** | Servidor SMTP caído o con latencia alta | Falla del relay SMTP externo o caída de Mailpit en local | Timeout explícito de conexión SMTP (5s) y registro del estado `fallido` en BD | `502 Bad Gateway` / `500 Internal Server Error` sin colgar el hilo | Prueba de timeout contra servidor SMTP no disponible (Hito 1/2) |

| **MF-06** | Envío duplicado de correos por reintentos | Servicio consumidor (ej. Grupo 9) reintenta tras timeout de red | Soporte de cabecera `Idempotency-Key` para deduplicar solicitudes repetidas | `200 OK` / `201 Created` retornando el registro previo sin reenviar correo | Dos `POST` consecutivos con el mismo `Idempotency-Key` |

| **MF-07** | Solicitud OPTIONS / CORS desde orígenes no autorizados | Peticiones preflight desde navegadores externos | Cabeceras CORS centralizadas (`configurarCORS`) restringidas a `FRONTEND_URL` | `200 OK` solo para origen configurado | `OPTIONS /api/v1/notificaciones` en `main_test.go` |

| **MF-08** | Inyección SQL en asunto o cuerpo del correo | Entrada maliciosa en campos de texto libre | Consultas 100% parametrizadas con marcadores nativos (`$1, $2, ...`) | Prevención total a nivel de driver `pgx` | Análisis estático SonarQube y pruebas de integración |

---

## 2. Estrategia de Verificación y Salud

1. **Salud del Proceso (`GET /healthz`):**
   * Verifica que el servidor HTTP esté vivo y aceptando solicitudes.
   * Retorna `200 OK` con cuerpo `ok`.

2. **Salud de la Dependencia (`GET /readyz`):**
   * Ejecuta un ping activo a PostgreSQL con timeout de 2 segundos.
   * Retorna `200 OK` (`ready`) si la base responde, o `503 Service Unavailable` si no está disponible.

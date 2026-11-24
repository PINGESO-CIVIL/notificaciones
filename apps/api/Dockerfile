FROM golang:1.26-alpine AS build

WORKDIR /app

# Copiar archivos de dependencias desde apps/api
COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download

# Copiar el código fuente y compilar binario reproducible
COPY apps/api/ ./
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -buildvcs=false -trimpath \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /app/notificaciones .

# Runtime ligero
FROM alpine:3.23
WORKDIR /app
RUN apk --no-cache upgrade && \
    addgroup -S nonroot -g 10001 && adduser -S nonroot -G nonroot -u 10001
COPY --from=build --chown=10001:10001 /app/notificaciones /app/notificaciones

USER 10001:10001
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8080/healthz || exit 1
CMD ["/app/notificaciones"]

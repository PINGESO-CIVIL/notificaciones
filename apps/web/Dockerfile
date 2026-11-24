# Etapa 1: Construcción (Build)
FROM node:22-alpine AS build-stage
WORKDIR /app
COPY apps/web/package.json apps/web/package-lock.json ./
RUN npm ci
COPY apps/web/ ./
RUN npm run build

# Etapa 2: Producción (Nginx sin privilegios)
FROM nginxinc/nginx-unprivileged:alpine AS production-stage
USER 0:0
RUN apk --no-cache upgrade
COPY --from=build-stage /app/dist /usr/share/nginx/html
COPY apps/web/nginx.conf /etc/nginx/conf.d/default.conf
USER 101:101
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8080/ || exit 1
CMD ["nginx", "-g", "daemon off;"]

# 1) Builder
FROM golang:1.25.3-alpine AS builder
WORKDIR /app
RUN apk add --no-cache build-base git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Собираем бинарь в корень файловой системы билдера
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./cmd/db_demo

# 2) Runtime
FROM alpine:3.20
WORKDIR /app
RUN adduser -D app && chown -R app:app /app
USER app
# Копируем ИМЕННО файл бинаря
COPY --from=builder /server /app/server

EXPOSE 8080 6060
ENV PORT=8080 JWT_SECRET=super-secret JWT_ISSUER=day2 JWT_AUDIENCE=day2-clients
ENTRYPOINT ["/app/server"]
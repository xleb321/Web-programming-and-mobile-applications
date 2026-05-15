# Этап сборки
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем файлы модуля и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/rugram-api ./cmd/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/rugram-api .

EXPOSE 4200

CMD ["./rugram-api"]
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем файлы с зависимостями
COPY go.mod go.sum ./

# Настраиваем GOPROXY для лучшей загрузки
RUN go env -w GOPROXY=https://proxy.golang.org,direct

# Проверяем наличие файлов и загружаем зависимости
RUN ls -la && \
    if [ -f go.mod ]; then \
        go mod download; \
    else \
        echo "go.mod not found!" && exit 1; \
    fi

# Копируем остальной код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env* ./
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main"]
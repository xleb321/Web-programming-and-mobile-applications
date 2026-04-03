# Используем официальный образ Go
FROM golang:1.25-alpine

# Устанавливаем необходимые системные пакеты
RUN apk add --no-cache git ca-certificates tzdata

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы с зависимостями
COPY go.mod go.sum ./

# Настраиваем GOPROXY (важно для загрузки)
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://proxy.golang.org,direct

# Загружаем зависимости (с дополнительной диагностикой)
RUN go mod download || (echo "Failed to download deps" && cat go.mod && exit 1)

# Копируем весь код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["/app/main"]
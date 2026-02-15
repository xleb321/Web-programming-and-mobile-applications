# Используем официальный образ Go
FROM golang:1.21-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы с зависимостями
COPY go.mod ./
RUN go mod download

# Копируем исходный код
COPY *.go ./

# Собираем приложение
RUN go build -o /new-year-counter

# Указываем порт
EXPOSE 3000

# Запускаем приложение
CMD ["/new-year-counter"]
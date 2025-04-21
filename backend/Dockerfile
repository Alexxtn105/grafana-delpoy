# Используем официальный образ Go
FROM golang:1.23-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git make

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы модуля
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарник из builder-стадии
COPY --from=builder /app/app .

# Копируем статические файлы (если есть)
# COPY --from=builder /app/static ./static

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./app"]
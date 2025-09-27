FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарный файл
COPY --from=builder /app/main .

# Создаем директорию для временных файлов
RUN mkdir -p /tmp

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
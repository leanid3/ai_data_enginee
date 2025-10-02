FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Обновляем зависимости и собираем приложение
RUN go mod tidy && go build -o main ./cmd/http_main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Копируем бинарный файл
COPY --from=builder /app/main .

# Копируем proto файлы (если нужны для runtime)
COPY --from=builder /app/proto ./proto
# Копируем gen файлы профайлера
COPY --from=builder /app/profiler-gen ./profiler-gen

EXPOSE 8080 50051

CMD ["./main"]

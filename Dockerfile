# ---- Этап 1: сборка ----
FROM golang:1.26-alpine AS build
WORKDIR /src
 
# Сначала зависимости: Docker кэширует этот слой
COPY go.mod go.sum ./
RUN go mod download
 
COPY . .
# CGO выключен: modernc.org/sqlite написан на чистом Go
RUN CGO_ENABLED=0 go build -o /todo .
 
# ---- Этап 2: итоговый образ ----
FROM alpine:3.20
WORKDIR /app
 
COPY --from=build /todo /app/todo
COPY web /app/web
 
RUN mkdir -p /data
 
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
 
EXPOSE 7540
CMD ["/app/todo"]
 
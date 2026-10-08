# Tahap 1: build binary Go
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app ./cmd

# Tahap 2: image kecil untuk menjalankan
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY --from=build /app /app
EXPOSE 8080
CMD ["/app"]
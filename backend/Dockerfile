# --- Stage 1: Build ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy file dependency & vendor
COPY go.mod go.sum ./

RUN go mod download
RUN go mod vendor

# Copy source code (cmd & internal)
COPY cmd ./cmd
COPY internal ./internal

# Build binary (gunakan flag -mod=vendor)
RUN go build -mod=vendor -o main cmd/api/main.go

# --- Stage 2: Runtime ---
FROM alpine:latest

WORKDIR /root/

# Copy binary dari stage builder
COPY --from=builder /app/main .

# Copy folder migrations (WAJIB ADA untuk database)
COPY migrations ./migrations

# Install timezone & certificates (Penting untuk HTTPS & Waktu Jakarta)
RUN apk --no-cache add ca-certificates tzdata

# Set Environment Default
ENV PORT=8080
ENV APP_ENV=production

EXPOSE 8080

CMD ["./main"]
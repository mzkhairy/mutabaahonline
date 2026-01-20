# --- Stage 1: Build ---
FROM golang:alpine AS builder

WORKDIR /app

# Install Git (Wajib untuk download module)
RUN apk add --no-cache git

# 1. Copy file definisi module
COPY go.mod go.sum ./

# 2. Download library ke cache sistem
RUN go mod download

# 3. Copy source code
COPY . .

# 4. Build Binary (HAPUS flag -mod=vendor)
RUN go build -o main cmd/api/main.go

# --- Stage 2: Runtime ---
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY migrations ./migrations

# Install timezone & certificates
RUN apk --no-cache add ca-certificates tzdata

ENV PORT=8080
ENV APP_ENV=production

EXPOSE 8080

CMD ["./main"]
# --- Stage 1: Build Vue ---
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package.json
COPY package*.json ./
RUN npm install

# Copy source code frontend
COPY . .

# Terima argumen URL API saat build (PENTING)
ARG VITE_API_BASE_URL
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

# Build production
RUN npm run build

# --- Stage 2: Serve with Nginx ---
FROM nginx:alpine

# Copy hasil build ke folder Nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy konfigurasi Nginx custom
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
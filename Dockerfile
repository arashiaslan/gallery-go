# Stage 1: Build frontend
FROM node:22-alpine AS frontend-builder

WORKDIR /frontend

# Copy frontend files
COPY frontend/package*.json ./
RUN npm ci

COPY frontend . 

# Build frontend
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.24-alpine AS backend-builder

WORKDIR /build

# Copy Go dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend code
COPY backend/ ./

# Replace the dev placeholder dist with the real frontend build
RUN rm -rf ./internal/web/dist && mkdir -p ./internal/web/dist
COPY --from=frontend-builder /frontend/dist ./internal/web/dist

# Build Go binary with embedded frontend
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Stage 3: Runtime image (minimal)
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=backend-builder /build/server .

# Create data directory for SQLite database and uploads
RUN mkdir -p ./data

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run server
CMD ["./server"]

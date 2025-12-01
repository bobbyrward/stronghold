# Stage 1: Build Vue frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/stronghold .

# Stage 3: Production image
FROM alpine:3.21
RUN apk --no-cache add ca-certificates ffmpeg
WORKDIR /app
COPY --from=backend-builder /out/stronghold .
COPY --from=backend-builder /app/web/dist ./web/dist
EXPOSE 8000
CMD ["./stronghold", "www"]

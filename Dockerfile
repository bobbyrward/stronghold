FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/stronghold .

# Stage 2: Create the final image
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /out/stronghold .
CMD ["./stronghold", "/etc/stronghold/config.yaml"]

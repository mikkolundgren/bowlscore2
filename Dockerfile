# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod init bowlscore && go mod tidy
RUN go build -o bowlscore .

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates libc6-compat
WORKDIR /app

# Copy binary and assets
COPY --from=builder /app/bowlscore .
COPY --from=builder /app/static ./static
COPY --from=builder /app/certs ./certs

# Create directory for database file
RUN mkdir data

# Set environment variables
ENV DB_PATH=/app/data/bowling_scores.db

# Expose HTTPS port
EXPOSE 8443

# Run the application
CMD ["./bowlscore"]

#To build and run:
#docker build -t bowlscore .
#docker run -p 8443:8443 -v ./data:/app/data bowlscore

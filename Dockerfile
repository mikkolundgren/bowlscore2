# Build stage
FROM --platform=linux/amd64 golang:1.23-alpine AS builder
WORKDIR /app
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o bowlscore .

# Runtime stage
FROM --platform=linux/amd64 alpine:latest
RUN apk --no-cache add ca-certificates libc6-compat
WORKDIR /app

# Copy binary and assets
COPY --from=builder /app/bowlscore .
COPY --from=builder /app/static ./static

# Create directory for database file
RUN mkdir data

# Set environment variables
ENV DB_PATH=/app/data/bowling_scores.db

# Expose HTTPS port
EXPOSE 8080

# Run the application
CMD ["./app/bowlscore"]

#To build and run:
#docker build -t bowlscore .
#docker run -p 8443:8443 -v ./data:/app/data bowlscore

FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY templates/ ./templates/
COPY static/ ./static/
COPY content/ ./content/

EXPOSE 8080

CMD ["./server"]

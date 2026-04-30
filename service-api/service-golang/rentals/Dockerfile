FROM golang:1.24-alpine AS builder
WORKDIR /workspace

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/rentals ./cmd/rentals

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/rentals /app/rentals
EXPOSE 8096
ENTRYPOINT ["/app/rentals"]

# Build stage for the sales service.
FROM golang:1.24-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /out/sales ./cmd/sales

# Runtime stage kept small for local and remote execution.
FROM alpine:3.21

WORKDIR /app

COPY --from=build /out/sales /usr/local/bin/sales

EXPOSE 8087

ENTRYPOINT ["sales"]

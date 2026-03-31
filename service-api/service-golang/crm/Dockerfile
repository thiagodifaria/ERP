# Build stage for the crm service.
FROM golang:1.24-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /out/crm ./cmd/crm

# Runtime stage kept small for local and remote execution.
FROM alpine:3.21

WORKDIR /app

COPY --from=build /out/crm /usr/local/bin/crm

EXPOSE 8083

ENTRYPOINT ["crm"]

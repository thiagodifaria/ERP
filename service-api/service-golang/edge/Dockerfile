# Build stage for the edge service.
FROM golang:1.24-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /out/edge ./cmd/edge

# Runtime stage kept small for local and remote execution.
FROM alpine:3.21

WORKDIR /app

COPY --from=build /out/edge /usr/local/bin/edge

EXPOSE 8080

ENTRYPOINT ["edge"]

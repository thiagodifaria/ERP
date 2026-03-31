# Build stage for the webhook ingress service.
FROM rust:1 AS build

WORKDIR /src

COPY Cargo.toml ./
COPY src ./src

RUN cargo build --release

# Runtime stage focused on a small executable image.
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=build /src/target/release/webhook-hub /usr/local/bin/webhook-hub

EXPOSE 8082

ENTRYPOINT ["webhook-hub"]

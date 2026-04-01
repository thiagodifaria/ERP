FROM elixir:1.17-alpine

RUN apk add --no-cache build-base git

WORKDIR /app

COPY mix.exs ./mix.exs
COPY config ./config

RUN mix local.hex --force && \
    mix local.rebar --force && \
    mix deps.get

COPY lib ./lib

ENV MIX_ENV=prod

RUN mix deps.compile

EXPOSE 8085

CMD ["mix", "run", "--no-halt"]

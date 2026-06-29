# Build stage

FROM golang:1.26-alpine3.24 AS build-stage

WORKDIR /app

COPY --link go.mod go.sum* .
RUN go mod download

COPY --link --exclude=web . .

RUN go build -o web-app ./cmd/server/main.go


# Final run stage

FROM alpine:3.24 AS final-stage

WORKDIR /app

COPY --link web web
COPY --link --from=build-stage /app/web-app .

EXPOSE 80

CMD ["./web-app"]

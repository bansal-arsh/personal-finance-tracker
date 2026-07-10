# Build stage

FROM golang:1.26-alpine3.24 AS build-stage

WORKDIR /server

COPY --link go.mod go.sum ./
RUN go mod download

COPY --link --exclude=web . .

RUN go build -o web-server ./cmd/server/main.go


# Final run stage

FROM alpine:3.24 AS final-stage

WORKDIR /server

COPY --link web web
COPY --link --from=build-stage /server/web-server .

CMD ["./web-server"]

# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS dev
RUN apk --no-cache add make gcc musl-dev
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 3000
CMD ["go", "run", "./cmd/api"]

FROM golang:1.26-alpine AS builder
WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /out/api ./api

EXPOSE 3000

USER nobody:nobody

ENTRYPOINT ["./api"]

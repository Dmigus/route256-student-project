FROM golang:1.22-alpine AS build-stage

WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build -o bin/app ./cmd

FROM alpine AS release-stage

WORKDIR /app

COPY --from=build-stage /app/bin/app /app/bin/app
COPY --from=build-stage /app/configs /app/configs/

EXPOSE 8080

CMD ["./bin/app"]
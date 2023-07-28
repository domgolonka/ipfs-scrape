FROM golang:1.20-alpine as builder

RUN apk update && apk upgrade
RUN apk --no-cache add git

RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on

COPY . ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blockparty main.go

# Run container
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add tzdata

# The config file to use when running.
ARG config

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/blockparty .

ENTRYPOINT ["./blockparty"]

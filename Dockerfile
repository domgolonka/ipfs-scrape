FROM golang:1.20.3-alpine as builder

RUN apk update && apk upgrade
RUN apk --no-cache add git make gcc musl-dev


WORKDIR . /app
ARG ENV

ADD . /app

ARG ep=api

RUN #scripts/build-info.sh > ./build.json

RUN cd /app && make build-$ep-static
RUN file="$(ls -1 /app)" && echo $file

# Run container
FROM alpine:latest as launch
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add tzdata

# The config file to use when running.
ARG config
ARG ep=core
# The environment is set to argument for the entrypoint to work.
ENV APP=$ep

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/bin/$ep .

ENTRYPOINT ./${APP}
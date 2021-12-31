# syntax=docker/dockerfile:1

FROM golang:1.17.5 AS build

# Create a user so Docker doesn't run as root
ENV APP_USER app
ENV APP_HOME $GOPATH/src/app

RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER \
	&& mkdir -p $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME

WORKDIR $APP_HOME
USER $APP_USER

# Copy
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY cmd/webserver ./cmd/webserver

RUN go install ./...

FROM frolvlad/alpine-glibc

RUN apk --update upgrade
RUN apk add sqlite

# removing apk cache
RUN rm -rf /var/cache/apk/*

# Create a user so Docker doesn't run as root
ENV APP_USER app
ENV APP_HOME /app

ARG GROUP_ID
ARG USER_ID
ARG PORT

RUN addgroup -g $GROUP_ID -S $APP_USER && adduser -u $USER_ID -S $APP_USER -G $APP_USER \
	&& mkdir -p $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME

USER $APP_USER
WORKDIR $APP_HOME

COPY --from=build /go/bin/webserver ./webserver

ENV PORT $PORT
EXPOSE $PORT

ENTRYPOINT ["./webserver"]

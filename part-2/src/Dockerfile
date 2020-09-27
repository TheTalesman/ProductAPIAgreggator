FROM golang:1.14.3-alpine as customgo

ENV SERVICE_NAME my_project
ENV APP /src/${SERVICE_NAME}/
ENV WORKDIR ${GOPATH}${APP}
WORKDIR $WORKDIR
ADD . $WORKDIR


RUN apk update && apk upgrade && apk add --no-cache bash git openssh
RUN apk add --no-cache \
        libc6-compat

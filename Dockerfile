FROM alpine
LABEL MAINTAINER="heixiushamao@gmail.com"

ENV DEFAULT_DB_NAME short-url.db
ENV DEFAULT_SCHEMA http
ENV DEFAULT_DOMAIN_NAME 127.0.0.1:8080
ENV TIME_ZONE=Asia/Singapore
RUN ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone
RUN apk add --no-cache tzdata

##CA证书，https请求
RUN apk add --no-cache ca-certificates

WORKDIR /go
COPY main /go
RUN mkdir data

ENTRYPOINT ./main -d=$DEFAULT_DB_NAME -s=$DEFAULT_SCHEMA -dn=$DEFAULT_DOMAIN_NAME
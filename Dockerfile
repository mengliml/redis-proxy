# DEVELOPMENT 
FROM golang:1.9.2-alpine

ARG app_env
ENV APP_ENV $app_env
ENV REDIS_ADDRESS=redis:6379
ENV GLOBAL_EXPIRY=60000
ENV CAPACITY=1000
ENV PORT=8080
ENV MAX_CLIENTS=10

# Install git and dep
RUN apk add --no-cache git
RUN go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/redis-proxy
WORKDIR /go/src/redis-proxy

# DEVELOPMENT:
RUN dep ensure
RUN go build
CMD ./redis-proxy -global-expiry 60000 -capacity 1000 -port 8080 -max-clients 10

EXPOSE 8080


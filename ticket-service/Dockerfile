FROM golang:1.15-alpine3.13
WORKDIR ticket-service
RUN apk update \
    && apk add sqlite \
    && apk add socat \
    && apk add build-base
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ticket-service
EXPOSE 8080
CMD ["./ticket-service"]
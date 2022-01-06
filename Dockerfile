FROM golang:1.16-alpine as build
WORKDIR /app
RUN apk add --update alpine-sdk
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /bad-pixel

FROM alpine:3.14.3
WORKDIR /
COPY --from=build /bad-pixel ./
EXPOSE 8080
ENTRYPOINT ["/bad-pixel"]
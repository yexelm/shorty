FROM golang:1.16-alpine AS build

LABEL stage=builder
RUN apk update && apk upgrade && \
    apk add --no-cache bash git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build ./cmd/main.go

FROM alpine:latest
COPY --from=build ./src/main main
EXPOSE 8080
CMD ["./main"]

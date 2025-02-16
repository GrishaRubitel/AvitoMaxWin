FROM golang:1.23.4-alpine AS builder
WORKDIR /usr/local/src
RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["./go.mod", "./go.sum", "./"]
RUN go mod download

COPY ./ ./
RUN go build -o ./bin/go .

FROM alpine
COPY .env .env
COPY --from=builder /usr/local/src/bin/go /
RUN apk --no-cache add bash
EXPOSE 8080
CMD [ "/go" ]
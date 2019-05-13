###############
# FIRST STAGE #
###############
FROM golang:1.12-alpine as builder

# Installing dependencies
RUN apk add git gcc musl-dev --update
RUN go get -u github.com/maxbrunsfeld/counterfeiter

# Bootstrapping modules dependencies
RUN mkdir -p /src/gopaytest/payments
WORKDIR /src/gopaytest/payments
COPY go.mod go.mod
COPY go.sum go.sum
RUN go get -d

# Copying source files after `go get` to retain modules cache as often as possible
COPY . /src/gopaytest/payments

# Running tests
RUN go generate ./...
RUN go test ./...

# Compiling binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api


################
# SECOND STAGE #
################
FROM scratch

ARG HTTP_PORT
ENV HTTP_PORT ${HTTP_PORT}
EXPOSE ${HTTP_PORT}

# Copy api binary from first step
COPY --from=builder /src/gopaytest/payments/api api

CMD ["./api"]

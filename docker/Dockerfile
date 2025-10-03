FROM golang:1.25-alpine

RUN apk add --update \
    curl \
    libc6-compat

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ${GOPATH}/bin v2.5.0

RUN go install mvdan.cc/gofumpt@latest
RUN go install github.com/golang/mock/mockgen@latest
RUN go install github.com/godoc-lint/godoc-lint/cmd/godoclint@v0.10.1

WORKDIR /app

COPY . .
RUN go mod download

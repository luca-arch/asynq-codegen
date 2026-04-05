FROM golang:1-alpine3.22 AS builder

WORKDIR /build

COPY package.go ./

RUN go mod init github.com/luca-arch/asynq-codegen-example \
    && go mod tidy \
    && go mod verify \
    && go get -tool github.com/luca-arch/asynq-codegen \
    && go generate ./... \
    && go build -o /bin/my-command .

ENTRYPOINT ["/bin/my-command"]

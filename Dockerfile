FROM golang:1.19-alpine3.17 as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-w -s"

FROM scratch
COPY --from=builder /workspace/ravelin-3ds-demo /workspace/ravelin-3ds-demo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8085

ENTRYPOINT ["/workspace/ravelin-3ds-demo"]
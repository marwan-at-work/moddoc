FROM golang:1.12 AS builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN GOPROXY=https://proxy.golang.org CGO_ENABLED=0 go build -o=/app/moddoc

FROM busybox

RUN mkdir /app

COPY --from=builder /app/moddoc /app/moddoc

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /app

ENTRYPOINT ["/app/moddoc"]
FROM golang:1.12 AS builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -mod=vendor -o=/app/goproxydoc

FROM busybox

RUN mkdir /app

COPY --from=builder /app/goproxydoc /app/goproxydoc

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /app/frontend/dist /app/frontend/dist

WORKDIR /app

ENTRYPOINT ["/app/goproxydoc"]
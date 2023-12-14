FROM golang:latest as builder

WORKDIR /codebase

COPY bin/server .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /codebase/server ./

ENTRYPOINT ["./server"]

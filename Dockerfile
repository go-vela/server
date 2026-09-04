# SPDX-License-Identifier: Apache-2.0

FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b as certs

RUN apk add --update --no-cache ca-certificates

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080

ENV GODEBUG=netdns=go

ADD release/vela-server /bin/

CMD ["/bin/vela-server"]

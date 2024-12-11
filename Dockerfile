# SPDX-License-Identifier: Apache-2.0

FROM alpine:3.20.3@sha256:beefdbd8a1da6d2915566fde36db9db0b524eb737fc57cd1367effd16dc0d06d as certs

RUN wget --quiet --output-document=/etc/ssl/certs/ca-certificates.crt "http://browserconfig.target.com/tgt-certs/tgt-ca-bundle.crt"

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080

ENV GODEBUG=netdns=go

ADD release/vela-server /bin/

CMD ["/bin/vela-server"]

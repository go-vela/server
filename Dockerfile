# SPDX-License-Identifier: Apache-2.0

FROM alpine:3.22.0@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715 as certs

RUN wget --quiet --output-document=/etc/ssl/certs/ca-certificates.crt "http://browserconfig.target.com/tgt-certs/tgt-ca-bundle.crt"

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080

ENV GODEBUG=netdns=go

ADD release/vela-server /bin/

CMD ["/bin/vela-server"]

FROM docker.io/golang:alpine as builder

COPY . /src
WORKDIR /src
RUN apk add git make && make

FROM docker.io/alpine

LABEL maintainer="George <zhoreeq@users.noreply.github.com>"

RUN mkdir -p /etc/ipfd && chown -R nobody:nobody /etc/ipfd

COPY --from=builder /src/ipfd /usr/bin/ipfd
COPY static /etc/ipfd/static
COPY templates /etc/ipfd/templates
COPY config.example /etc/ipfd/config

USER nobody
VOLUME [ "/etc/ipfd" ]
WORKDIR /etc/ipfd
EXPOSE 8000

CMD ["/usr/bin/ipfd", "-config", "/etc/ipfd/config"]

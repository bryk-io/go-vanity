FROM busybox:1.31.1

ARG VERSION_TAG
LABEL maintainer="Ben Cessa <ben@bryk.io>"
LABEL version=${VERSION_TAG}

COPY govanity_${VERSION_TAG}_linux_amd64 /usr/bin/govanity
COPY ca-roots.crt /etc/ssl/certs/ca-roots.crt

VOLUME ["/tmp", "/etc/govanity"]

EXPOSE 9090/tcp

ENTRYPOINT ["/usr/bin/govanity"]
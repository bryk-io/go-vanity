FROM busybox:1.31.1

# Metadata
ARG VERSION_TAG
LABEL maintainer="Ben Cessa <ben@pixative.com>"
LABEL version=${VERSION_TAG}

# Install trusted certificate roots
COPY ca-roots.crt /etc/ssl/certs/

# Run as an unprivileged user
ENV USER=guest
ENV UID=10001
RUN adduser -h /home/${USER} -g "container-user" -s /bin/sh -D -u 10001 ${USER} ${USER}
USER ${USER}:${USER}

# Expose required ports and volumes
VOLUME ["/tmp", "/etc/govanity"]
EXPOSE 9090/tcp

# Add application binary and use it as default entrypoint
COPY govanity_${VERSION_TAG}_linux_amd64 /usr/bin/govanity
ENTRYPOINT ["/usr/bin/govanity"]

FROM docker.pkg.github.com/bryk-io/base-images/shell:0.1.0

# Metadata
ARG VERSION
LABEL maintainer="Ben Cessa <ben@pixative.com>"
LABEL version=${VERSION}

# Expose required ports and volumes
VOLUME ["/tmp", "/etc/govanity"]
EXPOSE 9090/tcp

# Add application binary and use it as default entrypoint
COPY govanity_linux_amd64 /usr/bin/govanity
ENTRYPOINT ["/usr/bin/govanity"]

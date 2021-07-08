FROM ghcr.io/bryk-io/shell:0.2.0

# Expose required ports and volumes
VOLUME ["/tmp", "/etc/govanity"]
EXPOSE 9090/tcp

# Add application binary and use it as default entrypoint
COPY govanity /usr/bin/govanity
ENTRYPOINT ["/usr/bin/govanity"]

FROM golang:1.9

ADD . /go/src/github-webhook

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN cd /go/src/github-webhook;make

# Run the outyet command by default when the container starts.
USER 10000
WORKDIR /go/src/github-webhook
ENTRYPOINT ./github-webhook

# Document that the service listens on port 8080.
EXPOSE 5000

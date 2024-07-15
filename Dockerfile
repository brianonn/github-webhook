FROM golang:1.21

ADD . /go/src/github-webhook

# Build the github-webhook command inside the container.
# the Makefile will take care of dependancies too
RUN cd /go/src/github-webhook && make

# Run the github-webhook command by default when the container starts.
USER 10000
WORKDIR /go/src/github-webhook
ENTRYPOINT ./github-webhook

# Document that the service listens on port 5000.
EXPOSE 5000

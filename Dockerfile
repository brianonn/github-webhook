FROM golang:1.21 as builder

ADD . /build

# Build the github-webhook command inside the container.
# the Makefile will take care of dependancies too
RUN make -C /build all

FROM scratch
COPY --from=builder /build/github-webhook /github-webhook

# Run the github-webhook command by default when the container starts.
USER 10000
WORKDIR /
CMD [ "/github-webhook" ]

# Document that the service listens on port 5000.
EXPOSE 5000

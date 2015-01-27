# Dockerfile

# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
ADD . /go/src/app

# Fetch dependencies.
# Build the crawler command inside the container
RUN go get golang.org/x/tools/cmd/vet
RUN make

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/src/app/bin/crawler

# Document that the service listens on port 8080.
EXPOSE 8080

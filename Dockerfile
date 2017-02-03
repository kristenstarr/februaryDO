FROM ubuntu:latest

RUN apt-get update && apt-get install --no-install-recommends -y \
    ca-certificates \
    curl \
    mercurial \
    git-core

# Obtain latest version of Golang (currently 1.7)
# Update package index and install go + git
RUN apt-get update -q && apt-get install -yq golang git-core

# Set up GOPATH
RUN mkdir /go
ENV GOPATH /go

EXPOSE 8080

WORKDIR /go/src/github.com/kristenfelch/pkgindexer

COPY . .

# Build our API.
RUN go build

# Use a CMD here, instead of ENTRYPOINT, for easy overwrite in docker ecosystem.
CMD go run main.go

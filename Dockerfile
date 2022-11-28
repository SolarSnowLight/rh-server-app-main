FROM golang

# WORKDIR /server-app-main

RUN go version
ENV GOPATH=/

COPY ./ ./

# Update packages
RUN apt-get update

# Install postgresql-client
RUN apt-get -y install postgresql-client

# Make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# Download all extensions
RUN go mod download

# Build go app
RUN go build -o server-app-main ./cmd/main.go

# Start go app
CMD ["./server-app-main"]
FROM golang:1.14 as builder
WORKDIR /go/src/github.com/vpoletaev11/fileHostingSite

# Copy the local project files to the container's workspace.
ADD . /go/src/github.com/vpoletaev11/fileHostingSite

# Build the project inside the container.
# RUN GOOS=linux go build  .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o .


# Execute the binary
FROM scratch
EXPOSE 8080

COPY --from=builder /go/src/github.com/vpoletaev11/fileHostingSite   /
ENTRYPOINT ["/fileHostingSite"]
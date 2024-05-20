# Start from the latest golang base image
FROM golang:1.21-bullseye as base

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download && go mod verify

# Copy the source from the current directory to the Working Directory inside the container
COPY cmd/repositoryembedder/ .
COPY pkg/common/ pkg/common/
COPY pkg/cache/ pkg/cache/
COPY pkg/embedder/ pkg/embedder/
COPY pkg/database/ pkg/database/
COPY api/ api/

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

FROM gcr.io/distroless/static-debian11

COPY --from=base /main .

# Command to run the executable
CMD ["./main"]

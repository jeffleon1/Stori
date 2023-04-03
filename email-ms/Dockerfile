
FROM golang AS builder

WORKDIR /src

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64


# Download dependencies
COPY go.mod go.sum /
RUN go mod download

# Add source code
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/api

# Multi-Stage production build
FROM alpine AS production
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
# Retrieve the binary from the previous stage
COPY --from=builder /go/bin/app /go/bin/app
RUN chmod a+x /go/bin/app

# Set the binary as the entrypoint of the container
ENTRYPOINT ["/go/bin/app"]
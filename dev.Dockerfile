FROM golang:1.22-alpine3.19 AS builder
ENV CGO_ENABLED 0
ARG BUILD_REF

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

# Copy the source code into the container.
COPY . /service
# Build the service binary.
WORKDIR /service/cmd/api
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
    go build -ldflags "-X main.build=${BUILD_REF}" -o server

# Run the Go Binary in Alpine.
FROM alpine:3.19
ARG BUILD_DATE
ARG BUILD_REF
WORKDIR /service
COPY --from=builder /service/cmd/api/server .
CMD ["./server"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="buy-better-admin-api" \
      org.opencontainers.image.revision="${BUILD_REF}"

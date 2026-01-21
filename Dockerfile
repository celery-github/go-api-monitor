# Build
FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/monitor ./cmd/monitor

# Run
FROM gcr.io/distroless/static:nonroot
COPY --from=build /bin/monitor /monitor
WORKDIR /
USER nonroot:nonroot
ENTRYPOINT ["/monitor"]

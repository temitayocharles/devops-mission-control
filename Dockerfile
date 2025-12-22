FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ops-tool .

# Distroless base for minimal size
FROM gcr.io/distroless/base-debian11

COPY --from=builder /app/ops-tool /
ENTRYPOINT ["/ops-tool"]

# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o kokovox .

# Runtime stage
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/kokovox /kokovox

EXPOSE 5108

USER nonroot:nonroot

ENTRYPOINT ["/kokovox"]

# Multi-stage build for the PromptSieve Go API
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /promptsieve-server ./cmd/promptsieve-server

FROM gcr.io/distroless/static-debian12
LABEL org.opencontainers.image.title="PromptSieve API"
LABEL org.opencontainers.image.description="Open-source sensitive data scanner for AI prompts"
COPY --from=builder /promptsieve-server /
EXPOSE 8080
ENTRYPOINT ["/promptsieve-server"]

FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /rocketchat-mcp ./cmd/rocketchat-mcp/

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /rocketchat-mcp /rocketchat-mcp
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/rocketchat-mcp"]
CMD ["-transport", "http"]

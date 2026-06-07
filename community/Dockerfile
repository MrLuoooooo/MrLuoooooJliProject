FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ .
RUN CGO_ENABLED=0 go build -o /community-server .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /community-server /community-server
COPY server/config.yaml /config.yaml
RUN mkdir -p /tmp/logs /tmp/uploads

EXPOSE 1807
CMD ["/community-server"]

FROM golang:1.27.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go test -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /short-urls ./cmd/short-urls

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /short-urls .
EXPOSE 8080
CMD ["./short-urls"]
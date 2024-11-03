FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o build/do-dyndns

FROM alpine
COPY --from=builder /src/build/do-dyndns /usr/local/bin/
CMD ["/usr/local/bin/do-dyndns"]

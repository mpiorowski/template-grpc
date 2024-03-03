#  Development
FROM golang:1.22 as dev

WORKDIR /server

COPY go.mod go.sum /server/
RUN go mod download

# Build
FROM golang:1.22 as build
WORKDIR /server

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux go build -v -o / /server/...

# Production
FROM debian:bookworm-slim as prod
WORKDIR /

COPY --from=build /powerit /powerit

# Install tls certificates
RUN apt-get update && apt-get install -y ca-certificates

CMD ["/service-auth"]

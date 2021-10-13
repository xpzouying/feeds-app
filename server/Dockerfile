FROM golang:latest AS builder
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /build/feeds-server ./cmd/main.go


FROM alpine:latest
COPY --from=builder /build/feeds-server /feeds-server
EXPOSE 8080
LABEL Name=feeds-server Version=0.0.3

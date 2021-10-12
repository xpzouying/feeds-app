FROM golang AS builder
WORKDIR /go/src/app
COPY ./server/ ./server/
COPY go.mod .
COPY go.sum .
RUN CGO_ENABLED=0 go build -o /app/feeds-server ./server/cmd/main.go


FROM alpine:latest
COPY --from=builder /app/feeds-server /feeds-server
LABEL Name=feeds-server Version=0.0.1

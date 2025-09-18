FROM golang:1.25-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o carrent cmd/carrent/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/carrent .
EXPOSE 8080
CMD ["./carrent", "--debug"]
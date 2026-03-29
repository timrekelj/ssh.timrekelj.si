FROM golang:1.26.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o ssh-portfolio .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/ssh-portfolio .
RUN mkdir -p .ssh
EXPOSE 22
CMD ["./ssh-portfolio"]

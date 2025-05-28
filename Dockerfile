FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

FROM scratch
WORKDIR /app
COPY --from=builder /app/app /app/app
COPY config/config.yml /app/config.yml
EXPOSE 8080
ENTRYPOINT ["/app/app"]
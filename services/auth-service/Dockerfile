FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service main.go

FROM scratch
COPY --from=builder /auth-service /auth-service
USER 10001:10001
EXPOSE 8081
ENTRYPOINT ["/auth-service"]

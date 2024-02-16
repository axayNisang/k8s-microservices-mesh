FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /catalog-service main.go

FROM scratch
COPY --from=builder /catalog-service /catalog-service
USER 10001:10001
EXPOSE 8082
ENTRYPOINT ["/catalog-service"]

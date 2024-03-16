FROM golang:1.21.1-alpine AS test_application
WORKDIR /src
COPY ./src/ /src/
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.2 && swag init
CMD ["go", "test", "-v", "./..."]

FROM test_application AS builder
RUN go build -o /app/stealthy-backend

FROM golang:1.21.1-alpine AS application
ARG UID=1001
RUN adduser -u $UID -D app-user
COPY --from=builder /app/stealthy-backend /app/stealthy-backend
WORKDIR /app
RUN chown -R app-user:app-user /app && chmod u+x stealthy-backend
COPY ./config.yaml .
USER app-user
CMD ["/app/stealthy-backend"]

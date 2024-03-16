FROM golang:1.21.1-alpine AS check_code_app
RUN go install github.com/qiniu/checkstyle/gocheckstyle@v0.082 && \
    go install honnef.co/go/tools/cmd/staticcheck@2023.1.7
WORKDIR /app
COPY checkstyle.json /
CMD gocheckstyle -config /checkstyle.json ./ && \
    staticcheck ./...

FROM check_code_app as test_application
COPY ./src/ ./
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.2 && swag init
CMD go test -v ./... && \
    gocheckstyle -config /checkstyle.json ./ && \
    staticcheck ./...

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

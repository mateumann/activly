FROM golang:1.21.6-alpine3.18 AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./activly ./cmd/main.go

FROM alpine:latest
ARG username=activly
ARG groupname=users
ENV POSTGRES_HOST=postgres
ENV POSTGRES_PORT=5432
ENV POSTGRES_DB=activly
ENV POSTGRES_USER=activly
ENV POSTGRES_PASSWORD=""
RUN adduser -h /app -s /dev/null -G $groupname -D -u 20080 $username
WORKDIR /app
COPY --from=builder --chown=$username:$groupname /build/activly ./activly
EXPOSE 8080
USER 20080
ENTRYPOINT ["/app/activly"]

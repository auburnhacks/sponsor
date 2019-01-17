FROM golang:1.11 as builder
WORKDIR /build
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix nocgo \
    -o sponsor-server ./cmd/sponsor-server

FROM alpine:latest
WORKDIR /app
RUN apk update \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*
COPY --from=builder /build/sponsor-server ./
COPY ./migrations/. ./migrations/
COPY ./jwt_key_prod ./jwt_key_dev
RUN chmod +x sponsor-server
EXPOSE 8080
EXPOSE 10000
ENTRYPOINT [ "./sponsor-server" ]

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
RUN chmod +x sponsor-server
EXPOSE 8080
EXPOSE 10000
ENTRYPOINT [ "./sponsor-server" ]

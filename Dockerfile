FROM golang:1.12.6-alpine3.10 as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN apk add --no-cache git
RUN go build -o ./bin/server/server ./cmd/server

FROM alpine:3.10
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/ /app/

WORKDIR /app/bin/server
ENTRYPOINT ["./server"]
CMD []

EXPOSE 3000
EXPOSE 8080

FROM golang:1.21.1 AS builder


# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/kgantsov/todogo
COPY ./ ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
WORKDIR $GOPATH/src/github.com/kgantsov/todogo/
RUN swag init -g ./cmd/server/main.go -o ./docs
WORKDIR $GOPATH/src/github.com/kgantsov/todogo/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

FROM alpine

ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /

COPY --from=builder /app /
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/app --port 8780"]

FROM golang:1.9 AS builder

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/kgantsov/todogo
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY ./ ./
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

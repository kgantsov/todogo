FROM golang:latest

RUN mkdir -p /go/src/github.com/kgantsov/todogo/
ADD . /go/src/github.com/kgantsov/todogo/
WORKDIR /go/src/github.com/kgantsov/todogo/

RUN go get github.com/lib/pq
RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/jinzhu/gorm
RUN go get gopkg.in/gin-gonic/gin.v1
RUN go build

CMD ["/go/src/github.com/kgantsov/todogo/todogo"]
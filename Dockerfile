FROM golang:1.13
MAINTAINER Diode "diodebupt@163.com"
WORKDIR $GOPATH/src/github.com/Diode222/GomokuGameReferee
ADD . $GOPATH/src/github.com/Diode222/GomokuGameReferee
ENV GO111MODULE on
ENV GOPROXY https://goproxy.io
RUN go mod download && go build main.go
EXPOSE 10000
ENTRYPOINT  ["./main"]

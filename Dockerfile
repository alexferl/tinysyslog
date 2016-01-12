FROM golang

ADD . /go/src/github.com/admiralobvious/tinysyslog

RUN go get github.com/tools/godep
WORKDIR /go/src/github.com/admiralobvious/tinysyslog
RUN godep restore
RUN go install

RUN mkdir -p /logs

ENTRYPOINT /go/bin/tinysyslog --address 0.0.0.0:5140 --filesystem-filename /logs/syslog.log --log-file stdout

EXPOSE 5140 5140/udp

FROM golang

ADD . /go/src/github.com/admiralobvious/tinysyslog

WORKDIR /go/src/github.com/admiralobvious/tinysyslog
RUN go install

RUN mkdir -p /data/logs

ENTRYPOINT /go/bin/tinysyslog --address 0.0.0.0:5140 --filesystem-filename /data/logs/syslog.log --log-file stdout

EXPOSE 5140 5140/udp

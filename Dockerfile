FROM golang:onbuild

MAINTAINER Alexandre Ferland <aferlandqc@gmail.com>

RUN mkdir /data/logs

EXPOSE 5140 5140/udp

CMD app --address 0.0.0.0:5140 --filesystem-filename /data/logs/syslog.log --log-file stdout

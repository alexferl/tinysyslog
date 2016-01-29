FROM golang:onbuild

RUN mkdir /logs && chown nobody:nogroup -R /logs
USER nobody

EXPOSE 5140 5140/udp
CMD app --address 0.0.0.0:5140 --filesystem-filename /logs/syslog.log --log-file stdout

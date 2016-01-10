# tinysyslog
[![Go Report Card](http://goreportcard.com/badge/admiralobvious/seau)](http://goreportcard.com/report/admiralobvious/seau)

A tiny and simple syslog server with log rotation.

## Quickstart
To install tinysyslog itself:

    go get -u github.com/admiralobvious/tinysyslog
And then to run it (from your $GOPATH/bin folder):

    ./tinysyslog --log-file stdout
If tinysyslog started properly you should see:
```
INFO[0000] tinysyslog listening on 0.0.0.0:5140
```

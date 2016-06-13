# tinysyslog
[![Go Report Card](http://goreportcard.com/badge/admiralobvious/tinysyslog)](http://goreportcard.com/report/admiralobvious/tinysyslog)

A tiny and simple syslog server with log rotation. tinysyslog was born out of the need for a tiny (the binary is currently <10MB in size), easy to setup and use syslog server that simply writes every incoming log (RFC5424 format) to a file that is automatically rotated. tinysyslog is based on [go-syslog](https://github.com/mcuadros/go-syslog) and [lumberjack](https://github.com/natefinch/lumberjack).

## Quickstart
tinysyslog requires golang to work.

To install it on OS X:

    brew install go
On Ubuntu/Debian:

    sudo apt-get install golang
You will also need to set the GOPATH, see: https://golang.org/doc/code.html#GOPATH

To install tinysyslog itself:

    go get -u github.com/admiralobvious/tinysyslog
And then to run it (from your $GOPATH/bin folder):

    ./tinysyslog --log-file stdout
If tinysyslog started properly you should see:
```
INFO[0000] tinysyslog listening on 127.0.0.1:5140
```
You can take make sure logs are saved to the file by entering the following in a terminal:
```
nc -w0 -u 127.0.0.1 5140 <<< '<165>1 2016-01-01T12:01:21Z hostname appname 1234 ID47 [exampleSDID@32473 iut="9" eventSource="test" eventID="123"] message'
```

You should then see the following in `syslog.log`:
```
Jan  1 12:01:21 hostname appname[1234]: message
```

## Configuration
```
type Config struct {
	Address     string
	Filesystem  Filesystem
	LogFile     string
	LogFormat   string
	LogLevel    string
	MutatorType string
	SinkType    string
	SocketType  string
}

type Filesystem struct {
	Filename   string
	MaxAge     int
	MaxBackups int
	MaxSize    int
}
```

### Address
`--address`

IP and port to listen on. (default "127.0.0.1:5140")
### Filesystem
#### Filename
`--filesystem-filename`

File to write incoming logs to. (default "syslog.log")
#### MaxAge
`--filesystem-max-age`

Maximum age (in days) before a log is deleted. Set to '0' to disable. (default 30)
#### MaxBackups
`--filesystem-max-backups`

Maximum backups to keep. Set to '0' to disable. (default 10)
#### MaxSize
`--filesystem-max-size`

Maximum log size (in megabytes) before it's rotated. (default 100)
### LogFile
`--log-file`

The log file to write to. 'stdout' means log to stdout and 'stderr' means log to stderr. (default "tinysyslog.log")
### LogFormat
`--log-format`

The log format. Valid format values are: text, json. (default "text")
### LogLevel
`--log-level`

The granularity of log outputs. Valid level names are: debug, info, warning, error and critical. (default "info")
### MutatorType
`--mutator-type`

Mutator to transform logs as. Valid format values are: text, json. (default "text")
### SinkType
`--sink-type`

Sink to save logs to. (default "filesystem")
### SocketType
`--socket-type`

Type of socket to use, TCP or UDP. If no type is specified, both are used. (default "")

## Benchmarks
Nothing scientific here but with a simple client consisting of a for loop sending large messages as fast as possible over UDP:

`iostat -d 5`
```
    KB/t tps  MB/s
  127.61 585 72.95
  127.66 592 73.74
  126.41 591 72.98
  126.36 590 72.76
  124.76 615 74.95
```

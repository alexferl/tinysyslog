# tinysyslog
[![Go Report Card](http://goreportcard.com/badge/admiralobvious/tinysyslog)](http://goreportcard.com/report/admiralobvious/tinysyslog)

A tiny and simple syslog server with log rotation. tinysyslog was born out of the need for a tiny (the binary is currently <10MB in size), easy to setup and use syslog server that simply writes every incoming log (RFC5424 format) to a file that is automatically rotated. tinysyslog is based on [go-syslog](https://github.com/mcuadros/go-syslog) and [lumberjack](https://github.com/natefinch/lumberjack).

## Quickstart
To install tinysyslog:

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
Usage of ./tinysyslog:
      --address string               IP and port to listen on. (default "127.0.0.1:5140")
      --console-output string        Console to output too. Valid outputs are: stdout, stderr. (default "stdout")
      --filesystem-filename string   File to write incoming logs to. (default "syslog.log")
      --filesystem-max-age int       Maximum age (in days) before a log is deleted. (default 30)
      --filesystem-max-backups int   Maximum backups to keep. (default 10)
      --filesystem-max-size int      Maximum log size (in megabytes) before it's rotated. (default 100)
      --filter-type string           Filter to filter logs with. Valid filters are: regex. (default "regex")
      --log-file string              The log file to write to. 'stdout' means log to stdout and 'stderr' means log to stderr. (default "tinysyslog.log")
      --log-format string            The log format. Valid format values are: text, json. (default "text")
      --log-level string             The granularity of log outputs. Valid level names are: debug, info, warning, error and critical. (default "info")
      --mutator-type string          Mutator type to use. Valid mutators are: text, json. (default "text")
      --regex-filter string          Regex to filter with. No filtering by default.
      --sink-type string             Sink to save logs to. Valid sinks are: console, filesystem. (default "filesystem")
      --socket-type string           Type of socket to use, TCP or UDP. If no type is specified, both are used.
```

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
